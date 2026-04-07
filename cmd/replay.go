package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/whitingdeltix/deltix-cli/internal/api"
	"github.com/whitingdeltix/deltix-cli/internal/config"
	"github.com/whitingdeltix/deltix-cli/internal/output"
	"github.com/spf13/cobra"
)

var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Run playbook regression (no LLM)",
	Long:  "Replays saved playbooks deterministically without LLM. Fast and free.",
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, _ := cmd.Flags().GetString("app")
		playbookID, _ := cmd.Flags().GetString("playbook")
		wait, _ := cmd.Flags().GetBool("wait")

		if appID == "" && playbookID == "" {
			return fmt.Errorf("--app or --playbook is required")
		}

		cfg, _ := config.Load()
		if apiURL != "" {
			cfg.APIURL = apiURL
		}
		if !cfg.IsLoggedIn() {
			return fmt.Errorf("not logged in — run: uxvalidator login")
		}

		client := api.NewClient(cfg)

		// Get playbooks to run
		var playbooks []api.Spec
		if playbookID != "" {
			// Single playbook mode
			playbooks = []api.Spec{{ID: playbookID}}
		} else {
			// All playbooks for the app
			specs, err := client.ListSpecs(appID)
			if err != nil {
				return fmt.Errorf("failed to list playbooks: %w", err)
			}
			if len(specs) == 0 {
				output.Info("No playbooks found. Run a validation and save the spec first.")
				return nil
			}
			playbooks = specs
		}

		output.Header("UX Validator — Playbook Regression (no LLM)")
		if appID != "" {
			app, err := client.GetApp(appID)
			if err == nil {
				fmt.Printf("App:       %s (%s)\n", app.Name, app.BundleID)
			}
		}
		fmt.Printf("Playbooks: %d\n", len(playbooks))
		fmt.Println()

		passed := 0
		failed := 0

		for _, pb := range playbooks {
			name := pb.ID[:8]
			if pb.Name != nil {
				name = *pb.Name
			}
			steps := 0
			if pb.StepCount != nil {
				steps = *pb.StepCount
			}

			// Trigger playback
			resp, err := client.TriggerPlayback(pb.ID)
			if err != nil {
				output.Cross(fmt.Sprintf("%-30s  %sfailed to start%s", name, output.Red, output.Reset))
				failed++
				continue
			}

			if !wait {
				fmt.Printf("  ▶ %-30s  started (run: %s)\n", name, resp.PlaybackRunID[:8])
				continue
			}

			// Poll for completion
			startTime := time.Now()
			var result *api.PlaybackResult
			for {
				result, err = client.GetPlaybackResult(resp.PlaybackRunID)
				if err != nil {
					break
				}
				if result.Status == "passed" || result.Status == "failed" || result.Status == "error" || result.Status == "completed" {
					break
				}
				time.Sleep(2 * time.Second)
			}

			elapsed := time.Since(startTime)

			if err != nil {
				output.Cross(fmt.Sprintf("%-30s  %serror%s", name, output.Red, output.Reset))
				failed++
				continue
			}

			isPassed := result.Status == "passed" || result.Status == "completed"
			if isPassed {
				output.Check(fmt.Sprintf("%-30s  %d steps  %.1fs  %s", name, steps, elapsed.Seconds(), output.PassFail(true)))
				passed++
			} else {
				output.Cross(fmt.Sprintf("%-30s  %d steps  %.1fs  %s", name, steps, elapsed.Seconds(), output.PassFail(false)))
				if result.FailureReason != nil {
					output.Info(fmt.Sprintf("  %s", *result.FailureReason))
				}
				failed++
			}
		}

		fmt.Println()
		output.Divider()
		fmt.Printf("Result: %d/%d passed\n", passed, passed+failed)

		if failed > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	replayCmd.Flags().String("app", "", "App ID — replay all playbooks")
	replayCmd.Flags().String("playbook", "", "Specific playbook ID to replay")
	replayCmd.Flags().Bool("wait", true, "Wait for completion")
	rootCmd.AddCommand(replayCmd)
}
