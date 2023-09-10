package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GitCurrentBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func GitCurrentOrigin() (Origin, error) {
	return GitOrigin("")
}

func GitOrigin(path string) (Origin, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	if path != "" {
		cmd.Dir = path
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get origin url: %w", err)
	}
	return Origin(strings.TrimSpace(string(out))), nil
}

func IsGitRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func GitOrigins(dir string) (Origins, error) {
	var origins Origins
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() != dir {
			fullPath := filepath.Join(dir, "/", info.Name())

			if !IsGitRepo(fullPath) {
				return nil
			}

			origin, err := GitOrigin(fullPath)
			if err != nil {
				log.Printf("error getting origin for %s: %v", info.Name(), err)
				return nil
			}

			if origin != "" {
				origins = origins.Append(origin)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return origins, nil
}