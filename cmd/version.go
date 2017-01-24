package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version of prometheus_sql_exporter",
	RunE: func(_ *cobra.Command, args []string) error {
		if version == "" {
			version = "unversioned"
		}
		fmt.Println(version)
		return nil
	},
}
