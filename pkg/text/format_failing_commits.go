package text

import (
	"fmt"
	"os"
	"strings"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
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
	markdownTable := FailedCommitsMarkdownTable(commits)

	// Write markdown table to GITHUB_STEP_SUMMARY
	writeToGitHubStepSummary(markdownTable)

	return t
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

// FailedCommitsMarkdownTable generates a markdown table from the provided commits
func FailedCommitsMarkdownTable(commits []FailingCommit) string {
	var tableData [][]string
	for _, commit := range commits {
		var row = []string{
			commit.Hash,
			commit.Error.Error(),
			strings.ReplaceAll(strings.TrimSpace(commit.Message), "\n", "<br>"),
		}
		tableData = append(tableData, row)
	}

	prettyPrintedTable, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("hash", "failure", "text").
		Format(tableData)

	if err != nil {
		// ... do your thing
	}
	return prettyPrintedTable
}
