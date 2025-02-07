package chapter03

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

func Runner() {
	fmt.Println("*** Chapter 03 ***")

	InsertCommasInDecimal()
	CheckAnagrams()
}

func InsertCommasInDecimal() {
	fmt.Println("*** Running InsertCommasInDecimal ***")

	fmt.Println(commaWithFloatAndSignal("1"))
	fmt.Println(commaWithFloatAndSignal("12"))
	fmt.Println(commaWithFloatAndSignal("123"))
	fmt.Println(commaWithFloatAndSignal("1234"))
	fmt.Println(commaWithFloatAndSignal("12345"))
}

func comma(s string) string {
	var buf bytes.Buffer

	n := len(s) % 3
	if n == 0 {
		n = 3
	}

	for len(s) > 3 {
		buf.WriteString(s[:1])
		buf.WriteString(",")
		s = s[n:]
		n = 3
	}

	buf.WriteString(s)
	return buf.String()
}

func commaWithFloatAndSignal(s string) string {
	if strings.Contains(s, ".") || strings.Contains(s, "-") {
		return s
	}

	return comma(s)
}

func CheckAnagrams() {
	fmt.Println("*** Running CheckAnagrams ***")

	fmt.Println(isAnagram("dog", "god"))
	fmt.Println(isAnagram("banana", "orange"))
	fmt.Println(isAnagram("elbow", "below"))
}

func isAnagram(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	s1Array := []rune(s1)
	sort.Slice(s1Array, func(i, j int) bool {
		return s1Array[i] < s1Array[j]
	})

	s2Array := []rune(s2)
	sort.Slice(s2Array, func(i, j int) bool {
		return s2Array[i] < s2Array[j]
	})

	for i := 0; i < len(s1Array); i++ {
		if s1Array[i] != s2Array[i] {
			return false
		}
	}

	return true
}
