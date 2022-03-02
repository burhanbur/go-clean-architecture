package utils

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// Log struct
type Log struct {
}

const (
	typeError = "[ERROR]"
	typeInfo  = "[INFO]"
)

func writeLog(logType string, message string) {
	pc, fn, line, _ := runtime.Caller(2)

	// create "logs" directory if not exist
	path := "logs"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	f, _ := os.OpenFile("logs/logs-"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	errorTime := time.Now()
	errorText := fmt.Sprintf(logType+" %d-%02d-%02d %02d:%02d:%02d \n %s[%s:%d] \n %v\n",
		errorTime.Year(),
		errorTime.Month(),
		errorTime.Day(),
		errorTime.Hour(),
		errorTime.Minute(),
		errorTime.Second(),
		runtime.FuncForPC(pc).Name(),
		fn,
		line,
		message,
	)

	f.WriteString(errorText)

	f.Close()
}

// Error method
func (Log) Error(s string) {
	fmt.Println(s)
	writeLog(typeError, s)
}

// Info method
func (Log) Info(s string) {
	fmt.Println(s)
	writeLog(typeInfo, s)
}

// ErrorOutput function
func ErrorOutput(c *gin.Context, e error) {
	c.JSON(
		http.StatusInternalServerError,
		gin.H{
			"status":  false,
			"message": e.Error(),
		},
	)

	return
}
