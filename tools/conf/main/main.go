package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/golang/protobuf/proto"
	json_reader "v2ray.com/ext/encoding/json"
	"v2ray.com/ext/tools/conf"
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

	jsonConfig := &conf.Config{}
	decoder := json.NewDecoder(&json_reader.Reader{
		Reader: os.Stdin,
	})

	if err := decoder.Decode(jsonConfig); err != nil {
		os.Stderr.WriteString("failed to read json config: " + err.Error())
		return
	}

	pbConfig, err := jsonConfig.Build()
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
