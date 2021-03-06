package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Repository defines details to clone a dependency.
type Repository struct {
	// Dir contains the relative path where git will clone the repository from the root of the environment.
	Dir string
	// URL contains the URL of the repository.
	URL string
	// Commit contains the hash that will be checked out.
	Commit string
}

// Build clones and checks out the repository at the specified commit inside the environment.
func (r *Repository) Build(root string) error {
	dir := filepath.Join(root, r.Dir)

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir

	if stdout, err := cmd.Output(); err == nil {
		if r.Commit == "" || r.Commit == string(stdout) {
			return nil
		}
	} else {
		fmt.Println("getting", r.Dir)

		cmd = exec.Command("git", "clone", "-q", r.URL, r.Dir)
		cmd.Dir = root
		cmd.Stderr = os.Stderr
		cmd.Env = []string{
			"GIT_TERMINAL_PROMPT=0",
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to 'git clone %s %s', got %s", r.URL, r.Dir, err)
		}
	}

	if r.Commit == "" {
		return nil
	}

	cmd = exec.Command("git", "checkout", "-q", r.Commit)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to 'git checkout %s', got %s", r.Commit, err)
	}

	return nil
}
