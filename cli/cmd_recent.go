package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) recentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recent",
		Short: "List recently added PhilPapers entries",
		RunE: func(cmd *cobra.Command, _ []string) error {
			limit := a.effectiveLimit(20)
			papers, err := a.client.Recent(cmd.Context(), limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(papers, len(papers))
		},
	}
}
