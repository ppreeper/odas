/*
Copyright © 2024 Peter Preeper

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package internal

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

//go:embed commit.txt
var Commit string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "odas",
	Short: "Odoo Server Administration Tool",
	Long:  `Odoo Server Administration Tool`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Commit
	rootCmd.AddGroup(
		&cobra.Group{ID: "app", Title: "App Management"},
		&cobra.Group{ID: "backup", Title: "Backup Management"},
		&cobra.Group{ID: "config", Title: "Config Commands (requries sudo)"},
		&cobra.Group{ID: "database", Title: "Database Management"},
		&cobra.Group{ID: "instance", Title: "Instance Management"},
		&cobra.Group{ID: "repo", Title: "Repository Management"},
		&cobra.Group{ID: "user", Title: "Admin User Management"},
	)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
