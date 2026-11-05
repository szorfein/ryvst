package colour

import (
	"fmt"
	"math/rand"
	"strings"
)

var allColors = map[string]string{
	"reset":     "\033[0m",
	"bold":      "\033[1m",
	"underline": "\033[4m",
	"strike":    "\033[9m",
	"italic":    "\033[3m",
	"cRed":      "\033[31m",
	"cGreen":    "\033[32m",
	"cYellow":   "\033[33m",
	"cBlue":     "\033[34m",
	"cPurple":   "\033[35m",
	"cCyan":     "\033[36m",
	"cWhite":    "\033[37m",
}

var coloursArray = []string{
	"\033[31m",
	"\033[32m",
	"\033[33m",
	"\033[34m",
	"\033[35m",
	"\033[36m",
	"\033[37m",
}

// Usage:
// colour.Println("cRed", true, "time", "expire", "yes")
func Println(color string, clear bool, msg ...string) {
	if clear {
		fmt.Printf("%s%s %s\n", allColors[color], strings.Join(msg, " "), allColors["reset"])
	} else {
		fmt.Printf("%s%s\n", allColors[color], strings.Join(msg, " "))
	}
}

func RandPrintf(msg string, options ...string) {
	r := rand.Intn(5)
	fmt.Printf("%s%s%s", coloursArray[r], msg, allColors["reset"])
}
