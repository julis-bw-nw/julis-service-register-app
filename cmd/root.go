package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "julis-service-register-app",
		Short: "Webapp to allow independent user registration via an admin generated token",
	}
)

func init() {
	rootCmd.AddCommand(licenseCmd, runCmd)
}

// Execute executes the root command.
func Execute(license string) error {
	softwareLicense = license
	return rootCmd.Execute()
}
