package main

import (
	"bufio"
	"fmt"
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
		csvFile, err := os.Create("todo_list.csv")
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer csvFile.Close() // Ensure the file is closed after creation
	} else if err != nil {
		log.Fatalf("error checking file: %s", err)
	}

	read_todo(csvFile)

}

func write_todo() {

}

func read_todo(file *os.File) {

	file, err := os.Open("todo_list.csv")
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

}
