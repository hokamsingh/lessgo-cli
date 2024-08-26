---

# LessGo CLI

[![Go Version](https://img.shields.io/github/go-mod/go-version/hokamsingh/lessgo-cli)](https://golang.org/dl/)
![Version](https://img.shields.io/badge/version-v1.0.5-blue)

`lessgo-cli` is a command-line tool designed to scaffold and manage projects using the [LessGo](https://github.com/hokamsingh/lessgo) framework. It simplifies the process of setting up a new project with pre-defined directory structures, boilerplate code, and configurations.

## Features

- **Quick Project Initialization**: Scaffold a new LessGo project in seconds with a single command.
- **Pre-configured Boilerplate**: Includes a basic project structure with ready-to-use controllers, modules, and services.
- **Easy Version Management**: Check the CLI version with the `--version` flag.
- **Modular Structure**: Follow a modular approach similar to the one used in NestJS, tailored for Go.

## Installation

To install the `lessgo-cli`, run the following command:

```bash
go install github.com/hokamsingh/lessgo-cli@latest
```

Make sure `$GOPATH/bin` is added to your system's PATH to run the CLI tool from anywhere.

### Windows Users

For Windows users, you can add the GOPATH to your system’s PATH as follows:

1. Press `Win + R`, type `sysdm.cpl`, and press Enter.
2. Go to the "Advanced" tab and click "Environment Variables."
3. Under "System variables," find the `Path` variable and click "Edit."
4. Add a new entry with `$(go env GOPATH)\bin`.
5. Click "OK" to save and apply the changes.

### Linux Users

For Linux users, ensure that your `GOPATH` binaries directory (usually `~/go/bin`) is included in your system's `PATH`. Follow these steps:

1. Open your terminal.
2. Edit your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.) with your favorite text editor:
   ```bash
   nano ~/.bashrc
   ```
   or
   ```bash
   nano ~/.zshrc
   ```
3. Add the following line at the end of the file:
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
4. Save the file and exit the editor.
5. Reload your shell configuration:
   ```bash
   source ~/.bashrc
   ```
   or
   ```bash
   source ~/.zshrc
   ```
6. Now you should be able to run `lessgo-cli` from anywhere in your terminal.

## Usage

### Initialize a New Project

To create a new LessGo project, run:

```bash
lessgo-cli new myapp
```

Replace `myapp` with your desired project name. The CLI will create a new directory with the specified name, containing the scaffolded project structure.

### Check CLI Version

To check the version of the CLI, run:

```bash
lessgo-cli --version
```

This will display the current version of `lessgo-cli` you have installed.

### Project Structure

When you run `lessgo-cli new myapp`, the following project structure is created:

```
myapp/
├── app/
│   ├── cmd/
│   │   └── main.go
│   └── src/
│       ├── app_controller.go
│       ├── app_module.go
│       └── app_service.go
├── go.mod
└── go.sum
```

- **`cmd/main.go`**: Entry point of your application.
- **`src/`**: Contains controllers, modules, and services for your application.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

If you have any questions or feedback, feel free to reach out:

- **GitHub**: [hokamsingh/lessgo-cli](https://github.com/hokamsingh/lessgo-cli)
- **Email**: hokamsingh07@gmail.com

---
