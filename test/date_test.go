package test

import (
	"fmt"
	"github.com/metakeule/fmtdate"
	"testing"
)

func TestDate(t *testing.T) {
	date := "2016-01-24T163745+0000"
	startTime, err := fmtdate.Parse("YYYY-MM-DDThhmmss+0000", date)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(startTime)
}
