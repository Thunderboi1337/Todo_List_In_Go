package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14
const defaultWidth = 20

type item struct {
	title  string
	action string
}

type model struct {
	list            list.Model
	taskInput       textinput.Model
	prioInput       textinput.Model
	table           table.Model
	displayList     list.Model
	removeTable     table.Model
	removeList      list.Model
	tasks           [][]string
	choice          string
	addingTasks     bool
	displayingTasks bool
	removingTasks   bool
	quitting        bool
}

var (
	itemStyle           = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	quitTextStyle       = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	addInstructionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).PaddingTop(1)
	addHeadingStyle     = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("62")).
				MarginLeft(2).
				MarginTop(1).
				PaddingLeft(1).
				PaddingRight(1).
				Align(lipgloss.Left)
)

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (m model) startAddTasks() model {
	m.addingTasks = true
	m.taskInput.Placeholder = "Enter task"
	m.taskInput.Focus()
	return m
}

func (m model) startDisplayTasks() model {
	if len(m.tasks) == 0 {
		m.choice = "No tasks available."
	} else {
		columns := []table.Column{
			{Title: "ID", Width: 5},
			{Title: "Task", Width: 60},
			{Title: "Priority", Width: 25},
		}

		var rows []table.Row

		for id, taskList := range m.tasks {

			row := []string{
				fmt.Sprintf("%d", id+1), // Task ID
				taskList[0],             // Task Name
				taskList[1],             // Priority
			}

			rows = append(rows, row)
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
		)

		m.table = t
		m.displayingTasks = true
		m.choice = ""
	}

	return m
}

func (m model) startRemoveTasks() model {
	if len(m.tasks) == 0 {
		m.choice = "No tasks available."
	} else {
		columns := []table.Column{
			{Title: "ID", Width: 5},
			{Title: "Task", Width: 60},
			{Title: "Priority", Width: 25},
		}

		var rows []table.Row

		for id, taskList := range m.tasks {
			row := []string{
				fmt.Sprintf("%d", id+1), // Task ID
				taskList[0],             // Task Name
				taskList[1],             // Priority
			}
			rows = append(rows, row)
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
		)

		m.removeTable = t
		m.removingTasks = true
		m.choice = ""
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.removeList.SetWidth(msg.Width)
		m.displayList.SetWidth(msg.Width)
		m.table.SetWidth(msg.Width)
		m.removeTable.SetWidth(msg.Width) // Add this line
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			m.addingTasks = false
			m.displayingTasks = false
			m.removingTasks = false
			m.taskInput.SetValue("")
			m.prioInput.SetValue("")
			return m, nil

		case "enter":
			if m.addingTasks {
				if m.taskInput.Focused() {
					task := m.taskInput.Value()
					if task != "" {
						m.taskInput.Blur()
						m.prioInput.Focus()
						return m, nil
					} else {
						m.addingTasks = false
						return m, nil
					}

				} else if m.prioInput.Focused() {
					prio := m.prioInput.Value()
					if prio != "" {
						m.tasks = append(m.tasks, []string{m.taskInput.Value(), prio})
						m.prioInput.Blur()
						m.taskInput.SetValue("")
						m.prioInput.SetValue("")
						m.addingTasks = false
						return m, nil
					} else {
						m.tasks = append(m.tasks, []string{m.taskInput.Value(), prio})
						m.taskInput.SetValue("")
						m.prioInput.SetValue("")
						m.addingTasks = false
						return m, nil
					}
				}
			} else if m.removingTasks {
				// Remove the selected task from the table

				index := m.removeList.Index()
				if index >= 0 && index < len(m.tasks) {
					m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)

					// Update the removeList to reflect the changes
					var items []list.Item
					for _, task := range m.tasks {
						items = append(items, item{title: strings.Join(task, ", "), action: "remove"})
					}
					m.removeList.SetItems(items)
				}

				m.removingTasks = false
				m.choice = ""
				return m, nil
			} else if m.displayingTasks {
				// Exit task display mode and go back to the main menu
				m.displayingTasks = false
				m.choice = ""
				return m, nil
			}

			// Handle the main menu selection
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
				switch i.action {
				case "add":
					return m.startAddTasks(), nil
				case "display":
					return m.startDisplayTasks(), nil
				case "remove":
					return m.startRemoveTasks(), nil
				case "save_exit":
					task_save(m.tasks)
					m.quitting = true
					return m, tea.Quit
				}
			}
		}
	}

	// Handle updates based on the current state
	if m.addingTasks {
		if m.taskInput.Focused() {
			m.taskInput, cmd = m.taskInput.Update(msg)
		} else if m.prioInput.Focused() {
			m.prioInput, cmd = m.prioInput.Update(msg)
		}
	} else if m.removingTasks {
		m.removeTable, cmd = m.removeTable.Update(msg) // Update this line
	} else if m.displayingTasks {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.addingTasks {
		if m.taskInput.Focused() {
			return fmt.Sprintf(
				"%s\n\n%s\n\n%s",
				addHeadingStyle.Render("ADD_TASK"),
				m.taskInput.View(),
				addInstructionStyle.Render("(esc to cancel, enter to next)"),
			) + "\n"
		} else if m.prioInput.Focused() {
			return fmt.Sprintf(
				"%s\n\n%s\n\n%s",
				addHeadingStyle.Render("ADD_TASK"),
				m.prioInput.View(),
				addInstructionStyle.Render("(esc to cancel, enter to add)"),
			) + "\n"
		}
	}

	if m.quitting {
		task_save(m.tasks)
		return quitTextStyle.Render("")
	}

	if m.displayingTasks {
		return fmt.Sprintf(
			"%s\n\n%s%s",
			addHeadingStyle.Render("DISPLAY_TASKS"),
			m.table.View(),
			addInstructionStyle.Render("(esc to return"),
		) + "\n"
	}

	if m.removingTasks {
		return fmt.Sprintf(
			"%s\n\n%s%s",
			addHeadingStyle.Render("REMOVE_TASKS"),
			m.removeTable.View(),
			addInstructionStyle.Render("(esc to return"),
		) + "\n"

	}

	return "\n" + m.list.View()
}

