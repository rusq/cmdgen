package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestGenerateCompiles generates a project from the embedded template and
// verifies that it compiles in both release and debug modes.
func TestGenerateCompiles(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go binary not found in PATH")
	}

	dir := t.TempDir()
	outdir := filepath.Join(dir, "cmd", "foo")

	r := renderer{
		Command:         "foo",
		CommandVariable: "FooCommand",
		Package:         "github.com/acme/foo/cmd/foo",
	}
	if err := generate(r, output{OutputDir: outdir}); err != nil {
		t.Fatalf("generate: %s", err)
	}

	gomod := "module github.com/acme/foo\n\ngo 1.23\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644); err != nil {
		t.Fatal(err)
	}

	for _, tags := range [][]string{nil, {"-tags", "debug"}} {
		args := append([]string{"build"}, tags...)
		args = append(args, "./...")
		cmd := exec.Command("go", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("go %v: %s\n%s", args, err, out)
		}
	}
}

func TestRendererValidate(t *testing.T) {
	tests := []struct {
		name    string
		r       renderer
		wantVar string
		wantErr bool
	}{
		{"valid", renderer{Command: "foo", CommandVariable: "Foo", Package: "github.com/acme/foo"}, "Foo", false},
		{"lowercase is capitalised", renderer{Command: "foo", CommandVariable: "fooCmd", Package: "github.com/acme/foo"}, "FooCmd", false},
		{"non-letter start", renderer{Command: "foo", CommandVariable: "_foo", Package: "github.com/acme/foo"}, "", true},
		{"empty command", renderer{CommandVariable: "Foo", Package: "github.com/acme/foo"}, "", true},
		{"empty variable", renderer{Command: "foo", Package: "github.com/acme/foo"}, "", true},
		{"empty package", renderer{Command: "foo", CommandVariable: "Foo"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.r.CommandVariable != tt.wantVar {
				t.Errorf("CommandVariable = %q, want %q", tt.r.CommandVariable, tt.wantVar)
			}
		})
	}
}
