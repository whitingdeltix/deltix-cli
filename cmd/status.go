package cmd

import (
	"fmt"

	"github.com/whitingdeltix/deltix-cli/internal/api"
	"github.com/whitingdeltix/deltix-cli/internal/config"
	"github.com/whitingdeltix/deltix-cli/internal/output"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check run status",
	RunE: func(cmd *cobra.Command, args []string) error {
		runID, _ := cmd.Flags().GetString("run")
		if runID == "" {
			return fmt.Errorf("--run is required")
		}

		cfg, _ := config.Load()
		if apiURL != "" {
			cfg.APIURL = apiURL
		}
		if !cfg.IsLoggedIn() {
			return fmt.Errorf("not logged in — run: uxvalidator login")
		}

		client := api.NewClient(cfg)
		run, err := client.GetRun(runID)
		if err != nil {
			return err
		}

		fmt.Printf("Run %s%s%s\n", output.Bold, run.ID[:8], output.Reset)
		fmt.Printf("Status:    %s\n", output.Status(run.Status))
		fmt.Printf("Tasks:     %d/%d completed\n", run.CompletedCount, run.TaskCount)
		fmt.Printf("Device:    %s\n", run.DeviceType)
		fmt.Printf("Created:   %s\n", run.CreatedAt.Format("Jan 2, 3:04 PM"))
		if run.AggregateScore != nil {
			fmt.Printf("Score:     %s%d%s/100\n", output.ScoreColor(*run.AggregateScore), *run.AggregateScore, output.Reset)
		}
		return nil
	},
}

func init() {
	statusCmd.Flags().String("run", "", "Run ID (required)")
	rootCmd.AddCommand(statusCmd)
}
