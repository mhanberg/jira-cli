package list

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdList is the issue link list command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List lists issue link types configured in the Jira instance",
		Long:    "List lists issue link types — the relationship names you can use as the third argument to `jira issue link`.",
		Aliases: []string{"lists", "ls", "types"},
		Run:     List,
	}

	cmd.Flags().Bool("plain", false, "Display output in plain mode")
	cmd.Flags().Bool("no-headers", false, "Don't display table headers in plain mode. Works only with --plain")
	cmd.Flags().String("delimiter", "\t", "Custom delimiter for columns in plain mode. Works only with --plain")
	cmd.Flags().Bool("raw", false, "Print raw JSON output")
	cmd.Flags().Bool("csv", false, "Print output in CSV format")

	return cmd
}

// List runs the issue link list command.
func List(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	opts := readListOpts(cmd)

	types, err := func() ([]*jira.IssueLinkType, error) {
		s := cmdutil.Info("Fetching link types...")
		defer s.Stop()
		return api.DefaultClient(debug).GetIssueLinkTypes()
	}()
	cmdutil.ExitIfError(err)

	if len(types) == 0 {
		cmdutil.Failed("No link types found.")
		return
	}

	if raw {
		cmdutil.ExitIfError(cmdutil.RenderJSON(types))
		return
	}

	headers := []string{"ID", "NAME", "INWARD", "OUTWARD"}
	rows := make([][]string, 0, len(types))
	for _, t := range types {
		rows = append(rows, []string{t.ID, t.Name, t.Inward, t.Outward})
	}

	cmdutil.ExitIfError(cmdutil.RenderTable(headers, rows, opts))
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
