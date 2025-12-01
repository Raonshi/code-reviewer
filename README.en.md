# AI Code Review Agent CLI

**AI Code Review Agent CLI** is a tool that leverages Google Gemini AI to automatically analyze Git changes and provide professional review reports for code quality improvement.

This project aims to streamline the code review process by proactively identifying potential bugs, security vulnerabilities, and performance issues in developer code, and providing concrete improvement suggestions in Korean (or English).

## Tech Stack

*   **Language**: Go 1.25.4
*   **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
*   **AI Model**: Google Gemini (Configurable, default: `gemini-2.5-flash`) (via [Google GenAI SDK](https://github.com/googleapis/go-genai))
*   **Configuration**: JSON based local config

## Getting Started

### Prerequisites

*   **Go**: Version 1.25 or higher ([Download](https://go.dev/dl/))
*   **Git**: Required for project version control and tracking changes.
*   **Google AI API Key**: API key required for using the Gemini model. ([Get API Key](https://aistudio.google.com/app/apikey))

### Homebrew

You can install `code-reviewer` using Homebrew:

```bash
brew install Raonshi/tap/code-reviewer
```

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

On the first run, you will be prompted to enter your Google AI API Key and select an AI Model (default: `gemini-2.5-flash`).

## Usage
```bash
# Run code review (default: all changes)
./code-reviewer report

# Run code review on staged changes
./code-reviewer report --staged

# Run code review on unstaged changes
./code-reviewer report --unstaged

# Generate auto-fixes for code issues (NOT SUPPORTED YET)
./code-reviewer fix
```

## Key Features

*   **AI-Powered Code Review**: Uses Google Gemini models (configurable) to deeply analyze code changes (`git diff`).
*   **Korean Reports**: All analysis results and improvement suggestions are provided in Korean for easy understanding.
*   **Automatic Grading**: Changes are automatically graded into 4 levels: **Good**, **Not Bad**, **Bad**, **Need Check** to quickly identify importance.
*   **Concrete Improvement Suggestions**: Beyond simple pointers, it provides immediately applicable code snippets and refactoring guides for areas needing correction.
*   **Git Integration**: Automatically detects Staged and Unstaged changes in the current Git repository.