package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	title  string
	action string
}

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

type model struct {
	list            list.Model
	addList         textinput.Model
	displayList     list.Model
	removeList      list.Model
	textInput       textinput.Model
	tasks           [][]string
	choice          string
	addingTasks     bool
	displayingTasks bool
	removingTasks   bool
	quitting        bool
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
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.addingTasks {
				// Add the task and return to the main menu
				task := m.textInput.Value()
				if task != "" {
					m.tasks = append(m.tasks, []string{task})
				}
				m.addingTasks = false
				m.textInput.SetValue("") // Clear input after adding
				return m, nil
			}

			if m.removingTasks {
				// Remove the selected task
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
			}
			if m.displayingTasks || m.removingTasks {
				// Exit task display mode or remove mode and go back to the main menu
				m.displayingTasks = false
				m.removingTasks = false
				m.choice = ""
				return m, nil
			}

			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
				switch i.action {
				case "add":
					m.addingTasks = true
					m.textInput.Focus()
					return m, nil
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

	if m.addingTasks {
		m.textInput, cmd = m.textInput.Update(msg)
	} else if m.removingTasks {
		m.removeList, cmd = m.removeList.Update(msg)
	} else if m.displayingTasks {
		m.displayList, cmd = m.displayList.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) startDisplayTasks() model {
	if len(m.tasks) == 0 {
		m.choice = "No tasks available."
	} else {
		var items []list.Item
		for _, task := range m.tasks {
			items = append(items, item{title: strings.Join(task, ", "), action: "display"})
		}
		m.displayList.SetItems(items)
		m.displayingTasks = true
		m.choice = ""
	}
	return m
}

func (m model) startRemoveTasks() model {
	if len(m.tasks) == 0 {
		m.choice = "No tasks available."
	} else {
		var items []list.Item
		for _, task := range m.tasks {
			items = append(items, item{title: strings.Join(task, ", "), action: "remove"})
		}
		m.removeList.SetItems(items)
		m.removingTasks = true
		m.choice = ""
	}
	return m
}
func (m model) View() string {

	if m.addingTasks {
		return fmt.Sprintf(
			"Enter task:\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to cancel, enter to add)",
		) + "\n"
	}

	if m.quitting {
		task_save(m.tasks)
		return quitTextStyle.Render("Bye bye")
	}

	if m.removingTasks {
		return "\n" + m.removeList.View()
	}

	if m.displayingTasks {
		return "\n" + m.displayList.View()
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

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "TODO_LIST"

	addList := textinput.New()
	addList.Placeholder = "ADD_TASKS"
	addList.Focus()
	addList.CharLimit = 156
	addList.Width = 20

	textInput := textinput.New()
	textInput.Placeholder = "Enter task"
	textInput.Focus()
	textInput.CharLimit = 156
	textInput.Width = 20

	displayList := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
	displayList.Title = "DISPLAY_TASKS"

	removeList := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
	removeList.Title = "REMOVE_TASKS"

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{
		list:        l,
		addList:     addList,
		displayList: displayList,
		removeList:  removeList,
		textInput:   textInput,
		tasks:       tasks_list,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
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
		defer csvFile.Close()
	}

	// Pass the file to the read_todo function
	tasks := task_collect(csvFile)

	CLI(tasks)
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
	csvFile, err := os.OpenFile("todo_list.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
