package main

import (
	"bufio"
	"log"
	"os"
)

type Employee struct {
	ID  string
	Age int
}

func main() {
	var csvFile *os.File

	if _, err := os.Stat("todo_list.csv"); os.IsNotExist(err) {
		// File does not exist, create it
		csvFile, err = os.Create("todo_list.csv")
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer csvFile.Close() // Ensure the file is closed after creation
	} else if err != nil {
		log.Fatalf("error checking file: %s", err)
	} else {
		// File exists, open it
		csvFile, err = os.Open("todo_list.csv")
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		defer csvFile.Close() // Ensure the file is closed after reading
	}

	// Pass the file to the read_todo function
	read_todo(csvFile)
}

// Function that accepts the file as an argument
func read_todo(file *os.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Print each line in the file
		log.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
}
