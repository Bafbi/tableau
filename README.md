# Tableau

Tableau is a developer-centric, local-first Project Management tool that lives inside your git repository. It follows the Unix philosophy: it is a simple tool that does one thing well, stores data as text (Markdown), and is composable.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)

## Features

- **Git-Native**: Tasks are stored as Markdown files in a `.tableau` directory. Your project history matches your code history.
- **Local-First**: No API calls, no loading spinners. Instant interaction.
- **CLI & TUI**: 
  - **CLI** for scripting and quick capture.
  - **TUI** (Terminal User Interface) for visual planning (Kanban board).
- **Git Integration**: Automatically create and checkout feature branches linked to tasks.
- **Configurable**: Customize columns, colors, and workflow rules via `config.toml`.

## Installation

### From Source

Ensure you have [Go](https://go.dev/) installed (version 1.22 or later).

```bash
git clone https://github.com/bafbi/tableau.git
cd tableau
go install ./cmd/tableau
```

## Usage

### Initialization

Start by initializing Tableau in your project root:

```bash
tableau init
```

This creates a `.tableau` directory to store your tasks and configuration.

### Managing Tasks

Create a new task:
```bash
tableau new "Refactor authentication middleware"
```

List tasks:
```bash
tableau list
# Filter tasks
tableau list "status:todo priority:high"
```

Edit a task (opens in your `$EDITOR`):
```bash
tableau edit 1
```

### The Kanban Board (TUI)

Launch the interactive board:

```bash
tableau board
```

- **Navigation**: Use `h`, `j`, `k`, `l` (Vim style) or arrow keys.
- **Move Tasks**: Use `H` (left) and `L` (right) to move tasks between columns.
- **Details**: Press `Enter` to view task details and comments.
- **Quit**: Press `q` or `Esc`.

### Git Workflow

Start working on a task (creates a branch `feat/1-refactor...` and moves task to "Doing"):

```bash
tableau start 1
```

### Collaboration

Block a task:
```bash
tableau block 1
```

Add a comment:
```bash
tableau comment 1 "I'm looking into the JWT library now."
```

## Configuration

You can customize Tableau by creating a `.tableau/config.toml` file:

```toml
[project]
name = "My Project"

[columns]
todo = "Backlog"
doing = "In Progress"
review = "Code Review"
done = "Shipped"

[style]
border_color = "#7D56F4"
selected_color = "#ff00ff"

[git]
branch_prefix = "feat/"
```

## Contributing

We use [mise](https://mise.jdx.dev/) for development workflow management.

1.  Install `mise`.
2.  Clone the repository.
3.  Run setup:
    ```bash
    mise install
    ```

### Common Commands

- **Build**: `mise run build`
- **Test**: `mise run test`
- **Lint**: `mise run lint`
- **Dev Environment**: `mise run dev -- board` (Runs the app in a sandboxed `.tableau_dev` directory)

## License

MIT
