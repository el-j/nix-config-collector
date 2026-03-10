package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// Adapter implements git operations for managing generated Nix configs.
type Adapter struct {
	repoPath string
}

// New creates a new Git Adapter for the given repository path.
func New(repoPath string) *Adapter {
	return &Adapter{repoPath: repoPath}
}

// IsRepo returns true if the path is inside a git repository.
func (a *Adapter) IsRepo(ctx context.Context) bool {
	_, err := a.run(ctx, "git", "-C", a.repoPath, "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// Init initialises a new git repository at the adapter's path.
func (a *Adapter) Init(ctx context.Context) error {
	_, err := a.run(ctx, "git", "-C", a.repoPath, "init")
	return err
}

// AddAll stages all changes in the repository.
func (a *Adapter) AddAll(ctx context.Context) error {
	_, err := a.run(ctx, "git", "-C", a.repoPath, "add", ".")
	return err
}

// Commit creates a commit with the given message.
// Returns nil if there is nothing to commit.
func (a *Adapter) Commit(ctx context.Context, message string) error {
	out, err := a.run(ctx, "git", "-C", a.repoPath, "commit", "-m", message)
	if err != nil {
		// "nothing to commit" is not a real error
		if bytes.Contains(out, []byte("nothing to commit")) {
			return nil
		}
		return err
	}
	return nil
}

// Tag creates an annotated tag.
func (a *Adapter) Tag(ctx context.Context, tag, message string) error {
	_, err := a.run(ctx, "git", "-C", a.repoPath, "tag", "-a", tag, "-m", message)
	return err
}

// Push pushes the current branch and tags to origin.
func (a *Adapter) Push(ctx context.Context) error {
	if _, err := a.run(ctx, "git", "-C", a.repoPath, "push"); err != nil {
		return err
	}
	_, err := a.run(ctx, "git", "-C", a.repoPath, "push", "--tags")
	return err
}

// CurrentBranch returns the name of the current branch.
func (a *Adapter) CurrentBranch(ctx context.Context) (string, error) {
	out, err := a.run(ctx, "git", "-C", a.repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	branch := string(bytes.TrimSpace(out))
	return branch, nil
}

// CommitAndPush is a convenience method that stages all changes,
// commits with the given message, and pushes to origin.
func (a *Adapter) CommitAndPush(ctx context.Context, message string) error {
	if err := a.AddAll(ctx); err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	if err := a.Commit(ctx, message); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	if err := a.Push(ctx); err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil
}

func (a *Adapter) run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.Bytes(), fmt.Errorf("running %s %v: %w\noutput: %s", name, args, err, out.String())
	}
	return out.Bytes(), nil
}
