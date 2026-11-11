package download

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temporary local git repository for testing
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	srcDir, err := os.MkdirTemp("", "ryvst-test-src-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	runGit(t, srcDir, "init")
	runGit(t, srcDir, "config", "user.email", "test@example.com")
	runGit(t, srcDir, "config", "user.name", "Test User")

	dummyFile := filepath.Join(srcDir, "README.md")
	if err := os.WriteFile(dummyFile, []byte("# Test Repo"), 0o644); err != nil {
		t.Fatalf("failed to write dummy file: %v", err)
	}

	runGit(t, srcDir, "add", "README.md")
	runGit(t, srcDir, "commit", "-m", "Initial commit")

	cleanup := func() {
		os.RemoveAll(srcDir)
	}

	return srcDir, cleanup
}

// runGit is a helper to execute git commands in a specific directory
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, string(out))
	}
}

func TestGitClone(t *testing.T) {
	srcDir, cleanup := setupTestRepo(t)
	defer cleanup()

	destDir, err := os.MkdirTemp("", "ryvst-test-dest-*")
	if err != nil {
		t.Fatalf("failed to create temp dest dir: %v", err)
	}
	// Remove the dest dir so gitClone can create it fresh
	os.RemoveAll(destDir)

	err = gitClone(srcDir, destDir)
	if err != nil {
		t.Fatalf("gitClone failed: %v", err)
	}
	defer os.RemoveAll(destDir)

	readmePath := filepath.Join(destDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Errorf("expected README.md to exist in cloned repo, but it does not")
	}
}

func TestGitPull(t *testing.T) {
	// Save original working directory because gitPull uses os.Chdir
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(origDir) // Restore after test to prevent polluting other tests

	srcDir, cleanup := setupTestRepo(t)
	defer cleanup()

	destDir, err := os.MkdirTemp("", "ryvst-test-dest-*")
	if err != nil {
		t.Fatalf("failed to create temp dest dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	err = gitClone(srcDir, destDir)
	if err != nil {
		t.Fatalf("gitClone failed: %v", err)
	}

	err = gitPull(srcDir, destDir)
	if err != nil {
		t.Fatalf("gitPull failed: %v", err)
	}
}

func TestGit(t *testing.T) {
	// Save original working directory because Git -> gitPull uses os.Chdir
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(origDir)

	srcDir, cleanup := setupTestRepo(t)
	defer cleanup()

	destDir, err := os.MkdirTemp("", "ryvst-test-dest-*")
	if err != nil {
		t.Fatalf("failed to create temp dest dir: %v", err)
	}

	// Test 1: Destination does not exist -> should clone
	os.RemoveAll(destDir)

	// Note: Git() calls os.Exit(1) on failure.
	// We ensure the operation succeeds here to avoid terminating the test runner.
	Git(srcDir, destDir)

	readmePath := filepath.Join(destDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Errorf("expected README.md to exist after Git() clone, but it does not")
	}

	// Test 2: Destination exists -> should pull
	// Add a new commit to srcDir
	runGit(t, srcDir, "config", "user.email", "test@example.com")
	runGit(t, srcDir, "config", "user.name", "Test User")
	dummyFile2 := filepath.Join(srcDir, "UPDATE.md")
	if err := os.WriteFile(dummyFile2, []byte("# Update"), 0o644); err != nil {
		t.Fatalf("failed to write dummy file: %v", err)
	}
	runGit(t, srcDir, "add", "UPDATE.md")
	runGit(t, srcDir, "commit", "-m", "Second commit")

	// Call Git again, it should pull the new commit
	Git(srcDir, destDir)

	updatePath := filepath.Join(destDir, "UPDATE.md")
	if _, err := os.Stat(updatePath); os.IsNotExist(err) {
		t.Errorf("expected UPDATE.md to exist after Git() pull, but it does not")
	}

	os.RemoveAll(destDir)
}

func TestGitPullError(t *testing.T) {
	// Test gitPull directly with a non-existent directory to verify error handling
	err := gitPull("http://example.com/repo.git", "/non/existent/directory")
	if err == nil {
		t.Errorf("expected error when pulling into non-existent directory, got nil")
	}
}

func TestGitCloneError(t *testing.T) {
	// Test gitClone directly with an invalid repo URL to verify error handling
	err := gitClone("/invalid/repo/path", "/tmp/dest")
	if err == nil {
		t.Errorf("expected error when cloning from invalid path, got nil")
	}
}
