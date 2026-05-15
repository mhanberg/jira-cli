package list

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/adf"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/md"
)

// NewCmdList is the worklog list command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list ISSUE-KEY",
		Short:   "List lists worklogs on an issue",
		Long:    "List lists worklogs (timelog entries) on an issue, newest first.",
		Aliases: []string{"lists", "ls"},
		Args:    cobra.MinimumNArgs(1),
		Annotations: map[string]string{
			"help:args": "ISSUE-KEY\tIssue key, eg: ISSUE-1",
		},
		Run: List,
	}

	cmd.Flags().Bool("plain", false, "Display output in plain mode")
	cmd.Flags().Bool("no-headers", false, "Don't display table headers in plain mode. Works only with --plain")
	cmd.Flags().String("delimiter", "\t", "Custom delimiter for columns in plain mode. Works only with --plain")
	cmd.Flags().Bool("raw", false, "Print raw JSON output")
	cmd.Flags().Bool("csv", false, "Print output in CSV format")
	cmd.Flags().Uint("limit", 0, "Maximum number of worklogs to show (0 = all)")

	return cmd
}

// List runs the worklog list command.
func List(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	limit, err := cmd.Flags().GetUint("limit")
	cmdutil.ExitIfError(err)

	opts := readListOpts(cmd)
	key := cmdutil.GetJiraIssueKey(viper.GetString("project.key"), args[0])

	resp, err := func() (*jira.WorklogResult, error) {
		s := cmdutil.Info("Fetching worklogs...")
		defer s.Stop()
		return api.DefaultClient(debug).IssueWorklogs(key)
	}()
	cmdutil.ExitIfError(err)

	if len(resp.Worklogs) == 0 {
		cmdutil.Failed("No worklogs found on %s.", key)
		return
	}

	count := len(resp.Worklogs)
	if limit > 0 && uint(count) > limit {
		count = int(limit)
	}

	// Newest first.
	picked := make([]*jira.Worklog, 0, count)
	for i := len(resp.Worklogs) - 1; i >= len(resp.Worklogs)-count; i-- {
		picked = append(picked, resp.Worklogs[i])
	}

	if raw {
		cmdutil.ExitIfError(cmdutil.RenderJSON(picked))
		return
	}

	headers := []string{"ID", "AUTHOR", "STARTED", "TIME SPENT", "COMMENT"}
	rows := make([][]string, 0, len(picked))
	for _, w := range picked {
		author := w.Author.DisplayName
		if author == "" {
			author = w.Author.Name
		}
		rows = append(rows, []string{
			w.ID,
			author,
			cmdutil.FormatDateTimeHuman(w.Started, jira.RFC3339),
			w.TimeSpent,
			oneLine(worklogComment(w.Comment), opts.Plain),
		})
	}

	cmdutil.ExitIfError(cmdutil.RenderTable(headers, rows, opts))
}

func worklogComment(v interface{}) string {
	if adfNode, ok := v.(*adf.ADF); ok && adfNode != nil {
		return adf.NewTranslator(adfNode, adf.NewMarkdownTranslator()).Translate()
	}
	if s, ok := v.(string); ok {
		return md.FromJiraMD(s)
	}
	return ""
}

func oneLine(s string, plain bool) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	if !plain && len(s) > 60 {
		s = s[:57] + "..."
	}
	return s
}

func readListOpts(cmd *cobra.Command) cmdutil.TableRenderOptions {
	plain, _ := cmd.Flags().GetBool("plain")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	delim, _ := cmd.Flags().GetString("delimiter")
	csv, _ := cmd.Flags().GetBool("csv")

	return cmdutil.TableRenderOptions{
		Plain:     plain,
		NoHeaders: noHeaders,
		CSV:       csv,
		Delimiter: delim,
	}
}
