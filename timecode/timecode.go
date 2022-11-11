package timecode

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func almostEqual(a, b float64, precision int) bool {
	return (1 / math.Pow10(precision+1)) > math.Abs(a-b)
}

func gte(a, b float64, precision int) bool {
	if !almostEqual(a, b, precision) {
		return a > b
	}
	return true
}

func floatStr(a float64, precision int) string {
	return strconv.FormatFloat(a, 'f', precision, 64)
}

func processTimeCodeFile(in, out *os.File, precision int) (valid bool, err error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	currentLineNumber := 0
	var previousTimestamp float64
	firstTimestamp := true

	valid = true

	var writer *bufio.Writer

	if out != nil {
		writer = bufio.NewWriter(out)
		defer func(writer *bufio.Writer) {
			_ = writer.Flush()
		}(writer)
	}

	for scanner.Scan() {
		l := scanner.Text()
		if l != "" && !strings.HasPrefix(l, "#") {
			if ts, err := strconv.ParseFloat(l, 64); err == nil {
				if !firstTimestamp {
					if gte(previousTimestamp, ts, precision) {
						fmt.Printf("%d: not monotone increasing: %s -> %s", currentLineNumber, floatStr(previousTimestamp, precision), floatStr(ts, precision))
						valid = false
						if out != nil {
							fixedTimestamp := previousTimestamp + 1/math.Pow10(precision)
							fmt.Printf(" changing to %s", floatStr(fixedTimestamp, precision))
							ts = fixedTimestamp
							l = floatStr(ts, precision)
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
		if out != nil {
			_, err := writer.WriteString(l + "\n")
			if err != nil {
				return false, nil
			}
		}

		currentLineNumber++
	}
	return
}

func CheckTimeCodeFile(in *os.File, precision int) (bool, error) {
	return processTimeCodeFile(in, nil, precision)
}

func CleanTimecodeFile(in, out *os.File, precision int) (bool, error) {
	return processTimeCodeFile(in, out, precision)
}
