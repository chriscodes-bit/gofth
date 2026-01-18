package main

import (
	"fmt"
	"strconv"
	"strings"
)

var Stack []any
var output Output
var compiling bool
var currentWord string
var currentDef []string
var UserWords = Dictionary{}

type Dictionary = map[string]Word

type Output struct {
	lines []string
}

func (o *Output) Write(s string) {
	o.lines = append(o.lines, s)
}

func (o *Output) Clear() {
	o.lines = nil
}

func (o *Output) String() string {
	return strings.Join(o.lines, "\n")
}

func binaryOp(op func(a, b float64) float64) {
	b := pop()
	a := pop()
	af, aok := a.(float64)
	bf, bok := b.(float64)
	if !aok || !bok {
		output.Write("Error: operation requires two numbers")
		return
	}
	push(op(af, bf))
}
func push(v any) {
	Stack = append(Stack, v)
}
func pop() any {
	if len(Stack) == 0 {
		return nil
	}
	v := Stack[len(Stack)-1]
	Stack = Stack[:len(Stack)-1]
	return v
}
func parseNumber(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil
}
func printNum(v any) {
	if f, ok := v.(float64); ok {
		if f == float64(int(f)) {
			output.Write(fmt.Sprintf("%d", int(f)))
		} else {
			output.Write(fmt.Sprintf("%g", f))
		}
	} else {
		output.Write(fmt.Sprint(v))
	}
}

func makeWord(def []string) Word {
	// capture definition
	definition := make([]string, len(def))
	copy(definition, def)

	return Word{
		Name:        currentWord,
		Description: "User-defined word",
		Category:    "User",
		Execute: func() any {
			for _, cmd := range definition {
				if word, exists := UserWords[cmd]; exists {
					word.Execute()
				} else if word, exists := Builtins[cmd]; exists {
					word.Execute()
				} else if val, ok := parseNumber(cmd); ok {
					push(val)
				} else {
					output.Write(fmt.Sprintf("Unknown command in word: %s", cmd))
				}
			}
			return nil
		},
	}

}

func parseForthCode(code string) {
	commands := strings.FieldsSeq(code)
	for cmd := range commands {
		if compiling {
			if cmd == ";" {
				UserWords[currentWord] = makeWord(currentDef)
				compiling = false
				currentWord = ""
				currentDef = nil
			} else if currentWord == "" {
				currentWord = cmd
			} else {
				currentDef = append(currentDef, cmd)
			}
		} else {
			// Normal execution
			if cmd == ":" {
				compiling = true
			} else if word, exists := UserWords[cmd]; exists {
				word.Execute()
			} else if word, exists := Builtins[cmd]; exists {
				word.Execute()
			} else if val, ok := parseNumber(cmd); ok {
				push(val)
			} else {
				output.Write(fmt.Sprintf("Unknown command: %s", cmd))
			}
		}
	}
}
