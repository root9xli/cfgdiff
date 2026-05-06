package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"cfgdiff/internal/audit"
	"cfgdiff/internal/diff"
	"cfgdiff/internal/output"
	"cfgdiff/internal/parser"
)

var (
	auditLog  string
	format    string
	summary   bool
	sinceFlag string
)

var rootCmd = &cobra.Command{
	Use:   "cfgdiff <file1> <file2>",
	Short: "Diff and audit config file changes across environments",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Print audit log entries",
	RunE:  runAudit,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&auditLog, "audit-log", "", "path to audit log file")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "output format: text or json")
	rootCmd.PersistentFlags().BoolVarP(&summary, "summary", "s", false, "show summary only")

	auditCmd.Flags().StringVar(&sinceFlag, "since", "", "filter entries since timestamp (RFC3339)")
	rootCmd.AddCommand(auditCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	a, err := parser.Parse(args[0])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[0], err)
	}
	b, err := parser.Parse(args[1])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[1], err)
	}

	changes := diff.Compare(a, b)

	fmt := output.NewFormatter(os.Stdout)
	if format == "json" {
		if err := fmt.WriteJSON(changes, summary); err != nil {
			return err
		}
	} else {
		if err := fmt.WriteText(changes, summary); err != nil {
			return err
		}
	}

	if auditLog != "" {
		logger, err := audit.NewLogger(auditLog)
		if err != nil {
			return fmt.Errorf("opening audit log: %w", err)
		}
		defer logger.Close()
		if err := logger.Record(args[0], args[1], changes); err != nil {
			return fmt.Errorf("writing audit log: %w", err)
		}
	}
	return nil
}

func runAudit(cmd *cobra.Command, args []string) error {
	if auditLog == "" {
		return fmt.Errorf("--audit-log is required for audit command")
	}
	var since time.Time
	if sinceFlag != "" {
		var err error
		since, err = time.Parse(time.RFC3339, sinceFlag)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
	}
	entries, err := audit.ReadLog(auditLog, since)
	if err != nil {
		return err
	}
	audit.PrintEntries(os.Stdout, entries)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
