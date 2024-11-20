package utils

import (
	"bufio"
	"log"
	"os"
)

func ReadInput(filename string) ([]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Initialize a slice to store lines
	var lines []string

	// Use a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from file:", err)
		return nil, err
	}

	// Return the slice of lines
	return lines, nil
}
