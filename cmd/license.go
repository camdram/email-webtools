package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Print copyright license",
	Long:  "Prints email-webtools' copyright and licensing information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Copyright (c) 2019-2021 Members of the Camdram Web Team and other contributors.")
		fmt.Println("This software is released under the MIT License.")
		fmt.Println("Please visit https://github.com/camdram/email-webtools for more information.")
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
}
