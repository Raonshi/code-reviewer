# 바이너리 기본 이름
APP_NAME=code-reviewer

# 빌드 출력 디렉토리 (없으면 생성됨)
DIST_DIR=dist

# 기본 타겟: make만 입력하면 실행됨
all: clean build-mac build-win build-linux

# 1. Mac (Apple Silicon) 빌드
build-mac:
	@echo "Building for Mac (M-Series)..."
	mkdir -p $(DIST_DIR)/mac
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/mac/$(APP_NAME) main.go

# 2. Windows (Intel x64) 빌드
build-win:
	@echo "Building for Windows (Intel x64)..."
	mkdir -p $(DIST_DIR)/win
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(DIST_DIR)/win/$(APP_NAME) main.go

# 3. Linux (Intel x64) 빌드
build-linux:
	@echo "Building for Linux (Intel x64)..."
	mkdir -p $(DIST_DIR)/linux
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/linux/$(APP_NAME) main.go

# 청소: 기존 빌드 파일 삭제
clean:
	@echo "Cleaning up..."
	rm -rf $(DIST_DIR)