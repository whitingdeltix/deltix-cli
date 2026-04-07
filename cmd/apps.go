package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/whitingdeltix/deltix-cli/internal/api"
	"github.com/whitingdeltix/deltix-cli/internal/config"
	"github.com/whitingdeltix/deltix-cli/internal/output"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "List your registered apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		if apiURL != "" {
			cfg.APIURL = apiURL
		}
		if !cfg.IsLoggedIn() {
			return fmt.Errorf("not logged in — run: uxvalidator login")
		}

		client := api.NewClient(cfg)
		apps, err := client.ListApps()
		if err != nil {
			return err
		}

		if len(apps) == 0 {
			output.Info("No apps registered. Create one at the dashboard.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, output.Dim+"ID\tNAME\tBUNDLE\tTASKS\tSCORE"+output.Reset)
		for _, app := range apps {
			score := "—"
			if app.LastAggregateScore != nil {
				score = fmt.Sprintf("%s%d%s", output.ScoreColor(*app.LastAggregateScore), *app.LastAggregateScore, output.Reset)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				app.ID[:8], app.Name, app.BundleID, app.TaskCount, score)
		}
		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
}
