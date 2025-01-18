package text

import (
	"strings"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

// FailingCommit is just a formatted commit struct
type FailingCommit struct {
	Hash    string
	Message string
	Error   error
}

// FormatFailingCommits takes in slice of commit hashes and messages and formats it for nice output
func FormatFailingCommits(commits []FailingCommit) table.Writer {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"hash", "failure", "text"})

	builder := strings.Builder{}
	// Extra spacing to make it nicer
	builder.WriteString("\nFollowing commits failed the check: \n")

	for _, commit := range commits {
		t.AppendRow(table.Row{commit.Hash, commit.Error.Error(), commit.Message})
	}

	// Generate markdown table
	markdownTable := generateMarkdownTable(commits)

	// Write markdown table to GITHUB_STEP_SUMMARY
	writeToGitHubStepSummary(markdownTable)

	return t
}

func generateMarkdownTable(commits []FailingCommit) string {
    var sb strings.Builder
    sb.WriteString("| hash | failure | text |\n")
    sb.WriteString("|------|---------|------|\n")
    for _, commit := range commits {
        sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", commit.Hash, commit.Error.Error(), strings.TrimSpace(commit.Message)))
    }
    return sb.String()
}

func writeToGitHubStepSummary(content string) {
    summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
    if summaryFile != "" {
        f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_WRONLY, 0600)
        if err != nil {
            fmt.Printf("Error opening GITHUB_STEP_SUMMARY file: %v\n", err)
            return
        }
        defer f.Close()

        if _, err = f.WriteString(content); err != nil {
            fmt.Printf("Error writing to GITHUB_STEP_SUMMARY file: %v\n", err)
        }
    }
}

