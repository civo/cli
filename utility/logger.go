package utility

import (
	"fmt"

	"github.com/civo/cli/common"
)

func Println(args ...interface{}) {
	if !common.Quiet {
		fmt.Println(args...)
	}
}

func Printf(format string, args ...interface{}) {
	if !common.Quiet {
		fmt.Printf(format, args...)
	}
}
