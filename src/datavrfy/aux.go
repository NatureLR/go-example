package main

import (
	"fmt"

	"strconv"
	"strings"
)

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
	return fmt.Sprintf("%0."+strconv.Itoa(decimal)+"f", p)
}
