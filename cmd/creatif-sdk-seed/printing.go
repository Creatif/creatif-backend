package main

import (
	"fmt"
	"github.com/fatih/color"
)

func printNewlineSandwich(printer *color.Color, print string) {
	fmt.Println("")
	printer.Println(print)
}
