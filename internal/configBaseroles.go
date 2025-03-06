package internal

import (
	"embed"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

func aptInstall(packages ...string) {
	args := []string{"install", "-y", "--no-install-recommends"}
	args = append(args, packages...)
	CmdRun("apt-get", args...)
}

func npmInstall(packages ...string) {
	args := []string{"install", "-g"}
	args = append(args, packages...)
	CmdRun("npm", args...)
}

func roleCaddy() {
	url := "https://caddyserver.com/api/download?os=linux&arch=amd64&p=github.com%2Fcaddy-dns%2Fcloudflare"

	CmdRun("wget", "-qO", "/usr/local/bin/caddy", url)
	CmdRun("chmod", "+x", "/usr/local/bin/caddy")
	CmdRun("mkdir", "-p", "/etc/caddy")
}

func roleCaddyService(embedFS embed.FS) {
	fo, err := os.Create("/etc/systemd/system/caddy.service")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/caddy.service")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	CmdRun("sudo", "systemctl", "daemon-reload")
	CmdRun("sudo", "systemctl", "enable", "caddy.service")
	CmdRun("sudo", "mkdir", "-p", "/etc/caddy")
}

func roleGeoIP2DB() error {
	// install geolite databases
	geolite := [][]string{
		{"GeoLite2-ASN.mmdb", "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb"},
		{"GeoLite2-City.mmdb", "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb"},
		{"GeoLite2-Country.mmdb", "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-Country.mmdb"},
	}

	for _, geo := range geolite {
		CmdRun("wget", "-qO", "/usr/share/GeoIP/"+geo[0], geo[1])
	}

	return nil
}

func roleOdooService(embedFS embed.FS) {
	fo, err := os.Create("/etc/systemd/system/odoo.service")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/odoo.service")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	CmdRun("sudo", "systemctl", "daemon-reload")
	CmdRun("sudo", "systemctl", "enable", "odoo.service")
}

func roleOdooUser() {
	CmdRun("sudo", "groupadd", "-f", "-g", "1001", "odoo")
	CmdRun("sudo", "useradd", "-ms", "/bin/bash", "-g", "1001", "-u", "1001", "odoo")
	CmdRun("sudo", "usermod", "-aG", "sudo", "odoo")

	fo, err := os.Create("/etc/sudoers.d/odoo")
	cobra.CheckErr(err)
	fo.Close()

	data := map[string]string{}
	t, err := template.ParseGlob("odoo ALL=(ALL) NOPASSWD:ALL")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	CmdRun("sudo", "chown", "root:root", "/etc/sudoers.d/odoo")

	// SSH key
	CmdRun("sudo", "mkdir", "/home/odoo/.ssh")
	CmdRun("sudo", "chown", "-R", "odoo:odoo", "/home/odoo/.ssh")
}

func roleOdooDirs() {
	dirs := []string{"addons", "backups", "conf", "data", "odoo", "enterprise", "design-themes", "industry"}

	for _, dir := range dirs {
		CmdRun("mkdir", "-p", "/opt/odoo/"+dir)
	}
	CmdRun("sudo", "chown", "-R", "odoo:odoo", "/opt/odoo")
}

func rolePaperSize() {
	CmdRun("/usr/sbin/paperconfig", "-p", "letter")
}

func rolePGCat() {
	url := "https://www.preeper.org/pgcat"
	CmdRun("wget", "-qO", "/usr/local/bin/pgcat", url)
	CmdRun("chmod", "+x", "/usr/local/bin/pgcat")
}

func rolePGCatService(embedFS embed.FS) {
	fo, err := os.Create("/etc/systemd/system/pgcat.service")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/pgcat.service")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	CmdRun("sudo", "systemctl", "daemon-reload")
	CmdRun("sudo", "systemctl", "enable", "pgcat.service")
}

func rolePostgresqlRepo() {
	roleUpdate()
	aptInstall("postgresql-common", "apt-transport-https", "ca-certificates")
	CmdRun("/usr/share/postgresql-common/pgdg/apt.postgresql.org.sh", "-y")
	roleUpdate()
}

func rolePostgresqlClient(dbVersion string) {
	aptInstall("postgresql-client-" + dbVersion)
}

func rolePreeperRepo(embedFS embed.FS) {
	fo, err := os.Create("/etc/apt/sources.list.d/preeper.list")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/preeper.list")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	roleUpdate()
}

func roleUpdate() {
	CmdRun("/usr/local/bin/update")
}

func roleUpdateScript(embedFS embed.FS) {
	fo, err := os.Create("/usr/local/bin/update")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/update")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)
	err = os.Chmod("/usr/local/bin/update", 0o755)
	cobra.CheckErr(err)
}

func roleWkhtmltopdf() {
	// wkhtmltopdf
	url := "https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.jammy_amd64.deb"

	CmdRun("wget", "-qO", "wkhtmltox.deb", url)
	aptInstall("./wkhtmltox.deb")
	CmdRun("rm", "-rf", "wkhtmltox.deb")
}
