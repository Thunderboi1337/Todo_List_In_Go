package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
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
		csvFile, err = os.OpenFile("todo_list.csv", os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		// Ensure the file is closed after reading
		// defer a cool feature in Go, very nois
	}

	// Pass the file to the read_todo function
	tasks := task_collect(csvFile)
	csvFile.Close()

	CLI(tasks)
	/*
		exit := false

			for !exit {
				display_menu()

				option := 0
				_, err := fmt.Scanf("%d", &option)
				if err != nil {
					fmt.Println("Invalid input. Please enter a number.")
					// Clear input buffer
					var discard string
					fmt.Scanln(&discard)
					continue
				}

				switch option {
				case 1:
					tasks = task_write(tasks)
				case 2:
					task_display(tasks)
				case 3:

					tasks = task_remove(tasks)
				case 4:
					task_save(tasks)
					exit = true
				default:
					fmt.Println("Invalid option. Please choose a valid menu item.")
				}
			} */

}

// Function that accepts the file as an argument and reads data from file argument
func task_display(tasks [][]string) {
	num := 0
	for _, record := range tasks {
		num++
		fmt.Println(num, record)

	}

}

// Function that accepts the file as an argument and writes data to the file argument
func task_write(tasks [][]string) [][]string {

	reader := bufio.NewReader(os.Stdin)

	for {
		var task, priority string

		// Get task input from the user
		fmt.Println("Enter task: ")
		task, _ = reader.ReadString('\n')
		task = strings.TrimSpace(task)

		// Get priority input from the user
		fmt.Println("Enter priority: ")
		priority, _ = reader.ReadString('\n')
		priority = strings.TrimSpace(priority)

		// Convert the struct fields to a slice of strings
		record := []string{task, priority}
		tasks = append(tasks, record)

		// Ask if the user wants to add another task
		fmt.Println("Do you want to add another task? (yes/no)")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		if strings.ToLower(answer) != "yes" {
			break
		}
	}

	return tasks

}

func task_remove(tasks [][]string, tasks_remove int) [][]string {

	fmt.Println("Which task do you want to remove?")
	task_display(tasks)
	fmt.Println(len(tasks)+1, "to go back")
	_, err := fmt.Scanf("%d", &tasks_remove)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		// Clear input buffer
		var discard string
		fmt.Scanln(&discard)

	}

	if tasks_remove < 0 || tasks_remove >= len(tasks) {
		fmt.Println("Invalid index")
		return tasks
	}

	return append(tasks[:tasks_remove], tasks[tasks_remove+1:]...)
}

func task_collect(file *os.File) [][]string {

	reader := csv.NewReader(file)

	// Read all the CSV data into a slice of tasks
	tasks, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
	}

	return tasks
}

func task_save(tasks [][]string) {

	csvFile, err := os.OpenFile("todo_list.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	err = writer.WriteAll(tasks)
	if err != nil {
		fmt.Println("error writing CSV file: %w", err)
	}

}

func display_menu() {

	fmt.Println("----------TODO_LIST-----------")
	fmt.Println("Options:----------------------")
	fmt.Println("1. Add tasks------------------")
	fmt.Println("2. Display tasks--------------")
	fmt.Println("3. Remove tasks---------------")
	fmt.Println("4. Save & Exit----------------")
	fmt.Println("------------------------------")

}

func task_get() [][]string {

	return task
}
