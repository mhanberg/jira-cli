package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List lists boards in a project",
		Long:    "List lists boards in a project.",
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
	project := viper.GetString("project.key")

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

	boards, total, err := func() ([]*jira.Board, int, error) {
		s := cmdutil.Info(fmt.Sprintf("Fetching boards in project %s...", project))
		defer s.Stop()

		resp, err := api.DefaultClient(debug).Boards(project, jira.BoardTypeAll)
		if err != nil {
			return nil, 0, err
		}
		return resp.Boards, resp.Total, nil
	}()
	cmdutil.ExitIfError(err)

	// Total results in jira API response may not be present in older versions.
	if total == 0 {
		total = len(boards)
	}

	if total == 0 {
		fmt.Println()
		cmdutil.Failed("No boards found in project %q", project)
		return
	}

	if raw {
		out, err := json.MarshalIndent(boards, "", "  ")
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	v := view.NewBoard(
		boards,
		view.WithBoardPlain(plain),
		view.WithBoardNoHeaders(noHeaders),
		view.WithBoardCSV(csvOut),
		view.WithBoardDelimiter(delimiter),
	)

	cmdutil.ExitIfError(v.Render())
}
