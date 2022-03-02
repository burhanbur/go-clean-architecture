package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"blog/utils"

	"github.com/jinzhu/gorm"
)

type Pagination struct {
	Limit        int `json:"limit"`
	Page         int `json:"page"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

func (p *Pagination) PageData(offset int, limit int, totalRecord int) (*Pagination, error) {
	input := float64(totalRecord) / float64(limit)
	totalPage, err2 := p.TotalPage(input)

	if err2 != nil {
		return nil, err2
	}

	p.TotalRecords = totalRecord
	p.TotalPages = totalPage
	p.Page = offset
	p.Limit = limit

	return p, nil
}

func (p *Pagination) CountRecords(db *gorm.DB, tableName string, limit int) error {
	var count int
	err := db.Debug().Table(tableName).Count(&count)

	if err != nil {
		return err.Error
	}

	fmt.Println(count)
	input := float64(count) / float64(limit)
	totalPage, err2 := p.TotalPage(input)

	if err2 != nil {
		return err2
	}

	p.TotalRecords = count
	p.TotalPages = totalPage

	return nil
}

func (p *Pagination) TotalPage(inputNum float64) (int, error) {
	var err error
	result := 0

	utils.This(func() {
		s := strconv.FormatFloat(inputNum, 'f', -1, 64)
		z := strings.Split(s, string('.'))
		value1, _ := strconv.Atoi(z[0])
		var value int

		if value1 == 0 {
			value = value1
		}

		if len(z) >= 2 {
			value = value1 + 1
		} else {
			value = value1
		}

		result = value
		err = nil
	}).Catch(func(e utils.E) {
		result = 0
		err = errors.New("error get total page")
	})

	return result, err
}
