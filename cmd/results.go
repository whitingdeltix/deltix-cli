package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/whitingdeltix/deltix-cli/internal/api"
	"github.com/whitingdeltix/deltix-cli/internal/config"
	"github.com/whitingdeltix/deltix-cli/internal/output"
	"github.com/spf13/cobra"
)

var resultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Show run results",
	RunE: func(cmd *cobra.Command, args []string) error {
		runID, _ := cmd.Flags().GetString("run")
		jsonOut, _ := cmd.Flags().GetBool("json")
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
		results, err := client.GetRunResults(runID)
		if err != nil {
			return err
		}

		if jsonOut {
			data, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for _, r := range results {
			desc := "Task " + r.TaskID[:8]
			if r.TaskDescription != nil {
				desc = *r.TaskDescription
			}

			status := output.Green + "✓ pass" + output.Reset
			if !r.Succeeded {
				status = output.Red + "✗ fail" + output.Reset
			}

			fmt.Printf("\n%s  %s%s%s  %s%d%s/100  %d steps\n",
				status, output.Bold, desc, output.Reset,
				output.ScoreColor(r.AggregateScore), r.AggregateScore, output.Reset,
				r.Steps)

			// Score breakdown
			fmt.Printf("  Discoverability: %d  Efficiency: %d  Navigation: %d\n",
				r.Scores.Discoverability, r.Scores.Efficiency, r.Scores.NavigationClarity)
			fmt.Printf("  Feedback: %d  Confirmation: %d  Interruptions: %d\n",
				r.Scores.FeedbackClarity, r.Scores.ConfirmationClarity, r.Scores.InterruptionImpact)

			// Findings
			if len(r.Findings) > 0 {
				output.Section("findings")
				for _, f := range r.Findings {
					severity := f.Severity
					switch severity {
					case "critical", "major", "high":
						fmt.Printf("  %s%s%s  %s  %s\n", output.Red, severity, output.Reset, f.Category, f.Observation)
					case "medium", "moderate":
						fmt.Printf("  %s%s%s  %s  %s\n", output.Yellow, severity, output.Reset, f.Category, f.Observation)
					default:
						fmt.Printf("  %s%s%s  %s  %s\n", output.Dim, severity, output.Reset, f.Category, f.Observation)
					}
				}
			}
		}

		fmt.Println()
		fmt.Printf("Full report: %s\n", os.Getenv("UX_DASHBOARD_URL")+"/runs/"+runID)
		return nil
	},
}

func init() {
	resultsCmd.Flags().String("run", "", "Run ID (required)")
	resultsCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(resultsCmd)
}
