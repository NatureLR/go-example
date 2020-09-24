package tools

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//冒泡排序
func sortArr(arr []int, size int) []int {
	for i := 0; i < size; i++ {
		for j := 0; j < (size - 1 - i); j++ {
			if arr[j] > arr[j+1] {
				tmp := arr[j+1]
				arr[j+1] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func fmtTime(ts string) string {
	t, err := strconv.Atoi(ts)
	if err != nil {
		return ts
	}
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

func fmtDate(ts string) string {
	t, err := strconv.Atoi(ts)
	if err != nil {
		return ts
	}
	return time.Unix(int64(t), 0).Format("2006-01-02")
}

func fmtWeekDayTime(ts string) string {
	t, err := strconv.Atoi(ts)
	if err != nil {
		return ""
	}
	switch time.Unix(int64(t), 0).Weekday() {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	case time.Sunday:
		return "周日"
	}
	return ""
}

func fmtWeekDayDate(ts string) string {
	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		return ""
	}
	switch t.Weekday() {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	case time.Sunday:
		return "周日"
	}
	return ""
}

func parseTime(closetime, ordertime string) string {
	closeT := strings.Split(closetime, " ")
	orderT := strings.Split(ordertime, " ")
	if len(closeT) < 2 {
		return ""
	}
	return orderT[0] + " " + closeT[1]
}

func num(vs ...string) ([]float64, error) {
	//normalize numeric values
	_n := func(v string) string {
		v = strings.TrimSpace(v)
		if v == "" {
			return "0"
		}
		return v
	}
	var ns []float64
	for _, v := range vs {
		n, err := strconv.ParseFloat(_n(v), 64)
		if err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func sum(v1, v2 string) (float64, error) {
	n, err := num(v1, v2)
	if err != nil {
		return 0, err
	}
	return n[0] + n[1], nil
}

func sums(v1, v2 string, decimal int) string {
	s, err := sum(v1, v2)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%0."+strconv.Itoa(decimal)+"f", s)
}

func diff(v1, v2 string) (float64, error) {
	n, err := num(v1, v2)
	if err != nil {
		return 0, err
	}
	return n[0] - n[1], nil
}

func diffs(v1, v2 string, decimal int) string {
	d, err := diff(v1, v2)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%0."+strconv.Itoa(decimal)+"f", d)
}

func prod(v1, v2 string) (float64, error) {
	n, err := num(v1, v2)
	if err != nil {
		return 0, err
	}
	return n[0] * n[1], nil
}

func prods(v1, v2 string, decimal int) string {
	p, err := prod(v1, v2)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%f", p)
	//return fmt.Sprintf("%0."+strconv.Itoa(decimal)+"f", p)
}

func fmtMinute(second string) string {
	result, _ := strconv.ParseFloat(second, 64)
	r := result / 60
	return fmt.Sprintf("%0.0f", r)
}
