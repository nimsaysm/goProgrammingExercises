package chapter02

import (
	"fmt"
	"os"
	"strconv"
)

func Runner() {
	fmt.Println("*** Chapter 02 ***")
	ConvertToKelvin()
	PopCountWithLoop(255) // The return must be 8
}

func ConvertToKelvin() {
	fmt.Println("*** Running ConvertToKelvin ***")

	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cf: %v\n", err)
			os.Exit(1)
		}
		c := Celsius(t)
		k := Kelvin(t)
		fmt.Printf("%s = %s, %s = %s\n",
			c, CToK(c), k, KToC(k))
	}
}
