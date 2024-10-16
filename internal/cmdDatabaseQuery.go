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

	"github.com/spf13/cobra"
)

// hostsCmd represents the hosts command
var queryCmd = &cobra.Command{
	Use:     "query",
	Short:   "Query an Odoo model",
	Long:    `Query an Odoo model`,
	GroupID: "database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("query")
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
