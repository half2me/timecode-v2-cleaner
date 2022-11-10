package main

import (
	"bufio"
	"flag"
	"math"
	"os"
	"strconv"
	"strings"
)
import "fmt"

var precision = flag.Int("precision", 4, "max number of decimal places allowed in timecode")
var outFile = flag.String("out", "", "Output file with fixed timecode")

func almostEqual(a, b float64) bool {
	return (1 / math.Pow10(*precision+1)) > math.Abs(a-b)
}

func gte(a, b float64) bool {
	if !almostEqual(a, b) {
		return a > b
	}
	return true
}

func floatStr(a float64) string {
	return strconv.FormatFloat(a, 'f', *precision, 64)
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: timecode-cleaner [flags] <timecode-v2-file>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(flag.Arg(0))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	currentLineNumber := 0
	var previousTimestamp float64
	firstTimestamp := true

	foundErrors := false

	var writer *bufio.Writer

	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		writer = bufio.NewWriter(f)
		defer func(writer *bufio.Writer) {
			_ = writer.Flush()
		}(writer)
	}

	for scanner.Scan() {
		l := scanner.Text()
		if l != "" && !strings.HasPrefix(l, "#") {
			if ts, err := strconv.ParseFloat(l, 64); err == nil {
				if !firstTimestamp {
					if gte(previousTimestamp, ts) {
						fmt.Printf("%d: not monotone increasing: %s -> %s", currentLineNumber, floatStr(previousTimestamp), floatStr(ts))
						foundErrors = true
						if *outFile != "" {
							fixedTimestamp := previousTimestamp + 1/math.Pow10(*precision)
							fmt.Printf(" changing to %s", floatStr(fixedTimestamp))
							ts = fixedTimestamp
							l = floatStr(ts)
						}
						fmt.Println("")
					}
				}
				firstTimestamp = false
				previousTimestamp = ts
			} else {
				fmt.Printf("%d: invalid data: %s\n", currentLineNumber, l)
				l = "## " + l
			}
		}
		if *outFile != "" {
			_, err := writer.WriteString(l + "\n")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		currentLineNumber++
	}

	if *outFile != "" {
	}

	if foundErrors && *outFile == "" {
		os.Exit(2)
	}
}
