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

var playbooksCmd = &cobra.Command{
	Use:   "playbooks",
	Short: "List saved playbooks for an app",
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, _ := cmd.Flags().GetString("app")
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
		specs, err := client.ListSpecs(appID)
		if err != nil {
			return err
		}

		if len(specs) == 0 {
			output.Info("No playbooks saved. Run a validation and save the spec.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, output.Dim+"ID\tNAME\tSTEPS\tDIFFICULTY\tPLATFORM"+output.Reset)
		for _, s := range specs {
			name := "—"
			if s.Name != nil {
				name = *s.Name
			}
			steps := "—"
			if s.StepCount != nil {
				steps = fmt.Sprintf("%d", *s.StepCount)
			}
			diff := "—"
			if s.Difficulty != nil {
				diff = *s.Difficulty
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				s.ID[:8], name, steps, diff, s.Platform)
		}
		w.Flush()
		return nil
	},
}

func init() {
	playbooksCmd.Flags().String("app", "", "App ID (required)")
	rootCmd.AddCommand(playbooksCmd)
}
