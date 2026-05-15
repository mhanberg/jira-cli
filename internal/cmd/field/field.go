package field

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/field/list"
)

const helpText = `Field manages Jira fields. See available commands below.`

// NewCmdField is the field command.
func NewCmdField() *cobra.Command {
	cmd := cobra.Command{
		Use:         "field",
		Short:       "Field manages Jira fields",
		Long:        helpText,
		Aliases:     []string{"fields"},
		Annotations: map[string]string{"cmd:main": "true"},
		RunE:        field,
	}

	cmd.AddCommand(list.NewCmdList())

	return &cmd
}

func field(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
