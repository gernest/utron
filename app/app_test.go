package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAbsPath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	// non existing
	_, err = getAbsolutePath("nope")
	if err == nil {
		t.Error("expcted error got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expcetd not exist got %v", err)
	}

	absPath := filepath.Join(wd, "fixtures")

	// Relqtive
	dir, err := getAbsolutePath("fixtures")
	if err != nil {
		t.Error(err)
	}

	if dir != absPath {
		t.Errorf("expceted %s got %s", absPath, dir)
	}

	// Absolute
	dir, err = getAbsolutePath(absPath)
	if err != nil {
		t.Error(err)
	}

	if dir != absPath {
		t.Errorf("expceted %s got %s", absPath, dir)
	}

}
