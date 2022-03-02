package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SALT_BYTE_SIZE    = 24
	HASH_BYTE_SIZE    = 24
	PBKDF2_ITERATIONS = 1000
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func EncodeToString(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func OffsetLimit(offset, limit string) (map[string]int, error) {
	var number map[string]int
	// fmt.Println(offset, limit)
	number = map[string]int{}
	var errR error
	newOffset := 0
	newLimit := 0

	This(func() {
		newOffset, _ = strconv.Atoi(offset)
		newLimit, _ = strconv.Atoi(limit)
	}).Catch(func(err E) {
		log.Println("[account_usecase] Error string to int", err)
		errR = errors.New("")
	})

	newOffset -= 1
	if errR != nil || newLimit < 0 || newOffset < 0 {
		return nil, errR
	}

	number["limit"] = newLimit
	if newOffset > 0 {
		number["offset"] = newOffset * newLimit
		return number, nil
	}

	number["offset"] = newOffset
	fmt.Println(number)

	return number, nil
}

func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

func Hash(password string) (string, error) {
	salt := make([]byte, SALT_BYTE_SIZE)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		Log{}.Error(err.Error())
		return "", errors.New("Err generating random salt")
	}

	//todo: enhance: randomize itrs as well
	hbts := pbkdf2.Key([]byte(password), salt, PBKDF2_ITERATIONS, HASH_BYTE_SIZE, sha1.New)
	//hbtstr := fmt.Sprintf("%x", hbts)

	return fmt.Sprintf("%v:%v:%v",
		PBKDF2_ITERATIONS,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hbts)), nil
}

func Verify(raw, hash string) (bool, error) {
	hparts := strings.Split(hash, ":")

	itr, err := strconv.Atoi(hparts[0])
	if err != nil {
		fmt.Printf("wrong hash %v", hash)
		return false, errors.New("wrong hash, iteration is invalid")
	}
	salt, err := base64.StdEncoding.DecodeString(hparts[1])
	if err != nil {
		fmt.Print("wrong hash, salt error:", err)
		return false, errors.New("wrong hash, salt error:" + err.Error())
	}

	hsh, err := base64.StdEncoding.DecodeString(hparts[2])
	if err != nil {
		fmt.Print("wrong hash, hash error:", err)
		return false, errors.New("wrong hash, hash error:" + err.Error())
	}

	rhash := pbkdf2.Key([]byte(raw), salt, itr, len(hsh), sha1.New)
	return equal(rhash, hsh), nil
}

//bytes comparisons
func equal(h1, h2 []byte) bool {
	diff := uint32(len(h1)) ^ uint32(len(h2))
	for i := 0; i < len(h1) && i < len(h2); i++ {
		diff |= uint32(h1[i] ^ h2[i])
	}

	return diff == 0
}

func IsValidAssignTaskDate(startDate, endDate, todayDate time.Time) bool {
	for {
		if todayDate.After(endDate) {
			break
		}
		if todayDate.Before(startDate) {
			break
		}
		if todayDate.Equal(startDate) {
			return true
		}
		if todayDate.Equal(endDate) {
			return true
		}
		if todayDate.After(startDate) && todayDate.Before(endDate) {
			return true
		}
		startDate = startDate.Add(time.Hour * 24)
	}
	return false
}

func IsValidRoutineTask(startDate, todayDate time.Time, repeatTimes int, repeatEvery string) bool {
	for {
		if startDate.Equal(todayDate) {
			return true
		}
		if startDate.After(todayDate) {
			break
		}
		switch repeatEvery {
		case "D":
			startDate = startDate.Add(time.Hour * (24 * time.Duration(repeatTimes)))
		case "W":
			startDate = startDate.AddDate(0, 0, 7*repeatTimes)
		case "M":
			startDate = startDate.AddDate(0, 1*repeatTimes, 0)
		case "Y":
			startDate = startDate.AddDate(1*repeatTimes, 0, 0)
		}
	}
	return false
}

func IsContains(arrays []string, values string) bool {
	for _, val := range arrays {
		if val == values {
			return true
		}
	}
	return false
}

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func getenvInt(key string) (int, error) {
	s, err := getenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func getenvBool(key string) (bool, error) {
	s, err := getenvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
