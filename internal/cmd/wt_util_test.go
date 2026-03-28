package cmd

import (
	"testing"
)

func TestParseMainWorktreePath(t *testing.T) {
	porcelain := `worktree /projects/my-app/main
HEAD abc123
branch refs/heads/main

worktree /projects/my-app/feat/abc
HEAD def456
branch refs/heads/feat/abc

`
	got, err := parseMainWorktreePath(porcelain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/projects/my-app/main" {
		t.Errorf("got %q, want %q", got, "/projects/my-app/main")
	}
}

func TestParseMainWorktreePath_Empty(t *testing.T) {
	_, err := parseMainWorktreePath("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestRepoContainerFromMainPath(t *testing.T) {
	got := repoContainerFromMainPath("/projects/my-app/main")
	if got != "/projects/my-app" {
		t.Errorf("got %q, want %q", got, "/projects/my-app")
	}
}
