package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func CloneCiCdToolkit() error {

	const (
		cicdToolkitUrl = "git@github.com:zendesk/cicd-toolkit.git"
	)

	path := filepath.Join(os.Getenv("HOME"), ".zendesk", "cicd-toolkit")

	// if not found, clone the repo
	if _, err := os.Stat(path); err != nil {
		_ = os.MkdirAll(path, 0755)
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      cicdToolkitUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	} else {
		repo, err := git.PlainOpen(path)
		if err != nil {
			return err
		}

		wt, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("error calling Worktree(): %v", err)
		}

		err = wt.Pull(&git.PullOptions{
			Progress: os.Stdout,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "up-to-date") {
				return err
			}
		}
	}

	return nil
}
