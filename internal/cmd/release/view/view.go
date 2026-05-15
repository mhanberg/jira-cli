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

// NewCmdView is the release view command.
func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view VERSION-ID",
		Short: "View displays details of a project version (release)",
		Long:  "View displays details of a project version (release).",
		Example: `$ jira release view 10042

# Open the release in the default web browser
$ jira release view 10042 --web

# Print the raw JSON response
$ jira release view 10042 --raw`,
		Aliases: []string{"show"},
		Annotations: map[string]string{
			"help:args": "VERSION-ID\tID of the version, eg: 10042",
		},
		Args: cobra.MinimumNArgs(1),
		Run:  view,
	}

	cmd.Flags().BoolP(flagWeb, "w", false, "Open the release in the default web browser")
	cmd.Flags().Bool(flagRaw, false, "Print the raw Jira API response")

	return cmd
}

func view(cmd *cobra.Command, args []string) {
	id := args[0]

	debug, err := cmd.Flags().GetBool(flagDebug)
	cmdutil.ExitIfError(err)

	web, err := cmd.Flags().GetBool(flagWeb)
	cmdutil.ExitIfError(err)

	raw, err := cmd.Flags().GetBool(flagRaw)
	cmdutil.ExitIfError(err)

	client := api.DefaultClient(debug)
	version, err := func() (*jira.ProjectVersion, error) {
		s := cmdutil.Info("Fetching release details...")
		defer s.Stop()
		return client.GetVersion(id)
	}()
	cmdutil.ExitIfError(err)

	if web {
		url := fmt.Sprintf("%s/projects/%s/versions/%s", viper.GetString(configServer), viper.GetString("project.key"), id)
		fmt.Println(url)
		cmdutil.ExitIfError(browser.Browse(url))
		return
	}

	if raw {
		out, err := json.MarshalIndent(version, "", "  ")
		cmdutil.ExitIfError(err)
		fmt.Println(string(out))
		return
	}

	render(version)
}

func render(v *jira.ProjectVersion) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", v.ID)
	fmt.Fprintf(w, "Name:\t%s\n", v.Name)
	fmt.Fprintf(w, "Released:\t%v\n", v.Released)
	fmt.Fprintf(w, "Archived:\t%v\n", v.Archived)
	if v.Description != nil {
		fmt.Fprintf(w, "Description:\t%v\n", v.Description)
	}
	_ = w.Flush()
}
