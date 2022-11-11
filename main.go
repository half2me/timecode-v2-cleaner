package main

import (
	"flag"
	"github.com/half2me/timecode-v2-cleaner/timecode"
	"os"
)
import "fmt"

var precision = flag.Int("precision", 4, "max number of decimal places allowed in timecode")
var outFile = flag.String("out", "", "Output file with fixed timecode")

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

	in, err := os.Open(flag.Arg(0))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(in)

	if *outFile != "" {
		out, err := os.Create(*outFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(out)

		_, err = timecode.CleanTimecodeFile(in, out, *precision)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		valid, err := timecode.CheckTimeCodeFile(in, *precision)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !valid {
			os.Exit(2)
		}
	}
}
