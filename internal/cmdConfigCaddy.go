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
	"github.com/spf13/cobra"
)

// caddyCmd represents the caddy command
var caddyCmd = &cobra.Command{
	Use:     "caddy",
	Short:   "update caddyfile",
	Long:    `update caddyfile`,
	GroupID: "config",
	Run: func(cmd *cobra.Command, args []string) {
		domain := ""
		modlen := len(args)
		if modlen == 0 {
			domain = "local"
		} else {
			domain = args[0]
		}
		CaddyfileUpdate(domain)
	},
}

func init() {
	rootCmd.AddCommand(caddyCmd)
}
