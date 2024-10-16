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
	"path/filepath"

	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var repoUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "odoo repository update",
	Long:  `odoo repository update`,
	Run: func(cmd *cobra.Command, args []string) {
		repoDir := filepath.Join("/", "opt", "odoo")
		for _, repo := range OdooRepos {
			fmt.Fprintln(os.Stderr, "Updating", repo)
			dest := filepath.Join(repoDir, repo)

			pull := exec.Command("git", "pull", "--rebase")
			pull.Dir = dest
			if err := pull.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "git pull on %s %v", repo, err)
				return
			}
		}
	},
}

func init() {
	repoCmd.AddCommand(repoUpdateCmd)
}
