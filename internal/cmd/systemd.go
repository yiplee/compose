/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"html/template"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

// systemdCmd represents the systemd command
var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "Generate systemd uint file by current working dir",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = path.Base(dir)
		}

		bin, err := exec.LookPath(cmd.Root().Name())
		if err != nil {
			return err
		}

		return generateSystemdUnitFile(name, dir, bin, cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(systemdCmd)

	systemdCmd.Flags().StringP("name", "n", "", "name of the service")
}

func generateSystemdUnitFile(name, dir, bin string, w io.Writer) error {
	t, err := template.New(name).Parse(systemdUnitTemplate)
	if err != nil {
		return err
	}

	return t.Execute(w, map[string]any{
		"Name": name,
		"Dir":  dir,
		"Bin":  bin,
	})
}

const systemdUnitTemplate = `
[Unit]
Description={{.Name}} compose service

[Service]
WorkingDirectory={{.Dir}}
ExecStart={{.Bin}}

[Install]
WantedBy=multi-user.target
`
