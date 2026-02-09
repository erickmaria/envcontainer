package cmd

import (
	"fmt"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/pkg/updater"
	"github.com/ErickMaria/envcontainer/internal/pkg/version"
	"github.com/spf13/cobra"
)

type versionOptions struct {
	verbose bool
}

func versionCommand(projOpts projectOptions) *cobra.Command {
	opts := versionOptions{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show envcontainer version",
		Long:  "Display the currently installed version and check for updates.",
		Run: func(cmd *cobra.Command, args []string) {
			opts.run()
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Show changelog for available updates")

	return cmd
}

func (opts versionOptions) run() {
	current := version.Version
	fmt.Println("Version: " + current)

	// Check for updates (from cache only)
	info, err := updater.CheckForUpdate(current)
	if err != nil || !info.HasUpdate {
		return
	}

	// Show available update
	fmt.Printf("\n⚠ Update available: %s → %s\n", current, info.LatestVersion)

	// Show changelog if requested (from cache)
	if opts.verbose && len(info.Changelogs) > 0 {
		showVersionChangelog(info.Changelogs)
	}
}

func showVersionChangelog(logs []updater.ReleaseChangelog) {
	fmt.Println("\nChanges:")
	for _, log := range logs {
		fmt.Printf("  v%s\n", log.Version)
		if log.Body != "" {
			for _, line := range strings.Split(log.Body, "\n") {
				if line := strings.TrimSpace(line); line != "" {
					fmt.Printf("    %s\n", line)
				}
			}
		}
	}
}
