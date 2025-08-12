package cmd

// start.go - Student Command
//
// This command is used by students to start their test session.
//
// Workflow:
// 1. Prompts the student to enter their Student ID (used for reports and logs).
// 2. Verifies the provided starter key and unlock time to decrypt the archive.
//    - The test will not start before the scheduled unlock time.
// 3. Sets up a sandboxed environment using Linux namespaces (isolates processes, users, and filesystem).
// 4. Extracts the challenge folder into the sandbox and launches an interactive shell for the student.
// 5. Monitors activity and logs commands executed by the student.
// 6. When time ends or the student chooses to finish, generates a single-page PDF report with their results.
//
// Flags:
// -i  --id 			 	 	 Student ID (required)
// -a, --archive  		 Path to the encrypted archive file (required).
// -p, --password 		 Password used for encrypt the archive (required)
// -k, --key           Starter key used for encryption (required).
// -d, --duration      Total duration of the test in minutes (required).
// -o, --output        Directory to save logs and PDF report (optional, default: eval-results)
//
// Usage Example:
// eval start -a ./test.enc -p foo -k bar -d 1h30m -o ./results

import (
	"fmt"
	"os"
	"time"

	"github.com/ahmedYasserM/qo/pkg/archive"
	"github.com/ahmedYasserM/qo/pkg/logger"
	"github.com/ahmedYasserM/qo/pkg/sandbox"
	"github.com/spf13/cobra"
)

var (
	id            uint16
	archivePath   string
	utKeyStart    string
	passwordStart string
	testDuration  time.Duration
	outputLogDir  string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a test session in a sandboxed environment.",
	RunE: func(cmd *cobra.Command, args []string) error {

		if os.Geteuid() != 0 {
			logger.Error(fmt.Errorf("this program must be run as root"))
			os.Exit(1)
		}

		if err := sandbox.ExtractRootfs(); err != nil {
			return err
		}

		if err := archive.DecryptTarArchive(archivePath, passwordStart, utKeyStart); err != nil {
			return err
		}

		logger.Success(fmt.Sprintf("%s folder is unpacked and decrypted successfully.", archivePath))

		err := sandbox.StartSandBox()

		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Flags
	startCmd.Flags().Uint16VarP(&id, "id", "i", 0, "Student ID (required)")
	startCmd.Flags().StringVarP(&archivePath, "archive", "a", "", "Path to the encrypted archive file (required)")
	startCmd.Flags().StringVarP(&passwordStart, "password", "p", "", "Password used for encrypt the archive (required)")
	startCmd.Flags().StringVarP(&utKeyStart, "key", "k", "", "Starter key used for decryption (required)")
	startCmd.Flags().DurationVarP(&testDuration, "duration", "d", 0, "Total duration of the test (e.g., 90m, 1h30m) (required)")
	startCmd.Flags().StringVarP(&outputLogDir, "output", "o", "eval-results", "Output directory for logs and PDF reports")

	startCmd.MarkFlagRequired("id")
	startCmd.MarkFlagRequired("archive")
	startCmd.MarkFlagRequired("password")
	startCmd.MarkFlagRequired("key")
	startCmd.MarkFlagRequired("duration")

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
