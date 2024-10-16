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

var (
	restoreAny        bool
	restoreMove       bool
	restoreNeutralize bool
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Restore database and filestore or addons",
	Long:    `Restore database and filestore or addons`,
	GroupID: "backup",
	Args: func(cmd *cobra.Command, args []string) error {
		if restoreMove && restoreNeutralize {
			return fmt.Errorf("cannot move and neutralize at the same time")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		AdminRestore(
			restoreAny,
			restoreMove,
			restoreNeutralize,
		)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolVarP(&restoreAny, "any", "a", false, "any backup")
	restoreCmd.Flags().BoolVarP(&restoreMove, "move", "m", false, "move server")
	restoreCmd.Flags().BoolVarP(&restoreNeutralize, "neutralize", "n", false, "fully neutralize")
}