func CLI(tasks_list [][]string) {
	items := []list.Item{
		item{title: "Add tasks", action: "add"},
		item{title: "Display tasks", action: "display"},
		item{title: "Remove tasks", action: "remove"},
		item{title: "Save & Exit", action: "save_exit"},
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Gora"
	l.KeyMap.Quit.SetHelp("ctrl+c", "quit")
	l.KeyMap.Filter.Unbind()

	taskInput := textinput.New()
	taskInput.Placeholder = "Enter task"
	taskInput.Focus()
	taskInput.CharLimit = 156
	taskInput.Width = 20

	prioInput := textinput.New()
	prioInput.Placeholder = "Enter priority"
	prioInput.Focus()
	prioInput.CharLimit = 156
	prioInput.Width = 20

	displayList := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
	displayList.Title = "DISPLAY_TASKS"

	removeList := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
	removeList.Title = "REMOVE_TASKS"

	m := model{
		list:        l,
		displayList: displayList,
		removeList:  removeList,
		taskInput:   taskInput,
		prioInput:   prioInput,
		tasks:       tasks_list,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// Function that accepts the file as an argument and reads data from file argument
func task_collect(file *os.File) [][]string {
	reader := csv.NewReader(file)

	// Read all the CSV data into a slice of tasks
	tasks, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return nil
	}

	return tasks
}

func task_save(tasks [][]string) {
	var todoFile = getTodoFilePath()
	csvFile, err := os.OpenFile(todoFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	err = writer.WriteAll(tasks)
	if err != nil {
		fmt.Printf("error writing CSV file: %v", err)
	}
}

func getTodoFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// Build the path to the todo list file
	return filepath.Join(homeDir, ".local", "share", "Gora", "Gora_list.csv")
}

func main() {

	var csvFile *os.File

	var todoFile = getTodoFilePath()

	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		// File does not exist, create it
		csvFile, err = os.Create(todoFile)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer csvFile.Close() // Ensure the file is closed after creation
	} else if err != nil {
		log.Fatalf("error checking file: %s", err)
	} else {
		// File exists, open it
		csvFile, err = os.OpenFile(todoFile, os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		// Ensure the file is closed after reading
		defer csvFile.Close()
	}

	// Pass the file to the read_todo function
	tasks := task_collect(csvFile)

	CLI(tasks)

}
