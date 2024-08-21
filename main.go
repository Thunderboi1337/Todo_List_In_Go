package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type TodoList struct {
	Task     string
	Priority string
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
	task_write(csvFile)

	_, err := csvFile.Seek(0, 0)
	if err != nil {
		log.Fatalf("failed to reset file pointer: %s", err)
	}

	task_get(csvFile)

	_, err = csvFile.Seek(0, 0)
	if err != nil {
		log.Fatalf("failed to reset file pointer: %s", err)
	}

	task_remove(csvFile, 3)

}

// Function that accepts the file as an argument and reads data from file argument
func task_get(file *os.File) {

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++

		fmt.Println(i, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
}

// Function that accepts the file as an argument and writes data to the file argument
func task_write(file *os.File) {

	var task, priority string

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter task: ")
	task, _ = reader.ReadString('\n')
	task = task[:len(task)-1]

	fmt.Println("Enter priority: ")
	priority, _ = reader.ReadString('\n')
	priority = priority[:len(priority)-1]

	todo := TodoList{
		Task:     task,
		Priority: priority,
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{todo.Task, todo.Priority}
	err := writer.Write(record)
	if err != nil {
		log.Fatalf("error writing record to file: %s", err)
	}

}

func task_remove(file *os.File, taskRemove int) error {
	// Read all records from the file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %w", err)
	}

	// Check if taskRemove is within bounds
	if taskRemove < 0 || taskRemove >= len(records) {
		return fmt.Errorf("row index %d out of bounds", taskRemove)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatalf("failed to reset file pointer: %s", err)
	}

	// Remove the specific row
	records = append(records[:taskRemove], records[taskRemove+1:]...)

	// Write the updated records back to the file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return fmt.Errorf("error writing CSV file: %w", err)
	}

	return nil
}
