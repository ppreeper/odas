package internal

import (
	"html/template"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func aptInstall(packages ...string) {
	args := []string{"install", "-y", "--no-install-recommends"}
	args = append(args, packages...)
	cmdRun("apt-get", args...)
}

func npmInstall(packages ...string) {
	args := []string{"install", "-g"}
	args = append(args, packages...)
	cmdRun("npm", args...)
}

func cmdRun(cmd string, args ...string) {
	roleCmd := exec.Command(cmd, args...)
	roleCmd.Stdin = os.Stdin
	roleCmd.Stdout = os.Stdout
	roleCmd.Stderr = os.Stderr
	cobra.CheckErr(roleCmd.Run())
}

func roleCaddy() {
	url := "https://caddyserver.com/api/download?os=linux&arch=amd64&p=github.com%2Fcaddy-dns%2Fcloudflare"

	cmdRun("wget", "-qO", "/usr/local/bin/caddy", url)
	cmdRun("chmod", "+x", "/usr/local/bin/caddy")
	cmdRun("mkdir", "-p", "/etc/caddy")
}

func roleCaddyService() {
	fo, err := os.Create("/etc/systemd/system/caddy.service")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/caddy.service")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	cmdRun("sudo", "systemctl", "daemon-reload")
	cmdRun("sudo", "systemctl", "enable", "caddy.service")
	cmdRun("sudo", "mkdir", "-p", "/etc/caddy")
}

func roleGeoIP2DB() error {
	// install geolite databases
	geolite := [][]string{
		{"GeoLite2-ASN.mmdb", "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-ASN.mmdb"},
		{"GeoLite2-City.mmdb", "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb"},
		{"GeoLite2-Country.mmdb", "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-Country.mmdb"},
	}

	for _, geo := range geolite {
		cmdRun("wget", "-qO", "/usr/share/GeoIP/"+geo[0], geo[1])
	}

	return nil
}

func roleOdooService() {
	fo, err := os.Create("/etc/systemd/system/odoo.service")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/odoo.service")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	cmdRun("sudo", "systemctl", "daemon-reload")
	cmdRun("sudo", "systemctl", "enable", "odoo.service")
}

func roleOdooUser() {
	cmdRun("sudo", "groupadd", "-f", "-g", "1001", "odoo")
	cmdRun("sudo", "useradd", "-ms", "/bin/bash", "-g", "1001", "-u", "1001", "odoo")
	cmdRun("sudo", "usermod", "-aG", "sudo", "odoo")

	fo, err := os.Create("/etc/sudoers.d/odoo")
	cobra.CheckErr(err)
	fo.Close()

	data := map[string]string{}
	t, err := template.ParseGlob("odoo ALL=(ALL) NOPASSWD:ALL")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)

	cmdRun("sudo", "chown", "root:root", "/etc/sudoers.d/odoo")

	// SSH key
	cmdRun("sudo", "mkdir", "/home/odoo/.ssh")
	cmdRun("sudo", "chown", "-R", "odoo:odoo", "/home/odoo/.ssh")
}

func roleOdooDirs() {
	dirs := []string{"addons", "backups", "conf", "data", "odoo", "enterprise", "design-themes", "industry"}

	for _, dir := range dirs {
		cmdRun("mkdir", "-p", "/opt/odoo/"+dir)
	}
	cmdRun("sudo", "chown", "-R", "odoo:odoo", "/opt/odoo")
}

func rolePaperSize() {
	cmdRun("/usr/sbin/paperconfig", "-p", "letter")
}

func rolePostgresqlRepo() {
	roleUpdate()
	aptInstall("postgresql-common", "apt-transport-https", "ca-certificates")
	cmdRun("/usr/share/postgresql-common/pgdg/apt.postgresql.org.sh", "-y")
	roleUpdate()
}

func rolePostgresqlClient(dbVersion string) {
	aptInstall("postgresql-client-" + dbVersion)
}

func rolePreeperRepo() {
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
	cmdRun("/usr/local/bin/update")
}

func roleUpdateScript() {
	fo, err := os.Create("/usr/local/bin/update")
	cobra.CheckErr(err)
	defer fo.Close()

	data := map[string]string{}
	t, err := template.ParseFS(embedFS, "templates/update.sh")
	cobra.CheckErr(err)
	err = t.Execute(fo, data)
	cobra.CheckErr(err)
	err = os.Chmod("/usr/local/bin/update", 0o755)
	cobra.CheckErr(err)
}

func roleWkhtmltopdf() {
	// wkhtmltopdf
	url := "https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.jammy_amd64.deb"

	cmdRun("wget", "-qO", "wkhtmltox.deb", url)
	aptInstall("./wkhtmltox.deb")
	cmdRun("rm", "-rf", "wkhtmltox.deb")
}
