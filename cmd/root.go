package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var apiURL string

var rootCmd = &cobra.Command{
	Use:   "uxvalidator",
	Short: "UX Validator CLI — automated UX quality validation",
	Long: `UX Validator CLI

Validate your app's UX quality and run playbook regressions
from the command line or CI pipeline.

  uxvalidator validate --app <id>    Full validation (LLM)
  uxvalidator replay --app <id>      Playbook regression (no LLM)
  uxvalidator apps                   List your apps
  uxvalidator playbooks --app <id>   List saved playbooks`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API base URL (overrides config)")
}
