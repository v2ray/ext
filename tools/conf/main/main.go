package main

import (
	"flag"
	"os"

	"github.com/golang/protobuf/proto"
	"v2ray.com/ext/tools/conf/serial"
)

/*
var (
	inputFile = flag.String("input", "", "Path to the input file. StdIn if empty.")
	inputFormat = flag.String("iformat", "json", "Format of input file.")
	outputFile = flag.String("output", "", "Path to the output file. StdOut if empty.")
	outputFormat = flag.String("oformat", "pb", "Format of output file.")
)
*/

func main() {
	flag.Parse()

	pbConfig, err := serial.LoadJSONConfig(os.Stdin)
	if err != nil {
		os.Stderr.WriteString("failed to parse json config: " + err.Error())
		return
	}

	bytesConfig, err := proto.Marshal(pbConfig)
	if err != nil {
		os.Stderr.WriteString("failed to marshal proto config: " + err.Error())
		return
	}

	if _, err := os.Stdout.Write(bytesConfig); err != nil {
		os.Stderr.WriteString("failed to write proto config: " + err.Error())
		return
	}
}
