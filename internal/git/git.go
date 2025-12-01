package git

import (
	"errors"
	"os"
	"os/exec"
)

// IsRepo checks if the current directory is a git repository.
func IsRepo() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

// DiffMode defines the type of diff to retrieve
type DiffMode string

const (
	DiffModeAll      DiffMode = "all"
	DiffModeStaged   DiffMode = "staged"
	DiffModeUnstaged DiffMode = "unstaged"
)

// GetDiff retrieves the diff of changes based on the mode.
func GetDiff(mode DiffMode) (string, error) {
	args := []string{"diff"}
	switch mode {
	case DiffModeStaged:
		args = append(args, "--staged")
	case DiffModeAll:
		args = append(args, "HEAD")
	case DiffModeUnstaged:
		// Default behavior (git diff)
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// ApplyPatch applies a patch to the repository.
// This is a simplified implementation that expects a standard diff format.
// For complex patching, we might need to write the patch to a file and use `git apply`.
func ApplyPatch(patchContent string) error {
	// Create a temporary file for the patch
	f, err := os.CreateTemp("", "patch-*.diff")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(patchContent); err != nil {
		return err
	}
	f.Close()

	cmd := exec.Command("git", "apply", f.Name())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
