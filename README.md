# Todo List in Go
****Working progress****
- A simple command-line Todo List application built in Go, utilizing [Bubble Tea](https://github.com/charmbracelet/bubbletea) for an interactive text-based user interface. This project was created as a learning exercise to get familiar with Go and CSV data handling.

## Features

- **Interactive CLI**: Leverages Bubble Tea for an intuitive and responsive terminal user interface.
- **CSV Storage**: Tasks are stored in a CSV file for simplicity and ease of access.
- **Task Management**: Users can add, display, and remove tasks with associated priorities.

## Installation

To get started with the project:

1. **Clone the repository**:
    ```sh
    git clone https://github.com/thunderboi1337/Todo_List_In_Go.git
    cd Todo_List_In_Go
    ```

2. **Install dependencies**:
    ```sh
    go mod tidy
    ```

3. **Build the project**:
    ```sh
    go build main.go
    ```

4. **Run the application**:
    ```sh
    ./main
    ```

## Usage

Upon running the application, you can interact with it using the following commands:

- **Add tasks**: Enter a task description followed by its priority.
- **Display tasks**: View all current tasks along with their priorities.
- **Remove tasks**: Delete tasks by selecting them from the list.

### Example Usage

Here are some GIFs demonstrating how to use the application:
## Adding task
![Add and Complete Task](screenshots/Todolist.GIF)

## Display task
![Display Task](screenshots/TodolistDisplay.GIF)

## Remove task
![View and Delete Task](screenshots/TodolistRemove.GIF)

## Code Structure

- The application reads and writes tasks to a CSV file (`csv_data/todo_list.csv`).
- Tasks are managed in-memory during the session and saved upon exiting the program.
- The UI is divided into different states: adding tasks, displaying tasks, and removing tasks.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the text-based user interface.
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - For styling the terminal UI.

## License

This project is open source, and you're free to use, modify, and distribute it as you see fit.

## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, feel free to open an issue or submit a pull request. A link to the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework is included for further exploration.
