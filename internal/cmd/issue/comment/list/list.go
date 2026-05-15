package list

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/adf"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/jira/filter/issue"
	"github.com/ankitpokhrel/jira-cli/pkg/md"
)

// NewCmdList is the comment list command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list ISSUE-KEY",
		Short:   "List lists comments on an issue",
		Long:    "List lists comments on an issue, newest first.",
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
	cmd.Flags().Uint("limit", 0, "Maximum number of comments to show (0 = all)")

	return cmd
}

// List runs the comment list command.
func List(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	limit, err := cmd.Flags().GetUint("limit")
	cmdutil.ExitIfError(err)

	opts := readListOpts(cmd)
	key := cmdutil.GetJiraIssueKey(viper.GetString("project.key"), args[0])

	iss, err := func() (*jira.Issue, error) {
		s := cmdutil.Info("Fetching comments...")
		defer s.Stop()
		return api.ProxyGetIssue(api.DefaultClient(debug), key, issue.NewNumCommentsFilter(5000))
	}()
	cmdutil.ExitIfError(err)

	all := iss.Fields.Comment.Comments
	if len(all) == 0 {
		cmdutil.Failed("No comments found on %s.", key)
		return
	}

	// Newest first.
	count := len(all)
	if limit > 0 && uint(count) > limit {
		count = int(limit)
	}
	picked := make([]any, 0, count)
	for i := len(all) - 1; i >= len(all)-count; i-- {
		picked = append(picked, all[i])
	}

	if raw {
		cmdutil.ExitIfError(cmdutil.RenderJSON(picked))
		return
	}

	headers := []string{"ID", "AUTHOR", "CREATED", "BODY"}
	rows := make([][]string, 0, len(picked))
	for i := len(all) - 1; i >= len(all)-count; i-- {
		c := all[i]
		author := c.Author.DisplayName
		if author == "" {
			author = c.Author.Name
		}
		rows = append(rows, []string{
			c.ID,
			author,
			cmdutil.FormatDateTimeHuman(c.Created, jira.RFC3339),
			oneLine(commentBody(c.Body), opts.Plain),
		})
	}

	cmdutil.ExitIfError(cmdutil.RenderTable(headers, rows, opts))
}

func commentBody(v interface{}) string {
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
	if !plain && len(s) > 80 {
		s = s[:77] + "..."
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
