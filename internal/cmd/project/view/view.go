package view

import (
	"encoding/json"
	"fmt"
	"os"
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

// NewCmdView is the project view command.
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view PROJECT-KEY",
		Short: "View displays details of a project",
		Long:  "View displays details of a project.",
		Example: `$ jira project view MYPROJ

# Open the project in the default web browser
$ jira project view MYPROJ --web

# Print the raw JSON response
$ jira project view MYPROJ --raw`,
		Aliases: []string{"show"},
		Annotations: map[string]string{
			"help:args": "PROJECT-KEY\tProject key or id, eg: MYPROJ",
		},
		Args: cobra.MinimumNArgs(1),
		Run:  view,
	}

	cmd.Flags().BoolP(flagWeb, "w", false, "Open the project in the default web browser")
	cmd.Flags().Bool(flagRaw, false, "Print the raw Jira API response")

	return cmd
}

func view(cmd *cobra.Command, args []string) {
	key := args[0]

	debug, err := cmd.Flags().GetBool(flagDebug)
	cmdutil.ExitIfError(err)

	web, err := cmd.Flags().GetBool(flagWeb)
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool(flagRaw)
	cmdutil.ExitIfError(err)

	server := viper.GetString(configServer)

	if web {
		url := cmdutil.GenerateServerBrowseURL(server, key)
		fmt.Println(url)
		cmdutil.ExitIfError(browser.Browse(url))
		return
	}

	client := api.DefaultClient(debug)
	project, err := func() (*jira.Project, error) {
		s := cmdutil.Info("Fetching project details...")
		defer s.Stop()
		return client.GetProject(key)
	}()
	cmdutil.ExitIfError(err)

	if raw {
		out, err := json.MarshalIndent(project, "", "  ")
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	render(project)
}

func render(p *jira.Project) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "Key:\t%s\n", p.Key)
	fmt.Fprintf(w, "Name:\t%s\n", p.Name)
	if p.Type != "" {
		fmt.Fprintf(w, "Type:\t%s\n", p.Type)
	}
	if p.Lead.Name != "" {
		fmt.Fprintf(w, "Lead:\t%s\n", p.Lead.Name)
	}
	_ = w.Flush()
}
