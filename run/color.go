package main

import (
	"github.com/fatih/color"
)

func colorOK(format string, args ...interface{}) string {
	return color.GreenString(format, args...)
}

func colorErr(format string, args ...interface{}) string {
	return color.RedString(format, args...)
}
