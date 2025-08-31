package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func unpackString(str string) (string, error) {
	res := []rune{}

	chars := []rune(str)
	for i := 0; i < len(chars); i++ {
		if n, err := strconv.Atoi(string(chars[i])); err == nil {
			if len(res) == 0 {
				return "", errors.New("string has no letters")
			}

			res = append(res, []rune(strings.Repeat(string(res[len(res)-1]), n-1))...)
		} else {
			if chars[i] == '\\' && i < len(chars)-1 {
				i++
			}
			res = append(res, chars[i])
		}
	}

	return string(res), nil
}

func main() {
	s, err := unpackString("a4bc2d5e")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s)
}
