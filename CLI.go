package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
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
	choice          string
	quitting        bool
	tasks           [][]string
	displayingTasks bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.displayingTasks {
				// Exit task display mode and go back to the main menu
				m.displayingTasks = false
				m.choice = ""
				return m, nil
			}

			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
				switch i.action {
				case "add":
					// Here you can add task logic
				case "display":
					return m.displayTasks(), nil
				case "remove":
					return m.removeTasks(), nil
				case "save_exit":
					m.quitting = true
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) displayTasks() model {

	if len(m.tasks) == 0 {
		m.choice = "No tasks available."
	} else {
		var taskStrings []string
		for _, task := range m.tasks {

			taskStrings = append(taskStrings, strings.Join(task, ", "))
		}
		m.choice = "Tasks:\n" + strings.Join(taskStrings, "\n")
	}
	m.displayingTasks = true
	return m
}

func (m model) removeTasks() model {
	if len(m.tasks) == 0 {
		m.choice = "No tasks available."
	} else {
		var taskStrings []string
		for _, task := range m.tasks {
			taskStrings = append(taskStrings, strings.Join(task, ", "))
		}
		m.choice = "Tasks:\n" + strings.Join(taskStrings, "\n")
	}
	m.displayingTasks = true
	return m
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(m.choice + "\n\nPress Enter to go back to the main menu.")
	}

	if m.quitting {
		return quitTextStyle.Render("Bye bye")
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

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{
		list:  l,
		tasks: tasks_list, // Example tasks
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
