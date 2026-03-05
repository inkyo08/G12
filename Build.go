package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
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

	libs := []string{}
	libs = append(libs, findLib("SDL3")...)

	cmd := exec.Command("hc", append(append([]string{"--import-builtin", "-o", "Build/Game", "--verbose", "-O"}, libs...), files...)...)
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

func findLib(name string) []string {
	pkgOut, err := exec.Command("pkg-config", "--libs-only-L", "--libs-only-l", name).Output()
	if err != nil {
		log.Fatalf("[BUILD] Error: pkg-config %s failed: %v", name, err)
	}
	pkgFlags := []string{}
	for _, f := range strings.Fields(string(pkgOut)) {
		if strings.HasPrefix(f, "-L") {
			pkgFlags = append(pkgFlags, "-L", f[2:])
		} else if strings.HasPrefix(f, "-l") {
			pkgFlags = append(pkgFlags, "-l", f[2:])
		} else {
			pkgFlags = append(pkgFlags, f)
		}
	}
	return pkgFlags
}
