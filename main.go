package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const version = "v1.0.3"

func main() {
	// Check for the --version flag
	versionFlag := flag.Bool("version", false, "Print the version number")
	flag.Parse()

	if *versionFlag {
		fmt.Println("lessgo-cli version", version)
		os.Exit(0)
	}

	// Handle other commands (e.g., "new")
	if len(os.Args) < 2 {
		fmt.Println("Expected 'new' command.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		// Ask for the project name
		var projectName string
		fmt.Print("Enter project name: ")
		fmt.Scanln(&projectName)

		// ASCII art for Go's logo
		goLogo := `       *((((((((((                                                    (((*      
    @@@@@@@@@@@@@@                                                    @@@@@@@                                       @@@@@@                  @@@@                
   (@@@@    @@@@@@                                                       @@@@/                                 @@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@           
   *@@@@    @@@@@@        @@@@@@           @@@@@@@        @@@@@@@        @@@@,                              @@@@@@@@@@@@@@@@@@@@@  @@@@@@@@@@@@@@@@@@@@         
    @@@@    @@@@@@    *@@@@@@@@@@@/     &@@@@@@@@@@@&   @@@@@@@@@@@@&    @@@@                              @@@@@@@        @@@    @@@@@@@@        @@@@@@@        
    @@@@    @@@@@@  ,@@@@@@   /@@@@@   @@@@@&   (@@/   @@@@@%   (@@/     @@@@                             @@@@@@@    @@@@@@@@@@@@@@@@@@           @@@@@@@       
./@@@@@@    @@@@@@ ,@@@@@       @@@@@  @@@@@@.         @@@@@@            @@@@@@(                 @@@@@@@  @@@@@@@    @@@@@@@@@@@@@@@@@@            @@@@@@@       
(@@@@@@     @@@@@@ @@@@@@@@@@@@@@@@@@   /@@@@@@@@@@     (@@@@@@@@@@       @@@@@@                           @@@@@@@         @@@@@@@@@@@@@          @@@@@@@        
   *@@@@    @@@@@@ /@@@@@                    &@@@@@@@        &@@@@@@@    @@@@,                            @@@@@@@@@*   @@@@@@@@& @@@@@@@@@   /@@@@@@@@@         
    @@@@    @@@@@@  @@@@@@#       ,     @.      @@@@@   @.      @@@@@    @@@@                               @@@@@@@@@@@@@@@@@@    @@@@@@@@@@@@@@@@@@@           
   .@@@@    ,@@@@@@@  @@@@@@@@@@@@@@  @@@@@@@@@@@@@@% @@@@@@@@@@@@@@(    @@@@                                  @@@@@@@@@@@           @@@@@@@@@@@@               
   (@@@@      .@@@@@.    .&@@@@@*        #@@@@@@@(       %@@@@@@@/       @@@@/  
    @@@@@                                                               &@@@@   
     @@@@@@                                                           @@@@@@    
`
		fmt.Println(goLogo)
		fmt.Printf("🚀 Initializing your Less%s project...\n\n", projectName)

		// Create project directory structure
		projectDir := filepath.Join(".", projectName)
		err := os.MkdirAll(filepath.Join(projectDir, "app", "cmd"), os.ModePerm)
		if err != nil {
			fmt.Println("❌ Error creating project directories:", err)
			return
		}

		// Verify the creation of the directory
		if _, err := os.Stat(filepath.Join(projectDir, "app", "cmd")); os.IsNotExist(err) {
			fmt.Println("❌ Failed to create the app/cmd directory.")
			return
		}

		// Dynamically create the content of main.go using the project name
		mainGoContent := fmt.Sprintf(`package main

import (
	%s "%s/app/src"
	"log"
	"time"

	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

func main() {
	// Load Configuration
	cfg := LessGo.LoadConfig()
	serverPort := cfg.Get("SERVER_PORT", "8080")
	env := cfg.Get("ENV", "development")
	addr := ":" + serverPort
	// CORS Options
	corsOptions := LessGo.NewCorsOptions(
		[]string{"*"}, 
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 
		[]string{"Content-Type", "Authorization"}, 
	)
	// Parser Options
	size := LessGo.ConvertToBytes(int64(1024), LessGo.Kilobytes)
	parserOptions := LessGo.NewParserOptions(size * 5)

	// redis client
	rClient := LessGo.NewRedisClient("localhost:6379")

	// Initialize App with Middlewares
	App := LessGo.App(
		LessGo.WithCORS(*corsOptions),
		// LessGo.WithInMemoryRateLimiter(4, 50, 1*time.Second, 5*time.Minute), // Rate limiter
		// LessGo.WithRedisRateLimiter("localhost:6379", 10, time.Minute*5),
		LessGo.WithJSONParser(*parserOptions),
		LessGo.WithCookieParser(),                        // Cookie parser
		LessGo.WithCsrf(),                                // CSRF protection middleware
		LessGo.WithXss(),                                 // XSS protection middleware
		LessGo.WithCaching(rClient, 5*time.Minute, true), // Caching middleware using Redis
		// LessGo.WithFileUpload("uploads"), // Uncomment if you want to handle file uploads you could use this in a seprate controller
	)

	// Serve Static Files
	 folderPath := LessGo.GetFolderPath("uploads")
	 App.ServeStatic("/static/", folderPath)

	// Root Module
	rootModule := src.NewRootModule(App)
	LessGo.RegisterModules(App, []LessGo.IModule{rootModule})

	// Example Route
	App.Get("/ping", func(ctx *LessGo.Context) {
		ctx.Send("pong")
	})

	// Start the server
	log.Printf("Starting server on port %%s in %%s mode", serverPort, env)
	if err := App.Listen(addr); err != nil {
		log.Fatalf("Server failed: %%v", err)
	}
}
`, projectName, projectName)

		envContent := `SERVER_PORT =9004
					   ENV=development
					   JWT_SECRET=secret
                       REDIS_ADDR=localhost:6379`

		dotAirDotYml := `# .air.toml

# General settings
[build]
# The binary name generated by the build command.
bin = "app_binary"

# The path to your main Go package. Assuming your main.go is in the cmd directory.
cmd = "go build -o ./{{root_dir}}/cmd/app_binary ./{{root_dir}}/cmd"

# Directory where air will watch files for changes.
watch_dir = "./{{root_dir}}"

# Patterns of files to watch.
include_ext = ["go", "tpl", "tmpl", "html"]

# Exclude directories or files from being watched.
exclude_dir = ["uploads", "dist"]

# Execute the application after each rebuild.
run_cmd = "./{{root_dir}}/cmd/app_binary"

# Run the application as a daemon.
run_args = []

# Environment variables to pass to the application.
[build.envs]
ENV = "development"
PORT = "8080"

# Custom commands before and after build.
pre_build = "go fmt ./..."
post_build = ""

# Configuring logging.
[log]
time = true
color = true

# Live reload settings.
[live]
# Number of lines to show in the output.
output_lines = 100
`
		dockerComposeContent := `version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_PORT=8080
      - ENV=development
      - REDIS_ADDR=redis:6379
    depends_on:
      - redis

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
`
		dockerfileContent := `# Set build arguments to make it dynamic
ARG GO_VERSION=1.22-bullseye
ARG APP_DIR=/app
ARG CMD_DIR=${APP_DIR}/cmd
ARG OUTPUT_BINARY=/main
ARG PORT=8080

# Use the official Golang image to build the Go binary
FROM golang:${GO_VERSION} as builder

# Set the Current Working Directory inside the container
WORKDIR ${APP_DIR}

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Set the working directory to the cmd directory
WORKDIR ${CMD_DIR}

# Build the Go app and output it to the specified location
RUN go build -o ${OUTPUT_BINARY} .

# Start a new stage from scratch
FROM debian:bullseye-slim

# Set runtime environment variables
ARG OUTPUT_BINARY
ARG PORT

# Copy the pre-built binary file from the previous stage
COPY --from=builder ${OUTPUT_BINARY} ${OUTPUT_BINARY}

# Expose the port to the outside world
EXPOSE ${PORT}

# Command to run the executable
CMD ["/main"]
`

		makefileContent := `# Default Variables
GOFMT_FILES ?= $(shell find . -name '*.go')
APP_DIR ?= ./app
CMD_DIR ?= $(APP_DIR)/cmd
MAIN_FILE ?= main.go
BIN_DIR ?= ./bin
OUTPUT_BINARY ?= $(BIN_DIR)/main
COVERAGE_FILE ?= coverage.out

# External Tools (override these if you want to use different tools or versions)
GOFMT ?= gofmt
LINT_TOOL ?= golangci-lint
GOSEC ?= gosec
GO ?= go

# Targets
.PHONY: all fmt lint test build clean security run deps help

all: fmt lint test build

# Format Go files
fmt:
	@$(GOFMT) -s -w $(GOFMT_FILES)
	@echo "Code formatted with $(GOFMT)."

# Lint the project using golangci-lint
lint:
	@$(LINT_TOOL) run
	@echo "Linting completed with $(LINT_TOOL)."

# Run tests
test:
	@$(GO) test ./... -coverprofile=$(COVERAGE_FILE)
	@echo "Tests completed. Coverage report generated at $(COVERAGE_FILE)."

# Build the project
build:
	@$(GO) build -o $(OUTPUT_BINARY) $(CMD_DIR)/$(MAIN_FILE)
	@echo "Build completed. Output binary is at $(OUTPUT_BINARY)."

# Clean the build files
clean:
	@rm -rf $(BIN_DIR)
	@echo "Cleanup completed. Removed $(BIN_DIR)."

# Run gosec for security checks
security:
	@$(GOSEC) ./...
	@echo "Security check completed with $(GOSEC)."

# Run the application
run:
	@$(GO) run $(CMD_DIR)/$(MAIN_FILE)

# Install dependencies
deps:
	@$(GO) mod tidy
	@echo "Dependencies installed."

# Help target to display available commands
help:
	@echo "Makefile targets:"
	@echo "  fmt      - Format the Go source files using $(GOFMT)"
	@echo "  lint     - Lint the project using $(LINT_TOOL)"
	@echo "  test     - Run tests and generate a coverage report at $(COVERAGE_FILE)"
	@echo "  build    - Build the project, output binary at $(OUTPUT_BINARY)"
	@echo "  clean    - Remove build artifacts in $(BIN_DIR)"
	@echo "  security - Run security checks with $(GOSEC)"
	@echo "  run      - Run the application at $(CMD_DIR)/$(MAIN_FILE)"
	@echo "  deps     - Install dependencies"
	@echo "  all      - Run all the above targets"
`
		// Create and write to main.go
		mainGoPath := filepath.Join(projectDir, "app", "cmd", "main.go")
		err = os.WriteFile(mainGoPath, []byte(mainGoContent), os.ModePerm)
		if err != nil {
			fmt.Println("❌ Error creating main.go:", err)
			return
		}
		rContents := []string{envContent, makefileContent, dotAirDotYml, dockerComposeContent, dockerfileContent}
		for i, file := range rContents {
			filePath := filepath.Join(projectDir, file)
			content := rContents[i]
			err = os.WriteFile(filePath, []byte(content), os.ModePerm)
			if err != nil {
				fmt.Println("❌ Error creating", file, ":", err)
				return
			}
		}

		// Create src directory
		srcDir := filepath.Join(projectDir, "app", "src")
		err = os.MkdirAll(srcDir, os.ModePerm)
		if err != nil {
			fmt.Println("❌ Error creating src directory:", err)
			return
		}

		// Verify the creation of the src directory
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			fmt.Println("❌ Failed to create the app/src directory.")
			return
		}

		// Dynamically create the content of src files using the project name
		controllerContent := fmt.Sprintf(`package %s

import LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"

type RootController struct {
	LessGo.BaseController
	Path    string
	Service RootService
}

func NewRootController(s *RootService, path string) *RootController {
	return &RootController{
		Path:    path,
		Service: *s,
	}
}

func (rc *RootController) RegisterRoutes(r *LessGo.Router) {
	// r.Get("/hello", func(ctx *LessGo.Context) {
	// 	ctx.Send("Hello world")
	// })
}
`, projectName)

		moduleContent := fmt.Sprintf(`package %s

import (
	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type RootModule struct {
	LessGo.Module
}

func NewRootModule(r *LessGo.Router) *RootModule {
	modules := []LessGo.IModule{}
	LessGo.RegisterModules(r, modules)
	service := NewRootService()
	controller := NewRootController(service, "/")
	return &RootModule{
		Module: *LessGo.NewModule("Root", []interface{}{controller}, []interface{}{service}, modules),
	}
}
`, projectName)

		serviceContent := fmt.Sprintf(`package %s

type IRootService interface{}

type RootService struct {
	// Add any shared dependencies or methods here
}

func NewRootService() *RootService {
	return &RootService{}
}
`, projectName)

		// File contents for src files
		rootContents := []string{controllerContent, moduleContent, serviceContent}

		// Create src files in app/src directory
		srcFiles := []string{"app_controller.go", "app_module.go", "app_service.go"}
		for i, file := range srcFiles {
			filePath := filepath.Join(srcDir, file)
			content := rootContents[i]
			err = os.WriteFile(filePath, []byte(content), os.ModePerm)
			if err != nil {
				fmt.Println("❌ Error creating", file, ":", err)
				return
			}
		}

		// Initialize go.mod
		cmd := exec.Command("go", "mod", "init", projectName)
		cmd.Dir = projectDir
		err = cmd.Run()
		if err != nil {
			fmt.Println("❌ Error initializing go.mod:", err)
			return
		}

		cmd = exec.Command("go", "mod", "tidy")
		cmd.Dir = projectDir
		err = cmd.Run()
		if err != nil {
			fmt.Println("❌ Error running go mod tidy:", err)
			return
		}

		fmt.Println("🎉 Project scaffold created successfully!")
		fmt.Printf("🚀 Spin up your new LessGo app by running: go run %s/app/cmd/main.go\n", projectName)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
