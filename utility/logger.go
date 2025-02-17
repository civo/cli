package utility

import (
	"fmt"
)

func Println(args ...interface{}) {
	fmt.Println(args...)
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
