# Go (Golang) Basics

This document serves as a quick reference for the Go programming language, covering essential commands and core concepts.

## Main Commands

Here are the most frequently used `go` commands:

- **`go run <file.go>`**: Compiles and executes a Go program in one step. Useful for development and testing.
- **`go build`**: Compiles the packages and dependencies, creating an executable binary in the current directory.
- **`go test`**: Runs test functions (functions starting with `Test` in files named `*_test.go`).
- **`go fmt`**: Automatically formats Go source code according to the standard Go styling rules.
- **`go mod init <module-name>`**: Initializes a new module in the current directory, creating a `go.mod` file.
- **`go get <package-url>`**: Downloads and installs third-party packages and adds them as dependencies to `go.mod`.
- **`go mod tidy`**: Cleans up the `go.mod` and `go.sum` files by adding missing dependencies and removing unused ones.

## What is a Package?

A **package** is a way to group related Go source files together. 
- Every Go source file must belong to a package, declared at the very top of the file (e.g., `package main` or `package math`).
- All files within the same directory must belong to the same package.
- Variables, functions, and types starting with an uppercase letter are **exported** (accessible from other packages), while those starting with a lowercase letter are **unexported** (private to the package).
- The `main` package is special: it defines a standalone standalone executable program rather than a library.

## What is a Module?

A **module** is a collection of related Go packages that are versioned together as a single unit.
- A module is defined by a tree of Go source files with a `go.mod` file in the tree's root directory.
- The `go.mod` file defines the module's path (which also serves as the import path for its packages) and its dependency requirements.

## How are Dependencies Managed?

Dependency management in Go is handled via **Go Modules**.

1. **`go.mod`**: When you run `go mod init`, this file is created. It tracks the exact versions of the third-party modules your project depends on.
2. **`go.sum`**: When you add a new dependency (e.g., using `go get`), Go automatically updates the `go.mod` file and generates/updates a `go.sum` file. The `go.sum` file contains cryptographic hashes of the dependencies to ensure their integrity and that the exact same code is used every time the project is built.
3. **Fetching**: The `go get` command is used to download dependencies. Go downloads these to a local cache (usually in `$GOPATH/pkg/mod`) so they can be reused across different projects on the same machine.
4. **Maintenance**: Running `go mod tidy` ensures that your `go.mod` and `go.sum` accurately reflect what your project actually imports, keeping things clean.
