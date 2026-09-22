package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var help bool
	var inputPath, outputPath string

	flag.StringVar(&inputPath, "i", "", "Path to input file")
	flag.StringVar(&inputPath, "input", "", "Path to input file")
	flag.StringVar(&outputPath, "o", "", "Path to output file")
	flag.StringVar(&outputPath, "output", "", "Path to output file")

	flag.Parse()

	if help || inputPath == "" || outputPath == "" {
		printUsage()
		os.Exit(0)
	}

	if err := validatePathInput(inputPath); err != nil {
		fmt.Printf("Error validating input: %v\n", err)
		os.Exit(0)
	}

	if err := validatePathOutput(outputPath); err != nil {
		fmt.Printf("Error validating output: %v\n", err)
		os.Exit(0)
	}

	var mapping map[string]string
	if err := readInput(inputPath, &mapping); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(0)
	}

}

func printUsage() {
	fmt.Printf("Usage: mockdata -i <input file> -o <output file>\n")
	fmt.Printf("-i --input: File to read from input\n")
	fmt.Printf("-o --output: File to write to output\n")
}

func validatePathInput(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	return nil
}

func validatePathOutput(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	fmt.Println("output path:", path)
	confirmOverwrite()
	return nil
}

func confirmOverwrite() {
	fmt.Printf("Are you sure you want to overwrite? [y/N]")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" && response != "ya" {
		fmt.Println("Aborting...")
		os.Exit(0)
	}
}

func readInput(path string, mapping *map[string]string) error {
	if path == "" {
		return errors.New("No input file specified")
	}

	if mapping == nil {
		return errors.New("No mapping spectified input file")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileByte, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(fileByte) == 0 {
		return errors.New("Empty input file")
	}

	if err := json.Unmarshal(fileByte, &mapping); err != nil {
		return err
	}

	return nil
}
