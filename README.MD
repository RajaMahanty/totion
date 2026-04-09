# Totion

A modern terminal-based note-taking application built with Go. Totion provides a fast and intuitive interface for managing your notes directly from the command line.

## Features

- **Terminal User Interface (TUI)** - Built with Bubble Tea for a smooth, responsive CLI experience
- **Note Management** - Create, edit, and organize notes efficiently
- **File Listing** - Browse and manage your notes with a beautiful list interface
- **Text Editing** - Full-featured text area for note editing
- **Quick Navigation** - Keyboard shortcuts for fast navigation and actions

## Prerequisites

- Go 1.25.5 or higher

## Installation

Clone the repository and build from source:

```bash
git clone https://github.com/RajaMahanty/totion
cd totion
make build
```

Or simply run:

```bash
make run
```

This will build and execute the application directly.

## Usage

Start the application:

```bash
./totion
```

### Keyboard Shortcuts

- `q` / `Ctrl+C` - Quit the application
- `Esc` - Go back / Cancel current action
- Navigate through menus using arrow keys
- Create new notes with the input dialog

## Project Structure

- `main.go` - Main application logic and UI components
- `go.mod` - Go module dependencies
- `makefile` - Build automation

## Dependencies

- **Bubble Tea** - TUI framework for building interactive CLI applications
- **Lipgloss** - Styling and layout library for terminal UIs
- **Bubbles** - Collection of components for Bubble Tea (list, textarea, textinput)

## Development

The application uses a model-view-update architecture (Elm-style) provided by Bubble Tea for managing application state and rendering.

### Building

```bash
make build
```

### Running

```bash
make run
```

## Author

[Raja Mahanty](https://github.com/RajaMahanty)

## License

See LICENSE file for details.
