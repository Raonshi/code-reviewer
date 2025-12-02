# AI Code Review Agent CLI

**AI Code Review Agent CLI**는 Google Gemini AI를 활용하여 Git 변경 사항을 자동으로 분석하고, 코드 품질 향상을 위한 전문적인 리뷰 리포트를 제공하는 도구입니다.

이 프로젝트는 개발자가 작성한 코드의 잠재적인 버그, 보안 취약점, 성능 이슈를 사전에 식별하고, 사용자가 설정한 언어(기본값: 한국어)로 구체적인 개선 방안을 제시하여 코드 리뷰 프로세스를 효율화하는 것을 목표로 합니다.

## Tech Stack

*   **Language**: Go 1.25.4
*   **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
*   **AI Model**: Google Gemini (설정 가능, 기본값: `gemini-2.5-flash`) (via [Google GenAI SDK](https://github.com/googleapis/go-genai))
*   **Configuration**: JSON based local config

## Getting Started

### Prerequisites

*   **Go**: 1.25 버전 이상 ([Download](https://go.dev/dl/))
*   **Git**: 프로젝트 버전 관리 및 변경 사항 추적을 위해 필요합니다.
*   **Google AI API Key**: Gemini 모델 사용을 위한 API 키가 필요합니다. ([Get API Key](https://aistudio.google.com/app/apikey))

### Homebrew

Homebrew를 사용하여 `code-reviewer`를 설치할 수 있습니다:

```bash
brew install Raonshi/tap/code-reviewer
```

### Installation

프로젝트를 클론하고 의존성을 설치합니다.

```bash
# Clone the repository
git clone https://github.com/your-username/code-reviewer.git
cd code-reviewer

# Download dependencies
go mod download
```

### Build & Run

프로젝트를 빌드하고 실행합니다.

```bash
# Build the binary
go build -o code-reviewer main.go

# Run the help command to verify installation
./code-reviewer --help
```

첫 실행 시, Google AI API Key, 사용할 AI 모델(기본값: `gemini-2.5-flash`), 그리고 출력 언어(기본값: `Korean`)를 입력하라는 메시지가 표시됩니다.

## Usage
```bash
# 코드 리뷰 실행 (기본값: 모든 변경 사항)
./code-reviewer report

# Staged 변경 사항에 대한 코드 리뷰 실행
./code-reviewer report --staged

# Unstaged 변경 사항에 대한 코드 리뷰 실행
./code-reviewer report --unstaged

# 코드 변경 사항에 대한 기술 문서 생성
./code-reviewer document

# 코드 문제에 대한 자동 수정 생성 (제안된 수정 사항 출력)
./code-reviewer fix

# 설정 관리
./code-reviewer config list
./code-reviewer config get output_language
./code-reviewer config set output_language English
```

## Key Features

*   **AI 기반 코드 리뷰**: Google Gemini 모델(설정 가능)을 사용하여 코드 변경 사항(`git diff`)을 심층 분석합니다.
*   **다국어 지원**: 다양한 언어로 리포트를 출력할 수 있습니다 (설정 가능, 기본값: 한국어).
*   **설정 관리**: `config` 명령어를 통해 API 키, 모델, 출력 언어를 쉽게 관리할 수 있습니다.
*   **자동 등급 분류**: 변경된 기능별로 **Good**, **Not Bad**, **Bad**, **Need Check** 4단계 등급을 매겨 중요도를 한눈에 파악할 수 있습니다.
*   **구체적인 개선 제안**: 단순한 지적을 넘어, 수정이 필요한 부분에 대해 바로 적용 가능한 코드 스니펫과 리팩토링 가이드를 제공합니다.
*   **자동 수정 생성**: 식별된 문제에 대한 수정된 코드를 생성하여 제안합니다 (현재는 표준 출력으로 제공).
*   **문서 생성**: `document` 명령어를 사용하여 코드 변경 사항에 대한 기술 문서를 자동으로 생성합니다.
*   **Git 통합**: 현재 Git 저장소의 Staged 및 Unstaged 변경 사항을 자동으로 감지합니다.