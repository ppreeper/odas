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

// hostsCmd represents the hosts command
var psqlCmd = &cobra.Command{
	Use:     "psql",
	Short:   "Access the instance database",
	Long:    `Access the instance database`,
	GroupID: "database",
	Run: func(cmd *cobra.Command, args []string) {
		dbHost := GetOdooConf("db_host")
		dbPort := GetOdooConf("db_port")
		dbName := GetOdooConf("db_name")
		dbUser := GetOdooConf("db_user")
		dbPassword := GetOdooConf("db_password")

		pgCmd := exec.Command("psql", "-h", dbHost, "-p", dbPort, "-d", dbName, "-U", dbUser)
		pgCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))
		pgCmd.Stdin = os.Stdin
		pgCmd.Stdout = os.Stdout
		pgCmd.Stderr = os.Stderr
		if err := pgCmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to run psql %w", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(psqlCmd)
}
