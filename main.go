package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var help bool
	var inputPath, outputPath string

	flag.BoolVar(&help, "h", false, "cara menggunakannya")
	flag.BoolVar(&help, "help", false, "cara menggunakannya")

	flag.StringVar(&inputPath, "i", "", "Lokasi file JSON sebagai input")
	flag.StringVar(&inputPath, "input", "", "Lokasi file JSON sebagai input")
	flag.StringVar(&outputPath, "o", "", "Output file JSON sebagai output")
	flag.StringVar(&outputPath, "output", "", "Output file JSON sebagai output")

	flag.Parse()

	if help || inputPath == "" || outputPath == "" {
		printUsage()
		os.Exit(0)
	}
}

func printUsage() {
	fmt.Println("Usage: mockdata [-i | --input] <input file> [-o | --output] <output file>")
	fmt.Println("-i --input: File input berupa JSON sebagai input")
	fmt.Println("-o --output: File output JSON sebagai output")
}
