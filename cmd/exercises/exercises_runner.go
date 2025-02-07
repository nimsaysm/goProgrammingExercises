package main

import (
	"fmt"
	"github.com/nimsaysm/goProgrammingExercises/internal/chapter01"
	"github.com/nimsaysm/goProgrammingExercises/internal/chapter02"
	"github.com/nimsaysm/goProgrammingExercises/internal/chapter03"
)

func RunAllExercises() {
	fmt.Println("Choose which chapter run...")
	var input int
	fmt.Scan(&input)
	switch input {
	case 1:
		chapter01.Runner()
	case 2:
		chapter02.Runner()
	case 3:
		chapter03.Runner()
	default:
		fmt.Println("The chapter does not exist.")
	}
}
