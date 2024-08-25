package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
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
	// CORS Options
	corsOptions := LessGo.NewCorsOptions(
		[]string{"*"}, 
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 
		[]string{"Content-Type", "Authorization"}, 
	)

	// Initialize App
	App := LessGo.App(
		LessGo.WithCORS(*corsOptions),
		LessGo.WithRateLimiter(100, 1*time.Minute),
		LessGo.WithJSONParser(),
		LessGo.WithCookieParser(),
	)

	// Serve Static Files
	folderPath, err := LessGo.GetFolderPath("uploads")
	if err != nil {
		log.Fatalf("Error: %%v", err)
	}
	App.ServeStatic("/static/", folderPath)

	// Dependency Injection Container
	container := LessGo.NewContainer()

	// Register Services
	if err := container.Register(%s.NewRootService); err != nil {
		log.Fatalf("Error registering RootService: %%v", err)
	}

	// Register Modules
	if err := container.Register(%s.NewRootModule); err != nil {
		log.Fatalf("Error registering UserModule: %%v", err)
	}

	// Root Module
	rootModule := %s.NewRootModule(App)
	LessGo.RegisterModules(App, []LessGo.IModule{rootModule})

	// Example Route
	App.Get("/ping", func(ctx *LessGo.Context) {
		ctx.Send("pong")
	})

	// Start the server
	serverPort := cfg.Get("SERVER_PORT", "8080")
	env := cfg.Get("ENV", "development")
	log.Printf("Starting server on port %%s in %%s mode", serverPort, env)
	if err := App.Listen(":" + serverPort); err != nil {
		log.Fatalf("Server failed: %%v", err)
	}
}
`, projectName, projectName, projectName, projectName, projectName)

	// Create and write to main.go
	mainGoPath := filepath.Join(projectDir, "app", "cmd", "main.go")
	err = os.WriteFile(mainGoPath, []byte(mainGoContent), os.ModePerm)
	if err != nil {
		fmt.Println("❌ Error creating main.go:", err)
		return
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
}
