package view

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/browser"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	flagWeb   = "web"
	flagRaw   = "raw"
	flagDebug = "debug"

	configServer = "server"
)

// NewCmdView is the board view command.
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view BOARD-ID",
		Short: "View displays details of a board",
		Long:  "View displays details of a board.",
		Example: `$ jira board view 17

# Open the board in the default web browser
$ jira board view 17 --web

# Print the raw JSON response
$ jira board view 17 --raw`,
		Aliases: []string{"show"},
		Annotations: map[string]string{
			"help:args": "BOARD-ID\tID of the board, eg: 17",
		},
		Args: cobra.MinimumNArgs(1),
		Run:  view,
	}

	cmd.Flags().BoolP(flagWeb, "w", false, "Open the board in the default web browser")
	cmd.Flags().Bool(flagRaw, false, "Print the raw Jira API response")

	return cmd
}

func view(cmd *cobra.Command, args []string) {
	boardID, err := strconv.Atoi(args[0])
	if err != nil {
		cmdutil.Failed("Invalid board id %q: must be an integer", args[0])
		return
	}

	debug, err := cmd.Flags().GetBool(flagDebug)
	cmdutil.ExitIfError(err)

	web, err := cmd.Flags().GetBool(flagWeb)
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool(flagRaw)
	cmdutil.ExitIfError(err)

	server := viper.GetString(configServer)

	if web {
		url := fmt.Sprintf("%s/secure/RapidBoard.jspa?rapidView=%d", server, boardID)
		fmt.Println(url)
		cmdutil.ExitIfError(browser.Browse(url))
		return
	}

	client := api.DefaultClient(debug)
	board, err := func() (*jira.Board, error) {
		s := cmdutil.Info("Fetching board details...")
		defer s.Stop()
		return client.GetBoard(boardID)
	}()
	cmdutil.ExitIfError(err)

	if raw {
		out, err := json.MarshalIndent(board, "", "  ")
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	render(board)
}

func render(b *jira.Board) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%d\n", b.ID)
	fmt.Fprintf(w, "Name:\t%s\n", b.Name)
	if b.Type != "" {
		fmt.Fprintf(w, "Type:\t%s\n", b.Type)
	}
	_ = w.Flush()
}
