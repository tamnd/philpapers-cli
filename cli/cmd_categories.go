package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/philpapers-cli/philpapers"
)

func (a *App) categoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List available philosophy categories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cats := philpapers.Categories()
			return a.render(cats)
		},
	}
}
