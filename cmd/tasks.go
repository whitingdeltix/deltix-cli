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

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks for an app",
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
		tasks, err := client.ListTasks(appID)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			output.Info("No tasks defined. Add tasks at the dashboard.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, output.Dim+"ID\tDESCRIPTION\tCATEGORY\tSTEPS"+output.Reset)
		for _, t := range tasks {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
				t.ID[:8], t.Description, t.Category, t.MaxSteps)
		}
		w.Flush()
		return nil
	},
}

func init() {
	tasksCmd.Flags().String("app", "", "App ID (required)")
	rootCmd.AddCommand(tasksCmd)
}
