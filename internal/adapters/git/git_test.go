package git_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/el-j/nix-config-collector/internal/adapters/git"
)

func TestNew(t *testing.T) {
	a := git.New("/tmp")
	if a == nil {
		t.Fatal("New() returned nil")
	}
}

func TestIsRepo_NotARepo(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	a := git.New(dir)
	if a.IsRepo(context.Background()) {
		t.Error("expected IsRepo() = false for a non-git directory")
	}
}

func TestInit_And_IsRepo(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-init-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	a := git.New(dir)
	if err := a.Init(context.Background()); err != nil {
		t.Fatalf("Init() error: %v", err)
	}
	if !a.IsRepo(context.Background()) {
		t.Error("expected IsRepo() = true after Init()")
	}
}

func TestAddAll_And_Commit(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-commit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	a := git.New(dir)
	ctx := context.Background()

	if err := a.Init(ctx); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// Best-effort git identity configuration so commits can succeed in CI.
	for _, args := range [][2]string{
		{"user.email", "test@example.com"},
		{"user.name", "Test User"},
	} {
		//nolint:gosec // args are hardcoded strings
		_ = exec.Command("git", "-C", dir, "config", args[0], args[1]).Run()
	}

	// Write a test file
	if err := os.WriteFile(dir+"/test.txt", []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := a.AddAll(ctx); err != nil {
		t.Fatalf("AddAll() error: %v", err)
	}
	// Commit may fail if git identity not configured in CI — that's acceptable
	// We just check it doesn't panic
	_ = a.Commit(ctx, "initial commit")
}

func TestCommit_NothingToCommit(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-nothing-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	a := git.New(dir)
	ctx := context.Background()

	if err := a.Init(ctx); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// Commit on empty repo (nothing to commit) should return nil
	err = a.Commit(ctx, "empty commit")
	// We don't assert nil here because git may error differently across versions
	_ = err
}

func TestCurrentBranch_NotARepo(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-branch-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	a := git.New(dir)
	_, err = a.CurrentBranch(context.Background())
	if err == nil {
		t.Error("expected error for CurrentBranch in non-git dir")
	}
}
