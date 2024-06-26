package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ktr0731/go-fuzzyfinder"
)

func dirOr(path string) (string, error) {
	d := filepath.Dir(path)
	ss := []string{path, d}
	idx, err := fuzzyfinder.Find(ss, func(i int) string {
		p := ss[i]
		rel, _ := filepath.Rel(filepath.Dir(d), p)
		return filepath.ToSlash(rel)
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return "", err
	}
	return ss[idx], nil
}

func defaultOpen(path string) error {
	return exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

type Filer struct {
	path string
}

func (flr *Filer) Init(path string) {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		flr.path = path
		return
	}
	flr.path = "explorer.exe"
}

func (flr Filer) open(path string) error {
	return exec.Command(flr.path, path).Start()
}

func (flr Filer) OpenSmart(path string, curDir string) error {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return flr.open(path)
	}
	if filepath.Dir(path) == curDir {
		return defaultOpen(path)
	}
	p, err := dirOr(path)
	if err != nil {
		return err
	}
	if fi, err := os.Stat(p); err == nil && fi.IsDir() {
		return flr.open(p)
	}
	return defaultOpen(p)
}
