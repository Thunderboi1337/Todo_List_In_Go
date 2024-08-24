# Todo List in Go - Gora

**A Work in Progress**

This command-line Todo List application is built in Go and uses [Bubble Tea](https://github.com/charmbracelet/bubbletea) for an interactive and responsive text-based user interface. This project has been a valuable learning exercise, helping me gain hands-on experience with Go, CSV data handling, and application installation.

## Features

- **Interactive CLI**: Provides an intuitive terminal interface with Bubble Tea.
- **CSV Storage**: Tasks are managed using a CSV file for ease of use.
- **Task Management**: Users can add, display, and remove tasks with priorities.

## Installation

Follow these steps to set up the Todo List application on your system:

### Using Makefile

1. **Clone the repository**:
    ```sh
    git clone https://github.com/thunderboi1337/Todo_List_In_Go.git
    cd Todo_List_In_Go
    ```
    This command clones the project repository to your local machine and navigates into the project directory.

2. **Build the project**:
    ```sh
    make build
    ```
    This command compiles the Go source code into an executable binary named `Todo_List`. The `make build` step also runs `go mod tidy` to ensure all dependencies are properly handled.

3. **Install the application**:
    ```sh
    make install
    ```
    This command installs the application:
    - Moves the compiled binary to `/usr/local/bin`, making it available system-wide.
    - Creates necessary directories for configuration, data, and cache (`~/.config/Todo_List`, `~/.local/share/Todo_List`, and `~/.cache/Todo_List`).

4. **Create a desktop entry**:
    ```sh
    make desktop-entry
    ```
    This command creates a desktop entry for the application:
    - Adds a `.desktop` file to `/usr/share/applications/`, which integrates the application into the desktop environment's application menu.
    - This file contains metadata about the application, including its name, comment, executable path, icon path, and terminal usage.

### Without Makefile

If you prefer not to use the Makefile, you can follow these steps manually:

1. **Clone the repository**:
    ```sh
    git clone https://github.com/thunderboi1337/Todo_List_In_Go.git
    cd Todo_List_In_Go
    ```
    Clone the repository and navigate to the project directory.

2. **Install dependencies**:
    ```sh
    go mod tidy
    ```
    This command ensures that all required dependencies for the project are installed and properly managed.

3. **Build the project**:
    ```sh
    go build -o Gora main.go
    ```
    This command compiles the Go code into an executable named `Gora`.

4. **Run the application**:
    ```sh
    ./Gora
    ```
    This command runs the compiled application directly from the terminal.

## Usage

Once the application is running, you can use the following commands:

- **Add tasks**: Input a task description and priority to add new tasks.
- **Display tasks**: View all tasks and their priorities.
- **Remove tasks**: Select and delete tasks from the list.

### Example Usage

Here are some GIFs demonstrating how to use the application:

- **Adding a task**:  
  ![Add and Complete Task](screenshots/Todolist.GIF)

- **Displaying tasks**:  
  ![Display Task](screenshots/TodolistDisplay.GIF)

- **Removing a task**:  
  ![View and Delete Task](screenshots/TodolistRemove.GIF)

## Code Structure

- **CSV File**: The application reads from and writes to a CSV file located at `~/.local/share/Todo_List/todo_list.csv`.
- **In-Memory Management**: Tasks are managed in memory during the session and saved upon exiting.
- **UI States**: The user interface supports different states: adding, displaying, and removing tasks.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the text-based user interface.
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - For styling the terminal UI.

## License

Feel free to use and modify this project as you see fit. Enjoy!
