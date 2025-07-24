package cmd

// build.go - Instructor Command
//
// This command is used by instructors to prepare challenge folders for distribution.
//
// Workflow:
// 1. Validates the folder structure to ensure it contains properly defined levels and check scripts.
// 2. Packages the folder into a compressed archive (e.g., .tar.gz).
// 3. Encrypts the archive with a starter key and an unlock time.
//    - The starter key will be given to students at test start time.
//    - The archive cannot be decrypted before the unlock time.
// 4. Produces an encrypted archive file ready to share with students.
//
// Flags:
// -f, --folder        Path to the challenge folder to be packaged (required).
// -p, --password 		 Password used for encrypt the archive (required)
// -k, --key           Starter key used for encryption (required).
// -u, --unlock-time   Unlock time in human-friendly format: "YYYY-MM-DD HH:MM" (24-hour clock) (required).
// -o, --output        Path to save the encrypted archive (optional, default: eval-archive.enc)
//
// Usage Example:
// qo build -f ./challenges -p foo -k bar -u "2025-07-10 09:30" -o ./test.enc

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ahmedYasserM/qo/pkg/archive"
	"github.com/ahmedYasserM/qo/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	folderPath       string
	password         string
	utKey            string
	unlockTime       string
	outputArchiveDir string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Package and encrypt a challenge folder for student testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		layout := "2006-01-02 15:04"
		_, err := time.Parse(layout, unlockTime)
		if err != nil {
			return err
		}

		if err := archive.IsValidFolderStructure(folderPath); err != nil {
			return err
		}

		if err := archive.CreateEncryptedTarArchive(folderPath, outputArchiveDir, unlockTime, password, utKey); err != nil {
			logger.Error(err)
			return err
		}

		logger.Success(fmt.Sprintf("%s folder is archived and encrypted successfully.", filepath.Clean(outputArchiveDir)))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Flags
	buildCmd.Flags().StringVarP(&folderPath, "folder", "f", "", "Path to the challenge folder to be packaged (required)")
	buildCmd.Flags().StringVarP(&password, "password", "p", "", "Password used for encrypt the archive (required)")
	buildCmd.Flags().StringVarP(&utKey, "key", "k", "", "Starter key used for encryption (required)")
	buildCmd.Flags().StringVarP(&unlockTime, "unlock-time", "u", "", "Unlock time in human-friendly format: \"YYYY-MM-DD HH:MM\" (24-hour clock) (required)")
	buildCmd.Flags().StringVarP(&outputArchiveDir, "output", "o", "eval-archive.enc", "Path to save the encrypted archive")

	buildCmd.MarkFlagRequired("folder")
	buildCmd.MarkFlagRequired("password")
	buildCmd.MarkFlagRequired("key")
	buildCmd.MarkFlagRequired("unlock-time")
}
