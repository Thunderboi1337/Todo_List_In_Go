# Todo List in Go

A simple command-line Todo List application built in Go, utilizing [Bubble Tea](https://github.com/charmbracelet/bubbletea) for an interactive text-based user interface. This project was created as a learning exercise to get familiar with Go and CSV data handling.

## Features

- **Interactive CLI**: Leverages Bubble Tea for an intuitive and responsive terminal user interface.
- **CSV Storage**: Tasks are stored in a CSV file for simplicity and ease of access.
- **Task Management**: Users can add, display, and remove tasks with associated priorities.

## Installation

To get started with the project:

1. **Clone the repository**:
    ```sh
    git clone https://github.com/yourusername/todo_list_go.git
    cd todo_list_go
    ```

2. **Install dependencies**:
    ```sh
    go mod tidy
    ```

3. **Build the project**:
    ```sh
    go build -o todo_list
    ```

4. **Run the application**:
    ```sh
    ./todo_list
    ```

## Usage

Upon running the application, you can interact with it using the following commands:

- **Add tasks**: Enter a task description followed by its priority.
- **Display tasks**: View all current tasks along with their priorities.
- **Remove tasks**: Delete tasks by selecting them from the list.

### Example Usage

Here are some GIFs demonstrating how to use the application:

![Add and Complete Task](https://media.giphy.com/media/example1/giphy.gif)

![View and Delete Task](https://media.giphy.com/media/example2/giphy.gif)

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
