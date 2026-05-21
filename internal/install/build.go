package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runGoBuild(repoRoot, dest string) error {
	cmd := exec.Command("go", "build", "-o", dest, ".")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}
	if _, err := os.Stat(dest); err != nil {
		return fmt.Errorf("binary missing after build: %s", dest)
	}
	abs, err := filepath.Abs(dest)
	if err == nil {
		fmt.Printf("Built %s\n", abs)
	}
	return nil
}
