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
package main

import (
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/ppreeper/odas/internal"

	"github.com/urfave/cli/v2"
)

//go:generate sh -c "printf '%s (%s)' $(git tag -l | sort -V | tail -1) $(date +%Y%m%d)-$(git rev-parse --short HEAD)"
//go:embed commit.txt
var commit string

//go:embed templates/*
var templates embed.FS

func main() {
	oda := internal.NewODA("odas", "Odoo Server Administration Tool", commit, templates).GetOdooConf()

	app := &cli.App{
		Name:                 oda.Name,
		Usage:                oda.Usage,
		Version:              oda.Version,
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			// --------------------------------------------------------------------------
			// Admin User Management
			{
				Name:     "admin",
				Usage:    "Admin user management",
				Category: "Admin User Management",
				Subcommands: []*cli.Command{
					{
						Name:  "username",
						Usage: "Odoo Admin username",
						Action: func(cCtx *cli.Context) error {
							// direct connection to the database
							// set username of user id=2
							return oda.AdminUsername()
						},
					},
					{
						Name:  "password",
						Usage: "Odoo Admin password",
						Action: func(cCtx *cli.Context) error {
							// direct connection to the database
							// set password of user id=2
							return oda.AdminPassword()
						},
					},
					{
						Name:  "updateuser",
						Usage: "Odoo Update User",
						Action: func(cCtx *cli.Context) error {
							// direct connection to the database
							// set username and password of user
							return oda.UpdateUser()
						},
					},
				},
			},
			// --------------------------------------------------------------------------
			// App Management
			{
				Name:     "install",
				Usage:    "Install module(s)",
				Category: "App Management",
				Action: func(cCtx *cli.Context) error {
					modlen := cCtx.Args().Len()
					if modlen == 0 {
						return fmt.Errorf("no modules specified")
					}
					// cmdRun := exec.Command("odoo-bin", "shell", "-c", "/opt/odoo/conf/odoo.conf" install module(s))
					return oda.InstanceAppInstallUpgrade(true, cCtx.Args().Slice()...)
				},
			},
			{
				Name:     "upgrade",
				Usage:    "Upgrade module(s)",
				Category: "App Management",
				Action: func(cCtx *cli.Context) error {
					modlen := cCtx.Args().Len()
					if modlen == 0 {
						return fmt.Errorf("no modules specified")
					}
					// cmdRun := exec.Command("odoo-bin", "shell", "-c", "/opt/odoo/conf/odoo.conf" update module(s))
					return oda.InstanceAppInstallUpgrade(false, cCtx.Args().Slice()...)
				},
			},
			{
				Name:     "scaffold",
				Usage:    "Generates an Odoo module skeleton in addons",
				Category: "App Management",
				Action: func(cCtx *cli.Context) error {
					modlen := cCtx.Args().Len()
					if modlen == 0 {
						return fmt.Errorf("no modules specified")
					}
					// cmdRun := exec.Command("odoo-bin", "shell", "-c", "/opt/odoo/conf/odoo.conf" scaffold module)
					return oda.Scaffold(cCtx.Args().First())
				},
			},
			// --------------------------------------------------------------------------
			// Backup Restore
			{
				Name:     "backup",
				Usage:    "Backup database filestore and addons",
				Category: "Backup Management",
				Action: func(cCtx *cli.Context) error {
					// direct connection to the database
					// backup database, filestore and addons
					return oda.Backup()
				},
			},
			{
				Name:     "restore",
				Usage:    "Restore database and filestore or addons",
				Category: "Backup Management",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "any",
						Value: false,
						Usage: "any backup",
					},
					&cli.BoolFlag{
						Name:  "move",
						Value: false,
						Usage: "move server",
					},
					&cli.BoolFlag{
						Name:  "neutralize",
						Value: false,
						Usage: "fully neutralize the server",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.Bool("move") && cCtx.Bool("neutralize") {
						return fmt.Errorf("cannot move and neutralize at the same time")
					}
					// direct connection to the database
					return oda.Restore(
						cCtx.Bool("any"),
						cCtx.Bool("move"),
						cCtx.Bool("neutralize"),
					)
				},
			},
			{
				Name:     "trim",
				Usage:    "Trim database backups",
				Category: "Backup Management",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "limit",
						Value: 10,
						Usage: "number of backups to keep",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// directly look at the files in backup directory
					return oda.Trim(cCtx.Int("limit"), false)
				},
			},
			{
				Name:     "trimall",
				Usage:    "Trim all database backups",
				Category: "Backup Management",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "limit",
						Value: 10,
						Usage: "number of backups to keep",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// directly look at the files in backup directory
					return oda.Trim(cCtx.Int("limit"), true)
				},
			},
			// --------------------------------------------------------------------------
			// Config Commands (requries sudo)
			{
				Name:     "caddy",
				Usage:    "update caddyfile",
				Category: "Config Commands (requries sudo)",
				Action: func(cCtx *cli.Context) error {
					domain := ""
					modlen := cCtx.Args().Len()
					if modlen == 1 {
						domain = cCtx.Args().First()
					}
					// create if not exists directory /etc/caddy/
					// create file /etc/caddy/Caddyfile
					// insert hostname.domain into Caddyfile
					// set tls, if tls_key is present and not blank then use cloudflare
					// format Caddyfile
					return oda.CaddyfileUpdate(domain)
				},
			},
			{
				Name:     "hosts",
				Usage:    "update hosts file",
				Category: "Config Commands (requries sudo)",
				Action: func(cCtx *cli.Context) error {
					domain := ""
					modlen := cCtx.Args().Len()
					if modlen == 0 {
						domain = "local"
					} else {
						domain = cCtx.Args().First()
					}
					// update if not exists directory /etc/hosts
					// insert hostname.domain into /etc/hosts
					return oda.HostsUpdate(domain)
				},
			},
			{
				Name:     "pgcat",
				Usage:    "pgcat setup",
				Category: "Config Commands (requries sudo)",
				Action: func(cCtx *cli.Context) error {
					// create if not exists directory /etc/pgcat.toml
					// uses db_user, db_password, db_host, db_port, db_name

					return oda.PGCatUpdate()
				},
			},
			{
				Name:     "config",
				Usage:    "config management",
				Category: "Config Commands (requries sudo)",
				Action: func(cCtx *cli.Context) error {
					domain := ""
					modlen := cCtx.Args().Len()
					if modlen == 0 {
						domain = "local"
					} else {
						domain = cCtx.Args().First()
					}
					// create if not exists directory /etc/oda/
					// initialize config /etc/oda/odas.yaml
					// set domain, db_user, db_password, db_host, db_port, db_name
					// set odoo_user, odoo_password, odoo_port
					// set odoo_branch, odoo_version
					// set odoo_addons_path, odoo_data_dir
					// maybe have a small tui to set these values
					return oda.ConfigInit(domain)
				},
			},
			// Database Management
			{
				Name:     "psql",
				Usage:    "Access the instance database",
				Category: "Database Management",
				Action: func(cCtx *cli.Context) error {
					// direct connection to the database
					// using cmdRun psql to connect to the database
					return oda.PSQL()
				},
			},
			{
				Name:      "query",
				Usage:     "Query an Odoo model",
				Category:  "Database Management",
				UsageText: "query <model> [command options]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "domain",
						Aliases:     []string{"d"},
						Value:       "",
						Usage:       "domain filter",
						Destination: &oda.Q.Filter,
					},
					&cli.IntFlag{
						Name:        "offset",
						Aliases:     []string{"o"},
						Value:       0,
						Usage:       "offset",
						Destination: &oda.Q.Offset,
					},
					&cli.IntFlag{
						Name:        "limit",
						Aliases:     []string{"l"},
						Value:       0,
						Usage:       "limit records returned",
						Destination: &oda.Q.Limit,
					},
					&cli.StringFlag{
						Name:        "fields",
						Aliases:     []string{"f"},
						Value:       "",
						Usage:       "fields to return",
						Destination: &oda.Q.Fields,
					},
					&cli.BoolFlag{
						Name:        "count",
						Aliases:     []string{"c"},
						Value:       false,
						Usage:       "count records",
						Destination: &oda.Q.Count,
					},
					&cli.StringFlag{
						Name:        "username",
						Aliases:     []string{"u"},
						Value:       "admin",
						Usage:       "username",
						Destination: &oda.Q.Username,
					},
					&cli.StringFlag{
						Name:        "password",
						Aliases:     []string{"p"},
						Value:       "admin",
						Usage:       "password",
						Destination: &oda.Q.Password,
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() == 0 {
						return fmt.Errorf("no model specified")
					}
					oda.Q.Model = cCtx.Args().First()
					// direct connection to the database
					// uses odoorpc to query the database
					return oda.Query()
				},
			},
			// Instance Management
			{
				Name:     "start",
				Usage:    "Start the instance",
				Category: "Instance Management",
				Action: func(cCtx *cli.Context) error {
					return oda.OdooStart()
				},
			},
			{
				Name:     "stop",
				Usage:    "Stop the instance",
				Category: "Instance Management",
				Action: func(cCtx *cli.Context) error {
					return oda.OdooStop()
				},
			},
			{
				Name:     "restart",
				Usage:    "Restart the instance",
				Category: "Instance Management",
				Action: func(cCtx *cli.Context) error {
					return oda.OdooRestart()
				},
			},
			{
				Name:     "logs",
				Usage:    "Follow the logs",
				Category: "Instance Management",
				Action: func(cCtx *cli.Context) error {
					return oda.Logs()
				},
			},
			// Repository Management
			{
				Name:     "repo",
				Usage:    "odoo repository management",
				Category: "Repository Management",
				Action: func(cCtx *cli.Context) error {
					// does a git pull on the
					// "odoo", "enterprise", "design-themes", "industry"
					// repositories
					// stored in /opt/odoo
					return oda.RepoUpdate()
				},
			},
			// Utility Commands
			{
				Name:     "welcome",
				Usage:    "Welcome message",
				Category: "Utility Commands",
				Action: func(cCtx *cli.Context) error {
					// looks at the system and returns a welcome message
					return oda.Welcome()
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
