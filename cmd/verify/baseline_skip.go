package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var baselineSkipCmd = &cobra.Command{
	Use:   "baseline-skip",
	Short: "Manage and verify UI baseline artifact-hash skips",
	RunE: func(cmd *cobra.Command, args []string) error {
		write, _ := cmd.Flags().GetBool("write")

		if write {
			return writeBaselineHash()
		}

		return checkBaselineSkip()
	},
}

func addFilePath(paths *[]string, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		*paths = append(*paths, path)
	}
	return nil
}

func gatherFiles(dirs []string) ([]string, error) {
	var paths []string
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			return addFilePath(&paths, path, info, err)
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func hashFiles(paths []string) (string, error) {
	h := sha256.New()
	for _, path := range paths {
		h.Write([]byte(path))

		// #nosec G304 -- local maintenance utility
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(h, file)
		_ = file.Close()
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func computeSourceHash() (string, error) {
	paths, err := gatherFiles([]string{"ui", "tests/e2e"})
	if err != nil {
		return "", err
	}
	return hashFiles(paths)
}

func writeBaselineHash() error {
	hash, err := computeSourceHash()
	if err != nil {
		return fmt.Errorf("failed to compute source hash: %w", err)
	}

	dir := ".tester"
	// #nosec G301 -- directory permissions required for local testing
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	hashFile := filepath.Join(dir, "last_verified_hash")
	// #nosec G304 -- local maintenance utility
	if err := os.WriteFile(hashFile, []byte(hash), 0600); err != nil {
		return fmt.Errorf("failed to write hash file: %w", err)
	}

	fmt.Printf("✅ Baseline source hash successfully written: %s\n", hash)
	return nil
}

func checkBaselineSkip() error {
	// 1. Check if git status --porcelain is empty
	err := maintenance.VerifyGitClean(".")
	if err != nil {
		return fmt.Errorf("git state is dirty: %w", err)
	}

	// 2. Read last_verified_hash
	hashFile := filepath.Join(".tester", "last_verified_hash")
	// #nosec G304 -- local maintenance utility
	storedBytes, err := os.ReadFile(hashFile)
	if err != nil {
		return fmt.Errorf("last_verified_hash not found: %w", err)
	}
	storedHash := strings.TrimSpace(string(storedBytes))

	// 3. Compute current hash
	currentHash, err := computeSourceHash()
	if err != nil {
		return fmt.Errorf("failed to compute current source hash: %w", err)
	}

	if currentHash != storedHash {
		return fmt.Errorf("hash mismatch: current %s != stored %s", currentHash, storedHash)
	}

	fmt.Println("✅ UI baseline matches last verified state. Skip permitted.")
	return nil
}

func init() {
	baselineSkipCmd.Flags().Bool("write", false, "Write current source hash to last_verified_hash")
}
