# AI Code Review Agent CLI

**AI Code Review Agent CLI** is a tool that leverages Google Gemini AI to automatically analyze Git changes and provide professional review reports for code quality improvement.

This project aims to streamline the code review process by proactively identifying potential bugs, security vulnerabilities, and performance issues in developer code, and providing concrete improvement suggestions in Korean (or English).

## Tech Stack

*   **Language**: Go 1.25.4
*   **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
*   **AI Model**: Google Gemini 2.5 Flash (via [Google GenAI SDK](https://github.com/googleapis/go-genai))
*   **Configuration**: JSON based local config

## Getting Started

### Prerequisites

*   **Go**: Version 1.25 or higher ([Download](https://go.dev/dl/))
*   **Git**: Required for project version control and tracking changes.
*   **Google AI API Key**: API key required for using the Gemini model. ([Get API Key](https://aistudio.google.com/app/apikey))

### Installation

Clone the project and install dependencies.

```bash
# Clone the repository
git clone https://github.com/your-username/code-reviewer.git
cd code-reviewer

# Download dependencies
go mod download
```

### Build & Run

Build and run the project.

```bash
# Build the binary
go build -o code-reviewer main.go

# Run the help command to verify installation
./code-reviewer --help
```

On the first run, you will be prompted to enter your Google AI API Key.

## Usage
```bash
# Run code review (default behavior)
./code-reviewer review

# Run code review explicitly
./code-reviewer review report

# Run code review on staged changes
./code-reviewer review --staged

# Generate auto-fixes for code issues
./code-reviewer review fix
```

## Key Features

*   **AI-Powered Code Review**: Uses the Google Gemini 2.5 Flash model to deeply analyze code changes (`git diff`).
*   **Korean Reports**: All analysis results and improvement suggestions are provided in Korean for easy understanding.
*   **Automatic Grading**: Changes are automatically graded into 4 levels: **Good**, **Not Bad**, **Bad**, **Need Check** to quickly identify importance.
*   **Concrete Improvement Suggestions**: Beyond simple pointers, it provides immediately applicable code snippets and refactoring guides for areas needing correction.
*   **Git Integration**: Automatically detects Staged and Unstaged changes in the current Git repository.

## Project Structure

```
code-reviewer/
├── cmd/                # CLI command definitions (Cobra)
│   ├── root.go         # Root command and global settings
│   ├── review.go       # 'review' command (main feature)
│   ├── report.go       # 'report' subcommand
│   └── fix.go          # 'fix' subcommand
├── internal/           # Private application logic
│   ├── agent/          # AI agent logic (Gemini integration, prompt management)
│   ├── config/         # Config file loading and saving
│   └── git/            # Git command execution wrapper
├── main.go             # Application entry point
├── go.mod              # Go module definition
└── README.md           # Project documentation
```
