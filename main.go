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

	"github.com/muhamadnurfi/mockdata/data"
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

	if err := validateType(mapping); err != nil {
		fmt.Printf("Error validating mapping: %v\n", err)
		os.Exit(0)
	}

	result, err := generatingOutput(mapping)
	if err != nil {
		fmt.Printf("Error generating output: %v\n", err)
		os.Exit(0)
	}

	if err := writeOutput(outputPath, result); err != nil {
		fmt.Printf("Error write output: %v\n", err)
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
		return errors.New("no input file specified")
	}

	if mapping == nil {
		return errors.New("no input mapping specified")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	fileByte, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(fileByte) == 0 {
		return errors.New("no input file specified")
	}

	if err := json.Unmarshal(fileByte, &mapping); err != nil {
		return err
	}

	return nil
}

func validateType(mapping map[string]string) error {
	for _, value := range mapping {
		if !data.Supported[value] {
			return errors.New("type data not supported")
		}
	}

	return nil
}

func generatingOutput(mapping map[string]string) (map[string]any, error) {
	result := make(map[string]any)

	for key, dataType := range mapping {
		result[key] = data.Generate(dataType)
	}

	return result, nil
}

func writeOutput(path string, result map[string]any) error {
	if path == "" {
		return errors.New("path invalid")
	}

	flags := os.O_RDWR | os.O_CREATE | os.O_TRUNC // operasi | itu adalah Operasi Inclusive. Pada dasarnya os itu adalah integer yang nanti diubah menjadi binary
	file, err := os.OpenFile(path, flags, 0644)   // 0644 adalah pengaturan hak akses (file permission) pada sistem operasi berbasis Unix/Linux
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	resultByte, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return err
	}

	if _, err := file.Write(resultByte); err != nil {
		return err
	}

	return nil
}
