# Project Tableau: Architecture & Design Specification

**Tableau** is a developer-centric, local-first Project Management tool that lives inside your git repository. It follows the Unix philosophy: it is a simple tool that does one thing well, stores data as text, and is composable.

## 1. Core Philosophy
1.  **Git-Native:** State is stored in files. Your project history matches your code history.
2.  **Speed:** No API calls, no loading spinners. Instant interaction via Go.
3.  **Duality:** 
    *   **CLI** for scripting, batch operations, and quick capture.
    *   **TUI** for planning, sorting, and visual management.
4.  **Editor Agnostic:** Task content is Markdown. Use Vim, VS Code, or Zed to write specs.

---

## 2. Feature Specification

### 2.1. The Data Strategy (The "Database")
All state resides in a hidden `.tableau` directory at the root of the project.

**Directory Structure:**
```text
.tableau/
├── config.toml       # Custom columns, colors, git integration settings
├── .meta             # Stores the auto-increment ID counter
└── tasks/            # The actual data
    ├── 1-setup-ci.md
    ├── 2-refactor-auth.md
    └── 3-ui-design.md
```

**Task File Format (`1-setup-ci.md`):**
We use **TOML Frontmatter** for metadata (strictly typed) and **Markdown** for the body (free text).

```markdown
+++
id = 1
title = "Setup CI Pipeline"
status = "doing"      # todo, doing, done (configurable)
priority = "high"     # low, medium, high
labels = ["devops", "backend"]
assignee = "me"
created_at = 2023-11-01T10:00:00Z
updated_at = 2023-11-02T15:30:00Z
+++

## Context
We need to setup GitHub Actions.

## Acceptance Criteria
- [x] Linting passes
- [ ] Tests run on PR
```

### 2.2. The CLI (Command Line Interface)
Designed for speed and scripting.

*   **Creation (Smart Parsing):**
    `tableau new "Refactor login @high +auth"`
    *(Parses `@high` as priority and `+auth` as a label)*
*   **Manipulation:**
    `tableau move 12 done` (Moves ID 12 to Done)
    `tableau assign 12 @josh`
*   **Querying:**
    `tableau list --filter "label:bug status:todo"`
*   **Editing:**
    `tableau edit 12` (Opens `$EDITOR` on the file)

### 2.3. The TUI (Terminal User Interface)
Designed for management and overview.

*   **Kanban View:** Vertical columns based on status.
*   **Navigation:** Vim-bindings (`h/j/k/l`) to navigate cards.
*   **Interactive:**
    *   `Enter`: Open task details (renders Markdown).
    *   `e`: Edit task in `$EDITOR`.
    *   `Space`: Toggle visual selection (for batch operations).
    *   `/`: Fuzzy search tasks.
*   **Drag & Drop:** Keys (e.g., `H/L` or `Ctrl+Left/Right`) move the selected task across columns.

---

## 3. Technical Architecture

We will use a **Hexagonal Architecture (Ports & Adapters)** to decouple the storage format from the UI. This allows us to swap the TUI library or storage engine later if needed, and keeps the CLI/TUI logic clean.

### 3.1. Technology Stack (Go)
*   **Framework:** `cobra` (CLI commands).
*   **TUI Engine:** `bubbletea` (The Elm Architecture).
*   **Styling:** `lipgloss` (CSS-like styling for terminal).
*   **Parsing:** `adrg/frontmatter` (TOML+MD), `goldmark` (MD Rendering).
*   **Fuzzy Search:** `sahilm/fuzzy` (For the search bar).

### 3.2. Package Structure

```text
cmd/
  tableau/           # Main entry point (wires everything together)
    main.go
    root.go          # Cobra root command
    server.go        # TUI command
    task.go          # CLI CRUD commands

internal/
  domain/            # Pure business logic (No imports from ui or store)
    task.go          # Struct definitions (Task, Status, Priority)
    board.go         # Logic for moving tasks, validating transitions
    query.go         # The filter language parser

  adapters/
    repo/            # File System interactions
      fs.go          # Reads/Writes Markdown files
      parser.go      # Handles Frontmatter marshalling
    
    git/             # Git integration
      branch.go      # Checks out branches matching task IDs

  app/               # Application Services (The "Brain")
    service.go       # Orchestrates "Move Task" -> "Update Git" -> "Save File"

  ui/
    tui/             # Bubbletea components
      model.go       # The TUI state
      column.go      # Visual column component
      styles.go      # Lipgloss definitions
    
    renderer/        # CLI text output logic (tables, lists)
```

---

## 4. Advanced "Power" Features (The differentiator)

To make this better than a simple toy, we will implement these specific subsystems:

### 4.1. The Query Language Engine (`internal/domain/query.go`)
Instead of simple string matching, we build a tokenizer that parses a query string into a filter struct.

*   **Input:** `priority:high -label:wontfix sort:created_desc`
*   **Logic:**
    1.  Tokenize string by spaces.
    2.  Identify keys (`priority`), operators (`:`), and values (`high`).
    3.  Apply exclusion (`-` prefix).
    4.  Apply sort logic.

### 4.2. Git Context Integration (`internal/adapters/git`)
Tableau should be aware it runs inside a repo.

*   **Feature:** `tableau start <id>`
    *   Does the task have a `branch` field?
    *   **If yes:** `git checkout <branch>`
    *   **If no:** Create branch `feature/<id>-<slug>`, update task file with branch name, then checkout.
*   **Feature:** **Git Hooks**.
    *   Add a `pre-push` hook: Check if any task referenced in commits is still "In Progress".

### 4.3. Batch Operations
The TUI needs a "Multi-Select" mode.
1.  User presses `v` to enter visual mode.
2.  User presses `j/k` and `space` to select multiple tasks.
3.  User presses `A` (archive) or `D` (done).
4.  **Code Path:** `Service.BatchUpdate(ids, func(t *Task) { t.Status = "Done" })`.

---

## 5. Implementation Roadmap

### Phase 1: The Foundation (CLI + Data)
*   [ ] Set up Go module and directory structure.
*   [ ] Implement `internal/domain` structs.
*   [ ] Implement `internal/adapters/repo` to read/write Markdown/TOML.
*   [ ] Create `tableau init` (scaffold folder) and `tableau new` (create file).
*   [ ] Create `tableau list` (print table).

### Phase 2: The Visuals (TUI MVP)
*   [ ] Initialize Bubbletea model.
*   [ ] Create `Column` component using Lipgloss.
*   [ ] Implement "Load Tasks" into TUI memory.
*   [ ] Implement navigation (Left/Right/Up/Down).
*   [ ] Implement "Move" logic (optimistic UI update + background disk save).

### Phase 3: The Power (Logic & Integration)
*   [ ] Implement the **Query Engine** (filtering).
*   [ ] Implement `tableau start` (Git branching).
*   [ ] Implement "Edit" (shell out to `vim`/`nano`).
*   [ ] Add `config.toml` support for custom columns.

### Phase 4: Polish
*   [ ] Add "Blocked" logic (visual lock icon).
*   [ ] Add support for rendering Markdown tables/lists inside the TUI detail view.

---

## 6. Example Configuration (`config.toml`)

This allows the tool to be adaptable to different workflows.

```toml
[project]
name = "Tableau"

[columns]
# Map internal status keys to display titles
todo = "To Do"
doing = "In Progress"
review = "Code Review"
done = "Done"

[style]
border_color = "#7D56F4"
selected_color = "#ff00ff"

[git]
auto_branch = true
branch_prefix = "feat/"
```
