// Command web builds the browser bundle and optionally serves it on localhost.
package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	serve := flag.Bool("serve", false, "Serve the built game on localhost")
	port := flag.Int("port", 8080, "Local HTTP port")
	flag.Parse()
	if err := build(); err != nil {
		log.Fatal(err)
	}
	if !*serve {
		fmt.Println("Browser bundle: build/web")
		return
	}
	mime.AddExtensionType(".wasm", "application/wasm")
	fmt.Printf("Open http://localhost:%d (Ctrl+C to stop)\n", *port)
	files := http.FileServer(http.Dir("build/web"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		files.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *port), handler))
}

func build() error {
	if err := os.MkdirAll("build/web", 0755); err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "-o", "build/web/game.wasm", ".")
	for _, value := range os.Environ() {
		name, _, _ := strings.Cut(value, "=")
		if !strings.EqualFold(name, "GOOS") && !strings.EqualFold(name, "GOARCH") {
			cmd.Env = append(cmd.Env, value)
		}
	}
	cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build Wasm: %w", err)
	}
	root, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return err
	}
	runtime := filepath.Join(strings.TrimSpace(string(root)), "lib", "wasm", "wasm_exec.js")
	if err := copyFile(runtime, "build/web/wasm_exec.js"); err != nil {
		return err
	}
	for _, name := range []string{"index.html", "launcher.js", "game.html", "game.js", "styles.css"} {
		if err := copyFile(filepath.Join("web", name), filepath.Join("build/web", name)); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(source, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return os.WriteFile(target, data, 0644)
}
