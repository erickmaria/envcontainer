package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/pkg/updater"
	"github.com/ErickMaria/envcontainer/internal/pkg/version"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	force      bool
	verbose    bool
	forceCheck bool
}

func updateCommand(projOpts projectOptions) *cobra.Command {
	opts := updateOptions{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update envcontainer to the latest version",
		Long:  "Check for a new version and update envcontainer to the latest available release.",
		Run: func(cmd *cobra.Command, args []string) {
			opts.execute()
		},
	}

	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "Force update without confirmation")
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Show changelog of each version being updated")
	cmd.Flags().BoolVarP(&opts.forceCheck, "check", "c", false, "Force check for updates, bypass cache")

	return cmd
}

func (opts updateOptions) execute() {
	current := version.Version

	// If only checking, don't perform update
	if opts.forceCheck {
		checkOnly(current, opts.verbose)
		return
	}

	fmt.Println("Checking for updates...")

	// Check for new version
	info, err := updater.CheckForUpdate(current)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		os.Exit(1)
	}

	// Already up to date
	if !info.HasUpdate {
		fmt.Println("✓ You are on the latest version (" + current + ")")
		return
	}

	// Show new version available
	fmt.Printf("✓ New version available: %s → %s\n", current, info.LatestVersion)

	// Show changelog if requested
	if opts.verbose && len(info.Changelogs) > 0 {
		showChangelogFromInfo(info.Changelogs)
	}

	// Confirm update
	if !opts.force && !confirmUpdate() {
		fmt.Println("Update cancelled.")
		return
	}

	// Download and install
	if err := downloadAndReplace(info.LatestVersion); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Update complete!")

	// Clear cache
	_ = updater.ClearCache()
}

func checkOnly(current string, verbose bool) {
	fmt.Println("Checking for updates...")

	info, err := updater.CheckForUpdateForce(current, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		os.Exit(1)
	}

	if !info.HasUpdate {
		fmt.Println("✓ You are on the latest version (" + current + ")")
		return
	}

	fmt.Printf("✓ New version available: %s → %s\n", current, info.LatestVersion)
	if verbose && len(info.Changelogs) > 0 {
		showChangelogFromInfo(info.Changelogs)
	}
}

func confirmUpdate() bool {
	fmt.Print("Update now? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}

func showChangelogFromInfo(logs []updater.ReleaseChangelog) {
	fmt.Println("\nChangelog:")
	for _, log := range logs {
		fmt.Printf("  v%s\n", log.Version)
		if log.Body != "" {
			for _, line := range strings.Split(log.Body, "\n") {
				if strings.TrimSpace(line) != "" {
					fmt.Printf("    %s\n", line)
				}
			}
		}
	}
	fmt.Println()
}


func downloadAndReplace(newVersion string) error {
	// Get release info
	url := "https://api.github.com/repos/ErickMaria/envcontainer/releases/tags/v" + newVersion
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("release not found")
	}

	var release updater.Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("parse release: %w", err)
	}

	// Find download URL for current OS
	downloadURL := findAssetURL(release.Assets)
	if downloadURL == "" {
		return fmt.Errorf("no binary for %s", runtime.GOOS)
	}

	// Get current binary path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get binary path: %w", err)
	}

	// Download
	fmt.Println("Updating...")
	tmpZip, err := downloadToTemp(downloadURL)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer os.Remove(tmpZip)

	// Extract and replace
	return replaceExecutable(tmpZip, exePath)
}

func findAssetURL(assets []updater.Asset) string {
	osName := runtime.GOOS
	for _, asset := range assets {
		if strings.Contains(strings.ToLower(asset.Name), osName) {
			return asset.DownloadURL
		}
	}
	// Fallback to linux
	for _, asset := range assets {
		if strings.Contains(strings.ToLower(asset.Name), "linux") {
			return asset.DownloadURL
		}
	}
	return ""
}

func downloadToTemp(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "envcontainer-*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	return tmpFile.Name(), err
}

func replaceExecutable(zipPath, exePath string) error {
	// Extract zip
	tmpDir, err := os.MkdirTemp("", "envcontainer-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("unzip", "-o", zipPath, "-d", tmpDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	// Find binary
	var binPath string
	filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.Contains(info.Name(), "envcontainer") {
			binPath = path
			return filepath.SkipDir
		}
		return nil
	})

	if binPath == "" {
		return fmt.Errorf("binary not in archive")
	}

	// Replace: backup old → copy new
	backupPath := exePath + ".backup"
	newPath := exePath + ".new"

	if err := copyFile(binPath, newPath); err != nil {
		return err
	}

	if err := os.Chmod(newPath, 0o755); err != nil {
		os.Remove(newPath)
		return err
	}

	if err := os.Rename(exePath, backupPath); err != nil {
		os.Remove(newPath)
		return err
	}

	if err := os.Rename(newPath, exePath); err != nil {
		os.Rename(backupPath, exePath)
		return err
	}

	os.Remove(backupPath)
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
