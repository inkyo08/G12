package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

type Runner struct{}

func (r Runner) Game() {
	log.Println("[BUILD] Game: started")

	if err := os.MkdirAll("Build", 0755); err != nil {
		log.Fatalf("[BUILD] Error: failed to create Build directory: %v", err)
	}

	dirs := []string{"Sources/Game"}
	files := []string{}
	for _, dir := range dirs {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && filepath.Ext(path) == ".hylo" {
				files = append(files, path)
			}
			return nil
		})
	}

	cmd := exec.Command("hc", append([]string{"--import-builtin", "-o", "Build/Game", "--verbose", "-O"}, files...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("[BUILD] Error: Game build failed: %v", err)
	}

	log.Println("\n[BUILD] Game: succeeded")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("[BUILD] Error: exactly one argument required")
	}

	method := reflect.ValueOf(Runner{}).MethodByName(os.Args[1])
	if !method.IsValid() {
		log.Fatalf("[BUILD] Error: unknown target %q", os.Args[1])
	}
	method.Call(nil)
}
