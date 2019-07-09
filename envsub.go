package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	input = flag.String("i", "", "")
)

var usage = `Usage: envsubst -i <filename>
Examples of template variable definition:
 ${ENV_VAR}            Value of $ENV_VAR or empty if variable is not set.
 ${ENV_VAR=default}    Value of $ENV_VAR or "default" if variable is not set.
 ${ENV_VAR-}           Value of $ENV_VAR or skip whole line if variable is not set.
 `

// Evar environment variable handler
type Evar struct {
	start        uint
	end          uint
	name         string
	initialized  bool
	isDefault    bool
	defaultValue string
	willSkip     bool
}

// Init expression
func (e *Evar) Init(position uint) {
	e.Clear()
	e.start = position
	e.initialized = true
}

// Append add char to Evar (name of default value)
func (e *Evar) Append(ch string) {
	if !e.isDefault {
		e.name += ch
	} else {
		e.defaultValue += ch
	}
}

// End close expression
func (e *Evar) End(position uint) {
	e.end = position
}

// Clear prepares Evar for next expression
func (e *Evar) Clear() {
	e.start = 0
	e.end = 0
	e.name = ""
	e.initialized = false
	e.isDefault = false
	e.defaultValue = ""
	e.willSkip = false
}

// GetValue returns final value of expression
func (e *Evar) GetValue() string {
	envValue := os.Getenv(e.name)
	if len(envValue) > 0 {
		return envValue
	}

	if len(e.defaultValue) > 0 {
		return e.defaultValue
	}

	return ""
}

// Line handler
type Line struct {
	inputLine  string
	outputLine string
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}

	flag.Parse()

	var reader *bufio.Reader
	if *input != "" {
		file, err := os.Open(*input)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error to open file input: %s.", *input))
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				substInLine(line)
				break
			}

		}
		newLine, skip := substInLine(line)
		if !skip {
			fmt.Fprint(os.Stdout, newLine)
		}
	}
}

// substInLine returns line with substituted variables as string and false
// or empty string and tru in case line have to be skiped.
func substInLine(line string) (string, bool) {
	env := &Evar{}
	newLine := ""

	for idx, ch := range line {
		switch sw := ch; {
		case sw == '$':
			env.Init(uint(idx))

		case sw == '{':

		case sw == '}':
			env.End(uint(idx))
			//fmt.Println("Value: ", env.name, "Skip", fmt.Sprintf("%b", env.willSkip))
			value := env.GetValue()
			if len(value) == 0 && env.willSkip {
				env.Clear()
				return "", true
			}
			newLine += value
			env.Clear()

		case sw == '=' && env.initialized:
			env.isDefault = true

		case sw == '-' && env.initialized && !env.isDefault:
			env.willSkip = true

		default:
			if env.initialized {
				env.Append(fmt.Sprintf("%c", ch))
			} else {
				newLine += fmt.Sprintf("%c", ch)
			}
		}
	}
	return newLine, false
}
