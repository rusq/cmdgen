package main

import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed template
var fsys embed.FS

func generate(r renderer, out output) error {
	for _, fn := range []func() error{r.validate, out.validate} {
		if err := fn(); err != nil {
			return fmt.Errorf("invalid parameters: %w", err)
		}
	}

	if err := os.MkdirAll(out.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to created output directory: %w", err)
	}

	fsys, err := fs.Sub(fsys, "template")
	if err != nil {
		return fmt.Errorf("failed to get sub filesystem: %w", err)
	}

	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "." {
				return nil
			}
			if err := os.MkdirAll(filepath.Join(out.OutputDir, path), 0755); err != nil {
				return err
			}
			return nil
		}
		output := filepath.Join(out.OutputDir, path)
		if filepath.Ext(path) != ".go" {
			if err := copyfile(fsys, path, output); err != nil {
				return err
			}
			return nil
		}
		if err := replace(fsys, path, output, r.replacer().Replace); err != nil {
			return err
		}
		return nil
	})
}

func copyfile(fsys fs.FS, path, output string) error {
	src, err := fsys.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(output)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func replace(fsys fs.FS, path, output string, replacer func(string) string) error {
	contents, err := fs.ReadFile(fsys, path)
	if err != nil {
		return err
	}
	if err := os.WriteFile(output, []byte(replacer(string(contents))), 0644); err != nil {
		return err
	}
	return nil
}
