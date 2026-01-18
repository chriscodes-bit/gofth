package main

import (
	"fmt"
	"strings"
)

type Word struct {
	Name        string
	Description string
	Category    string
	Execute     func() any
}

var Builtins = map[string]Word{
	".": {
		Name:        ".",
		Description: "Pop and print top of stack",
		Category:    "Output",
		Execute: func() any {
			v := pop()
			printNum(v)
			return nil
		},
	},
	".s": {
		Name:        ".",
		Description: "Print entire stack non-destructively",
		Category:    "Output",
		Execute: func() any {
			var parts []string
			parts = append(parts, "Stack:")
			for _, v := range Stack {
				if f, ok := v.(float64); ok {
					if f == float64(int(f)) {
						parts = append(parts, fmt.Sprintf("%d", int(f)))
					} else {
						parts = append(parts, fmt.Sprintf("%g", f))
					}
				} else {
					parts = append(parts, fmt.Sprint(v))
				}
			}
			output.Write(strings.Join(parts, " "))
			return nil
		},
	},
	"dup": {
		Name:        "dup",
		Description: "Duplicate top of the stack",
		Category:    "Stack manipulation",
		Execute: func() any {
			v := pop()
			push(v)
			push(v)
			return nil
		},
	},
	"swap": {
		Name:        "swap",
		Description: "Swap the two top most values of the stack",
		Category:    "Stack manipulation",
		Execute: func() any {
			a := pop()
			b := pop()
			push(a)
			push(b)
			return nil
		},
	},
	"drop": {
		Name:        "drop",
		Description: "Destroys the top most value from the stack",
		Category:    "Stack manipulation",
		Execute: func() any {
			pop()
			return nil
		},
	},
	"+": {
		Name:        "+ (Addition)",
		Description: "Adds the two top most values of the stack together",
		Category:    "Math",
		Execute: func() any {
			binaryOp(func(a, b float64) float64 { return a + b })
			return nil
		},
	},
	"-": {
		Name:        "- (Subtraction)",
		Description: "Subtracts the two top most values of the stack together",
		Category:    "Math",
		Execute: func() any {
			binaryOp(func(a, b float64) float64 { return a - b })
			return nil
		},
	},
	"*": {
		Name:        "* (Multiplication)",
		Description: "Multiplies the two top most values of the stack together",
		Category:    "Math",
		Execute: func() any {
			binaryOp(func(a, b float64) float64 { return a * b })
			return nil
		},
	},
	"/": {
		Name:        "/ (Division)",
		Description: "Divides the two top most values of the stack together",
		Category:    "Math",
		Execute: func() any {
			binaryOp(func(a, b float64) float64 { return a / b })
			return nil
		},
	},
}
