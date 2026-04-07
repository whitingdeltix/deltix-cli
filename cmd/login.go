package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/whitingdeltix/deltix-cli/internal/api"
	"github.com/whitingdeltix/deltix-cli/internal/config"
	"github.com/whitingdeltix/deltix-cli/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with UX Validator",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()
		if apiURL != "" {
			cfg.APIURL = apiURL
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password := string(passwordBytes)

		client := api.NewClient(cfg)
		resp, err := client.Login(email, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		cfg.Token = resp.AccessToken
		cfg.UserID = resp.UserID
		cfg.Username = resp.Username

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		output.Check(fmt.Sprintf("Logged in as %s%s%s", output.Bold, resp.Username, output.Reset))
		output.Info(fmt.Sprintf("Token saved to ~/.uxvalidator/config.json"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
