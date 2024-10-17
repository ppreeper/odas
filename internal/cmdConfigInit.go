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

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "system initialization",
	Long:  `system initialization`,
	Run: func(cmd *cobra.Command, args []string) {
		confirm := AreYouSure("initialize the system and delete all data?")

		if !confirm {
			fmt.Println("System initialization aborted.")
			return
		}

		var odooVersions []string
		switch GetOSVersion() {
		case "24.04":
			odooVersions = []string{"18.0"}
		case "22.04":
			odooVersions = []string{"17.0", "16.0", "15.0"}
		}

		versionOptions := []huh.Option[string]{}
		for _, version := range odooVersions {
			versionOptions = append(versionOptions, huh.NewOption(version, version))
		}
		var version string
		var create bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Odoo Version").
					Options(versionOptions...).
					Value(&version),

				huh.NewConfirm().
					Title("Start Build?").
					Value(&create),
			),
		)
		if err := form.Run(); err != nil {
			return
		}
		if !create {
			fmt.Println("System initialization aborted.")
			return
		}
		BaseCreate(version, localDomain)
	},
}
var localDomain string

func init() {
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().StringVar(&localDomain, "domain", "local", "domain name")
}

func GetOdooConfig(version string) OdooConfig {
	for _, config := range OdooConfigs {
		if config.Version == version {
			return config
		}
	}
	return OdooConfig{}
}

func BaseCreate(version, localDomain string) {
	fmt.Println("Creating base image for Odoo version", version, "on", localDomain)
	config := GetOdooConfig(version)

	roleUpdateScript()

	rolePreeperRepo()

	roleUpdate()

	aptInstall(config.BaselinePackages...)

	roleOdooUser()

	roleOdooDirs()

	rolePostgresqlRepo()

	rolePostgresqlClient(OdooDatabase.Version)

	roleWkhtmltopdf()

	aptInstall(config.Odoobase...)

	npmInstall("rtlcss")

	roleGeoIP2DB()

	rolePaperSize()

	roleOdooService()

	roleCaddy()

	roleCaddyService()
}
