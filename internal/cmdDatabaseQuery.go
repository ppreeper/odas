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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/odoorpc/odoojrpc"
	"github.com/spf13/cobra"
)

var q QueryDef

// hostsCmd represents the hosts command
var queryCmd = &cobra.Command{
	Use:     "query",
	Short:   "Query an Odoo model",
	Long:    `Query an Odoo model`,
	GroupID: "database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "no model specified")
			return
		}
		q.Model = args[0]

		dbName := GetOdooConf("db_name")

		oc := odoojrpc.NewOdoo().
			WithHostname("127.0.0.1").
			WithPort(8069).
			WithDatabase(dbName).
			WithSchema("http").
			WithUsername(q.Username).
			WithPassword(q.Password)

		err := oc.Login()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating odoo rpc %w", err)
			return
		}

		umdl := strings.Replace(q.Model, "_", ".", -1)

		fields := parseFields(q.Fields)
		if q.Count {
			fields = []string{"id"}
		}

		filtp, err := parseFilter(q.Filter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}

		rr, err := oc.SearchRead(umdl, q.Offset, q.Limit, fields, filtp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}
		if q.Count {
			fmt.Fprintln(os.Stderr, "records:", len(rr))
		} else {
			jsonStr, err := json.MarshalIndent(rr, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			}
			fmt.Fprintln(os.Stderr, string(jsonStr))
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&q.Filter, "filter", "d", "", "domain filter")
	queryCmd.Flags().IntVarP(&q.Offset, "offset", "o", 0, "offset")
	queryCmd.Flags().IntVarP(&q.Limit, "limit", "l", 0, "limit records returned")
	queryCmd.Flags().StringVarP(&q.Fields, "fields", "f", "", "fields to return")
	queryCmd.Flags().BoolVarP(&q.Count, "count", "c", false, "count records")
	queryCmd.Flags().StringVarP(&q.Username, "username", "u", "admin", "username")
	queryCmd.Flags().StringVarP(&q.Password, "password", "p", "admin", "password")
}
