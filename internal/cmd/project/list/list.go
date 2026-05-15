package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List lists Jira projects",
		Long:    "List lists Jira projects that a user has access to.",
		Aliases: []string{"lists", "ls"},
		Run:     List,
	}

	cmd.Flags().Bool("plain", false, "Display output in plain mode")
	cmd.Flags().Bool("no-headers", false, "Don't display table headers in plain mode. Works only with --plain")
	cmd.Flags().String("delimiter", "\t", "Custom delimiter for columns in plain mode. Works only with --plain")
	cmd.Flags().Bool("raw", false, "Print raw JSON output")
	cmd.Flags().Bool("csv", false, "Print output in CSV format")

	return cmd
}

// List displays a list view.
func List(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	plain, err := cmd.Flags().GetBool("plain")
	cmdutil.ExitIfError(err)

	noHeaders, err := cmd.Flags().GetBool("no-headers")
	cmdutil.ExitIfError(err)

	delimiter, err := cmd.Flags().GetString("delimiter")
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool("raw")
	cmdutil.ExitIfError(err)

	csvOut, err := cmd.Flags().GetBool("csv")
	cmdutil.ExitIfError(err)

	projects, total, err := func() ([]*jira.Project, int, error) {
		s := cmdutil.Info("Fetching projects...")
		defer s.Stop()

		projects, err := api.DefaultClient(debug).Project()
		if err != nil {
			return nil, 0, err
		}
		return projects, len(projects), nil
	}()
	cmdutil.ExitIfError(err)

	if total == 0 {
		cmdutil.Failed("No projects found.")
		return
	}

	if raw {
		out, err := json.MarshalIndent(projects, "", "  ")
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	v := view.NewProject(
		projects,
		view.WithProjectPlain(plain),
		view.WithProjectNoHeaders(noHeaders),
		view.WithProjectCSV(csvOut),
		view.WithProjectDelimiter(delimiter),
	)

	cmdutil.ExitIfError(v.Render())
}
