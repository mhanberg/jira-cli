package view

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

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

// NewCmdView is the sprint view command.
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view SPRINT-ID",
		Short: "View displays details of a sprint",
		Long:  "View displays details of a sprint.",
		Example: `$ jira sprint view 123

# Open the sprint's board in the default web browser
$ jira sprint view 123 --web

# Print the raw JSON response
$ jira sprint view 123 --raw`,
		Aliases: []string{"show"},
		Annotations: map[string]string{
			"help:args": "SPRINT-ID\tID of the sprint, eg: 123",
		},
		Args: cobra.MinimumNArgs(1),
		Run:  view,
	}

	cmd.Flags().BoolP(flagWeb, "w", false, "Open the sprint's board in the default web browser")
	cmd.Flags().Bool(flagRaw, false, "Print the raw Jira API response")

	return cmd
}

func view(cmd *cobra.Command, args []string) {
	sprintID, err := strconv.Atoi(args[0])
	if err != nil {
		cmdutil.Failed("Invalid sprint id %q: must be an integer", args[0])
		return
	}

	debug, err := cmd.Flags().GetBool(flagDebug)
	cmdutil.ExitIfError(err)

	web, err := cmd.Flags().GetBool(flagWeb)
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool(flagRaw)
	cmdutil.ExitIfError(err)

	client := api.DefaultClient(debug)

	sprint, err := func() (*jira.Sprint, error) {
		s := cmdutil.Info("Fetching sprint details...")
		defer s.Stop()
		return client.GetSprint(sprintID)
	}()
	cmdutil.ExitIfError(err)

	if web {
		url := sprintWebURL(viper.GetString(configServer), sprint)
		fmt.Println(url)
		cmdutil.ExitIfError(browser.Browse(url))
		return
	}

	if raw {
		out, err := json.MarshalIndent(sprint, "", "  ")
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	render(sprint)
}

func render(sprint *jira.Sprint) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%d\n", sprint.ID)
	fmt.Fprintf(w, "Name:\t%s\n", sprint.Name)
	fmt.Fprintf(w, "State:\t%s\n", sprint.Status)
	if sprint.BoardID != 0 {
		fmt.Fprintf(w, "Board ID:\t%d\n", sprint.BoardID)
	}
	if sprint.StartDate != "" {
		fmt.Fprintf(w, "Start:\t%s\n", cmdutil.FormatDateTimeHuman(sprint.StartDate, time.RFC3339))
	}
	if sprint.EndDate != "" {
		fmt.Fprintf(w, "End:\t%s\n", cmdutil.FormatDateTimeHuman(sprint.EndDate, time.RFC3339))
	}
	if sprint.CompleteDate != "" {
		fmt.Fprintf(w, "Completed:\t%s\n", cmdutil.FormatDateTimeHuman(sprint.CompleteDate, time.RFC3339))
	}
	_ = w.Flush()
}

func sprintWebURL(server string, sprint *jira.Sprint) string {
	if sprint.BoardID == 0 {
		return server
	}
	return fmt.Sprintf("%s/jira/software/c/projects/_/boards/%d?selectedSprint=%d", server, sprint.BoardID, sprint.ID)
}
