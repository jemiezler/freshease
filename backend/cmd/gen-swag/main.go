// gen-swag.go — run with: go run ./cmd/gen-swag
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// findRepoRoot walks up until it finds go.mod
func findRepoRoot(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s upward", start)
		}
		dir = parent
	}
}

func main() {
	fmt.Println("Generating Swagger documentation for Freshease...")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd: %v", err)
	}

	// Ensure we run from the repo root (where go.mod lives)
	root, err := findRepoRoot(wd)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Entry file that has @title/@BasePath comments
	entry := filepath.ToSlash(filepath.Join("main.go"))

	// Output folder must match your import: _ "freshease/backend/internal/docs"
	out := filepath.ToSlash(filepath.Join("internal", "docs"))

	// Ensure swag exists (optional)
	if _, err := exec.LookPath("swag"); err != nil {
		fmt.Println("Installing swag CLI...")
		install := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
		install.Stdout, install.Stderr = os.Stdout, os.Stderr
		if err := install.Run(); err != nil {
			log.Fatalf("failed to install swag: %v", err)
		}
	}

	args := []string{
		"init",
		"-g", entry, // relative to repo root
		"--dir", ".", // IMPORTANT: relative, not absolute
		"--parseDependency",
		"--parseInternal",
		"--output", out, // relative to repo root
	}

	fmt.Printf("Running (in %s): swag %v\n", root, args)
	cmd := exec.Command("swag", args...)
	cmd.Dir = root // run from repo root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Swagger generation failed: %v", err)
	}

	fmt.Printf("Swagger docs generated at %s\n", filepath.Join(root, out))
}
