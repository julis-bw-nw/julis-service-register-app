package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	softwareLicense string

	licenseCmd = &cobra.Command{
		Use:   "license",
		Short: "Shows the license of the program",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(softwareLicense)
		},
	}
)
