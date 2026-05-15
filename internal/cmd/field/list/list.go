package list

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdList is the field list command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List lists fields configured in the Jira instance",
		Long:    "List lists fields (system + custom) configured in the Jira instance, including their IDs.",
		Aliases: []string{"lists", "ls"},
		Run:     List,
	}

	cmd.Flags().Bool("plain", false, "Display output in plain mode")
	cmd.Flags().Bool("no-headers", false, "Don't display table headers in plain mode. Works only with --plain")
	cmd.Flags().String("delimiter", "\t", "Custom delimiter for columns in plain mode. Works only with --plain")
	cmd.Flags().Bool("raw", false, "Print raw JSON output")
	cmd.Flags().Bool("csv", false, "Print output in CSV format")
	cmd.Flags().Bool("custom-only", false, "Only show custom fields")

	return cmd
}

// List runs the field list command.
func List(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	customOnly, err := cmd.Flags().GetBool("custom-only")
	cmdutil.ExitIfError(err)

	opts := readListOpts(cmd)

	fields, err := func() ([]*jira.Field, error) {
		s := cmdutil.Info("Fetching fields...")
		defer s.Stop()
		return api.DefaultClient(debug).GetField()
	}()
	cmdutil.ExitIfError(err)

	if customOnly {
		filtered := fields[:0]
		for _, f := range fields {
			if f.Custom {
				filtered = append(filtered, f)
			}
		}
		fields = filtered
	}

	if len(fields) == 0 {
		cmdutil.Failed("No fields found.")
		return
	}

	if raw {
		cmdutil.ExitIfError(cmdutil.RenderJSON(fields))
		return
	}

	headers := []string{"ID", "NAME", "CUSTOM", "TYPE", "ITEMS"}
	rows := make([][]string, 0, len(fields))
	for _, f := range fields {
		rows = append(rows, []string{
			f.ID,
			f.Name,
			fmt.Sprintf("%v", f.Custom),
			f.Schema.DataType,
			f.Schema.Items,
		})
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
