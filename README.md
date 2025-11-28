# AI Code Review Agent CLI

**AI Code Review Agent CLI**는 Google Gemini AI를 활용하여 Git 변경 사항을 자동으로 분석하고, 코드 품질 향상을 위한 전문적인 리뷰 리포트를 제공하는 도구입니다.

이 프로젝트는 개발자가 작성한 코드의 잠재적인 버그, 보안 취약점, 성능 이슈를 사전에 식별하고, 한국어로 구체적인 개선 방안을 제시하여 코드 리뷰 프로세스를 효율화하는 것을 목표로 합니다.

## Tech Stack

*   **Language**: Go 1.25.4
*   **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
*   **AI Model**: Google Gemini 2.5 Flash (via [Google GenAI SDK](https://github.com/googleapis/go-genai))
*   **Configuration**: JSON based local config

## Getting Started

### Prerequisites

*   **Go**: 1.25 버전 이상 ([Download](https://go.dev/dl/))
*   **Git**: 프로젝트 버전 관리 및 변경 사항 추적을 위해 필요합니다.
*   **Google AI API Key**: Gemini 모델 사용을 위한 API 키가 필요합니다. ([Get API Key](https://aistudio.google.com/app/apikey))

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

첫 실행 시, Google AI API Key를 입력하라는 메시지가 표시됩니다. 

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

*   **AI 기반 코드 리뷰**: Google Gemini 2.5 Flash 모델을 사용하여 코드 변경 사항(`git diff`)을 심층 분석합니다.
*   **한국어 리포트**: 모든 분석 결과와 개선 제안은 한국어로 제공되어 이해하기 쉽습니다.
*   **자동 등급 분류**: 변경된 기능별로 **Good**, **Not Bad**, **Bad**, **Need Check** 4단계 등급을 매겨 중요도를 한눈에 파악할 수 있습니다.
*   **구체적인 개선 제안**: 단순한 지적을 넘어, 수정이 필요한 부분에 대해 바로 적용 가능한 코드 스니펫과 리팩토링 가이드를 제공합니다.
*   **Git 통합**: 현재 Git 저장소의 Staged 및 Unstaged 변경 사항을 자동으로 감지합니다.

## Project Structure

```
code-reviewer/
├── cmd/                # CLI 명령어 정의 (Cobra)
│   ├── root.go         # 루트 명령어 및 전역 설정
│   ├── review.go       # 'review' 명령어 (메인 기능)
│   ├── report.go       # 'report' 서브 명령어
│   └── fix.go          # 'fix' 서브 명령어
├── internal/           # 비공개 애플리케이션 로직
│   ├── agent/          # AI 에이전트 로직 (Gemini 연동, 프롬프트 관리)
│   ├── config/         # 설정 파일 로드 및 저장
│   └── git/            # Git 명령어 실행 래퍼
├── main.go             # 애플리케이션 진입점
├── go.mod              # Go 모듈 정의
└── README.md           # 프로젝트 문서
```
