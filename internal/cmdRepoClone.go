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

var branchVersion string

// adminCmd represents the admin command
var repoBranchCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "odoo repository clone",
	Long:  `odoo repository clone`,
	Run: func(cmd *cobra.Command, args []string) {
		// var version string
		// var create bool

		// repoShorts, _ := RepoShortCodes("odoo")
		// versions := GetCurrentOdooRepos()
		// for _, version := range versions {
		// 	repoShorts = removeValue(repoShorts, version)
		// }
		// versionOptions := []huh.Option[string]{}
		// for _, version := range repoShorts {
		// 	versionOptions = append(versionOptions, huh.NewOption(version, version))
		// }

		// if len(versionOptions) == 0 {
		// 	fmt.Fprintln(os.Stderr, "no more branches to clone")
		// 	return
		// }

		// form := huh.NewForm(
		// 	huh.NewGroup(
		// 		huh.NewSelect[string]().
		// 			Title("Available Odoo Branches").
		// 			Options(versionOptions...).
		// 			Value(&version),

		// 		huh.NewConfirm().
		// 			Title("Clone Branch?").
		// 			Value(&create),
		// 	),
		// )
		// if err := form.Run(); err != nil {
		// 	fmt.Fprintf(os.Stderr, "odoo version form error %v", err)
		// 	return
		// }

		// if !create {
		// 	return
		// }

		// repoDir := viper.GetString("dirs.repo")

		// for _, repo := range OdooRepos {
		// 	source := filepath.Join(repoDir, repo)
		// 	dest := filepath.Join(repoDir, version, repo)
		// 	if err := CopyDirectory(source, dest); err != nil {
		// 		fmt.Fprintf(os.Stderr, "copy directory %s to %s failed %v", source, dest, err)
		// 		return
		// 	}

		// 	fetcher := exec.Command("git", "fetch", "origin")
		// 	fetcher.Dir = dest
		// 	if err := fetcher.Run(); err != nil {
		// 		fmt.Fprintf(os.Stderr, "git fetch origin on %s %v", repo, err)
		// 		return
		// 	}

		// 	checkout := exec.Command("git", "checkout", version)
		// 	checkout.Dir = dest
		// 	if err := checkout.Run(); err != nil {
		// 		fmt.Fprintf(os.Stderr, "git checkout on %s %v", repo, err)
		// 		return
		// 	}

		// 	pull := exec.Command("git", "pull", "origin", version)
		// 	pull.Dir = dest
		// 	if err := pull.Run(); err != nil {
		// 		fmt.Fprintf(os.Stderr, "git pull on %s %v", repo, err)
		// 		return
		// 	}
		// }
	},
}

func init() {
	repoCmd.AddCommand(repoBranchCloneCmd)
	repoBranchCloneCmd.Flags().StringVarP(&branchVersion, "version", "v", "", "Odoo version")
}
