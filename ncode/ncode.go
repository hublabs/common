package ncode

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	ranges = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length = 5
)

func Max() int64 {
	return int64(math.Pow(float64(len(ranges)), float64(length))) - 1
}
func SetMax(l int) {
	length = l
}
func Range() string {
	return ranges
}
func SetRange(r string) {
	ranges = r
}

func Encode(n int64) (string, error) {
	if n > Max() {
		return "", fmt.Errorf("Number is too large. Max: %d", Max())
	}
	var code string
	for i := length - 1; i >= 0; i-- {
		divisor := int64(math.Pow(float64(len(ranges)), float64(i)))
		quotient := n / divisor
		n = n % divisor
		code += string(ranges[quotient])
	}
	return code, nil
}

func Decode(code string) (int64, error) {
	if len(code) > length {
		return 0, errors.New("Invalid code")
	}
	var n int64
	for i := length; i > 0; i-- {
		n += int64(strings.IndexRune(ranges, rune(code[i-1]))) *
			int64(math.Pow(float64(len(ranges)), float64(length-i)))
	}
	return n, nil
}
