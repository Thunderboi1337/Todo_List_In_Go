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
		csvFile, err = os.Create("todo_list.csv")
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer csvFile.Close() // Ensure the file is closed after creation
	} else if err != nil {
		log.Fatalf("error checking file: %s", err)
	} else {
		// File exists, open it
		csvFile, err = os.OpenFile("todo_list.csv", os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		defer csvFile.Close() // Ensure the file is closed after reading
		// defer a cool feature in Go, very nois
	}

	// Pass the file to the read_todo function
	read_todo(csvFile)
	write_todo(csvFile)
	read_todo(csvFile)

}

// Function that accepts the file as an argument and reads data from file argument
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

// Function that accepts the file as an argument and writes data to the file argument
func write_todo(file *os.File) {

	fmt.Println("Enter something random: ")

	// var then variable name then variable type
	var first string

	// Taking input from user
	fmt.Scanln(&first)
	fmt.Println("Enter something random again: ")
	var second string
	fmt.Scanln(&second)

	longstrng := first
	longstrng += ", " + second

	fmt.Print(longstrng)

	writer := bufio.NewWriter(file)
	da, err := writer.WriteString(longstrng)
	if err != nil {
		log.Fatalf("error write file: %d", da)
	}

	err = writer.Flush()
	if err != nil {
		log.Fatalf("error flushing buffer: %s", err)
	}

}
