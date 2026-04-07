package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/whitingdeltix/deltix-cli/internal/api"
	"github.com/whitingdeltix/deltix-cli/internal/config"
	"github.com/whitingdeltix/deltix-cli/internal/output"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Run full UX validation (uses LLM)",
	Long:  "Triggers a full UX validation run using Claude to navigate the app and evaluate UX quality.",
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, _ := cmd.Flags().GetString("app")
		device, _ := cmd.Flags().GetString("device")
		threshold, _ := cmd.Flags().GetInt("threshold")
		taskIDs, _ := cmd.Flags().GetString("tasks")
		branch, _ := cmd.Flags().GetString("branch")
		commit, _ := cmd.Flags().GetString("commit")
		wait, _ := cmd.Flags().GetBool("wait")

		if appID == "" {
			return fmt.Errorf("--app is required")
		}

		cfg, _ := config.Load()
		if apiURL != "" {
			cfg.APIURL = apiURL
		}
		if !cfg.IsLoggedIn() {
			return fmt.Errorf("not logged in — run: uxvalidator login")
		}

		client := api.NewClient(cfg)

		// Get app info
		app, err := client.GetApp(appID)
		if err != nil {
			return fmt.Errorf("failed to get app: %w", err)
		}

		output.Header("UX Validator — Full Validation (LLM)")
		fmt.Printf("App:    %s (%s)\n", app.Name, app.BundleID)
		fmt.Printf("Device: %s\n", device)

		// Build run request
		req := api.TriggerRunRequest{
			DeviceType: device,
			CommitHash: commit,
		}
		if taskIDs != "" {
			req.TaskIDs = strings.Split(taskIDs, ",")
			fmt.Printf("Tasks:  %d selected\n", len(req.TaskIDs))
		} else {
			fmt.Printf("Tasks:  all (%d)\n", app.TaskCount)
		}
		if branch != "" {
			fmt.Printf("Branch: %s\n", branch)
		}

		// Trigger run
		run, err := client.TriggerRun(appID, req)
		if err != nil {
			return fmt.Errorf("failed to trigger run: %w", err)
		}

		fmt.Printf("\nRun started: %s\n", run.ID[:8])

		if !wait {
			fmt.Printf("Track progress: uxvalidator status --run %s\n", run.ID)
			return nil
		}

		// Poll for completion
		fmt.Println()
		for {
			run, err = client.GetRun(run.ID)
			if err != nil {
				return fmt.Errorf("failed to get run status: %w", err)
			}

			fmt.Printf("\r  %s  %d/%d tasks completed",
				output.Status(run.Status), run.CompletedCount, run.TaskCount)

			if run.Status == "completed" || run.Status == "failed" || run.Status == "cancelled" {
				fmt.Println()
				break
			}

			time.Sleep(3 * time.Second)
		}

		// Show results
		results, err := client.GetRunResults(run.ID)
		if err != nil {
			return fmt.Errorf("failed to get results: %w", err)
		}

		fmt.Println()
		for _, r := range results {
			desc := "Task " + r.TaskID[:8]
			if r.TaskDescription != nil {
				desc = *r.TaskDescription
			}
			if r.Succeeded {
				output.Check(fmt.Sprintf("%s  %s%d%s/100  %d steps",
					desc, output.ScoreColor(r.AggregateScore), r.AggregateScore, output.Reset, r.Steps))
			} else {
				output.Cross(fmt.Sprintf("%s  %sfailed%s  %d steps",
					desc, output.Red, output.Reset, r.Steps))
			}
		}

		output.Divider()

		score := 0
		if run.AggregateScore != nil {
			score = *run.AggregateScore
		}

		passed := 0
		failed := 0
		for _, r := range results {
			if r.Succeeded {
				passed++
			} else {
				failed++
			}
		}

		thresholdMet := threshold == 0 || score >= threshold
		thresholdStr := ""
		if threshold > 0 {
			if thresholdMet {
				thresholdStr = fmt.Sprintf(" (threshold: %d %s✓%s)", threshold, output.Green, output.Reset)
			} else {
				thresholdStr = fmt.Sprintf(" (threshold: %d %s✗%s)", threshold, output.Red, output.Reset)
			}
		}

		fmt.Printf("Result: %s%d%s/100%s\n", output.ScoreColor(score), score, output.Reset, thresholdStr)
		fmt.Printf("Passed: %d/%d  Failed: %d/%d\n", passed, passed+failed, failed, passed+failed)

		// Show top findings
		var allFindings []api.Finding
		for _, r := range results {
			allFindings = append(allFindings, r.Findings...)
		}
		if len(allFindings) > 0 {
			output.Section("findings")
			limit := 5
			if len(allFindings) < limit {
				limit = len(allFindings)
			}
			for _, f := range allFindings[:limit] {
				switch f.Severity {
				case "critical", "major", "high":
					output.Cross(fmt.Sprintf("%s  %s", f.Severity, f.Observation))
				case "medium", "moderate":
					output.Warn(fmt.Sprintf("%s  %s", f.Severity, f.Observation))
				default:
					output.Info(fmt.Sprintf("%s  %s", f.Severity, f.Observation))
				}
			}
			if len(allFindings) > limit {
				output.Info(fmt.Sprintf("+%d more findings", len(allFindings)-limit))
			}
		}

		fmt.Printf("\nFull report: %s/runs/%s\n", cfg.APIURL, run.ID)

		if !thresholdMet {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	validateCmd.Flags().String("app", "", "App ID (required)")
	validateCmd.Flags().String("device", "simulator", "Device type: simulator or real_device")
	validateCmd.Flags().String("tasks", "", "Comma-separated task IDs (default: all)")
	validateCmd.Flags().String("branch", "", "Branch name (stored on run)")
	validateCmd.Flags().String("commit", "", "Commit hash (stored on run)")
	validateCmd.Flags().Int("threshold", 0, "Minimum score — exit 1 if below")
	validateCmd.Flags().Bool("wait", true, "Wait for completion and show results")
	rootCmd.AddCommand(validateCmd)
}
