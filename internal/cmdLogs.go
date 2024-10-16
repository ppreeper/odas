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
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Follow the logs",
	Long:  `Follow the logs`,
	Run: func(cmd *cobra.Command, args []string) {
		command := exec.Command("sudo",
			"journalctl",
			"-u",
			"odoo.service",
			"-f",
		)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "error getting logs %w", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
