package cmd

import (
	"fmt"
	"mainmodule/tools"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup unCommit files",
	Run: func(cmd *cobra.Command, args []string) {
		currentWorkDir, err := tools.GetGitRootDir()
		if err != nil {
			fmt.Println(err)
			return
		}
		fileName := filepath.Base(currentWorkDir)
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		backupTargetDir := fmt.Sprintf("%s/%s", userHomeDir, fileName)

		tools.StartBackupUnCommitFiles(currentWorkDir, backupTargetDir)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
