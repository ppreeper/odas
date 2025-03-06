package internal

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/charmbracelet/huh"
)

func (o *ODA) CaddyfileUpdate(domain string) error {
	fqdn, hostname, _ := GetFQDN()
	if domain != "" {
		fqdn = hostname + "." + domain
	}

	tmpl, err := template.ParseFS(o.EmbedFS, "templates/Caddyfile")
	if err != nil {
		return err
	}

	data := struct {
		FQDN     string
		Hostname string
		Domain   string
	}{
		FQDN:     fqdn,
		Hostname: hostname,
		Domain:   domain,
	}

	f, err := os.Create("/etc/caddy/Caddyfile")
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo",
		"caddy", "fmt", "--overwrite", "/etc/caddy/Caddyfile",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("caddyfile format failed: %w", err)
	}

	return nil
}

func (o *ODA) HostsUpdate(domain string) error {
	fqdn, hostname, _ := GetFQDN()
	if domain != "" {
		fqdn = hostname + "." + domain
	}

	tmpl, err := template.ParseFS(o.EmbedFS, "templates/hosts")
	if err != nil {
		return err
	}

	data := struct {
		FQDN     string
		Hostname string
		Domain   string
	}{
		FQDN:     fqdn,
		Hostname: hostname,
		Domain:   domain,
	}

	f, err := os.Create("/etc/hosts")
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func (o *ODA) PGCatUpdate() error {
	rolePGCat()
	rolePGCatService(o.EmbedFS)

	tmpl, err := template.ParseFS(o.EmbedFS, "templates/pgcat.toml")
	if err != nil {
		return err
	}

	data := struct {
		DBHost     string
		DBPort     string
		DBName     string
		DBUsername string
		DBPassword string
	}{
		DBHost:     o.OdooDatabase.Name,
		DBPort:     "5432",
		DBName:     o.OdooConf.DbName,
		DBUsername: o.OdooConf.DbUser,
		DBPassword: o.OdooConf.DbPassword,
	}

	f, err := os.Create("/etc/pgcat.toml")
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func (o *ODA) ConfigInit(localDomain string) error {
	confirm := AreYouSure("initialize the system and delete all data?")

	if !confirm {
		return fmt.Errorf("system initialization aborted")
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
	dbVersions := []huh.Option[string]{}
	for _, version := range []string{"17", "16", "15"} {
		dbVersions = append(dbVersions, huh.NewOption(version, version))
	}

	var version string
	var create bool
	var hostname string
	production := false
	container := true
	var databaseVersion string
	cluster := false
	var hostPort string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Odoo Version").
				Options(versionOptions...).
				Value(&version),
			// Ask for hostname (default current hostname)
			huh.NewInput().
				Title("Please enter hostname:").
				Prompt(">").
				Value(&hostname),
			// Ask for domain (default local)
			huh.NewInput().
				Title("Please enter domainname:").
				Prompt(">").
				Value(&localDomain),
			// Is Production (default false)
			huh.NewConfirm().
				Title("Is Production?").
				Value(&production),
			// Is Container (default true)
			huh.NewConfirm().
				Title("Is Container?").
				Value(&container),
			// Database version (default 17)
			huh.NewSelect[string]().
				Title("PostgreSQL Version").
				Options(dbVersions...).
				Value(&databaseVersion),
			// Database cluster (default false)
			huh.NewConfirm().
				Title("Is Database HA Cluster?").
				Value(&cluster),
			// Database host and port (default localhost:5432)
			huh.NewInput().
				Title("Please enter database host:port").
				Prompt(">").
				Value(&hostPort),
			// if cluster is true, ask for additional hosts
			huh.NewConfirm().
				Title("Start Build?").
				Value(&create),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	if !create {
		return fmt.Errorf("system initialization aborted")
	}

	fmt.Println("Creating base image for Odoo version", version, "on", localDomain)
	return o.BaseCreate(version, localDomain)
}

func (o *ODA) GetOdooConfig(version string) OdooConfig {
	for _, config := range o.OdooConfigs {
		if config.Version == version {
			return config
		}
	}
	return OdooConfig{}
}

func (o *ODA) BaseCreate(version, localDomain string) error {
	fmt.Println("Creating base image for Odoo version", version, "on", localDomain)
	config := o.GetOdooConfig(version)

	roleUpdateScript(o.EmbedFS)

	rolePreeperRepo(o.EmbedFS)

	roleUpdate()

	aptInstall(config.BaselinePackages...)

	roleOdooUser()

	roleOdooDirs()

	rolePostgresqlRepo()

	rolePostgresqlClient(o.OdooDatabase.Version)

	roleWkhtmltopdf()

	aptInstall(config.Odoobase...)

	npmInstall("rtlcss")

	roleGeoIP2DB()

	rolePaperSize()

	roleOdooService(o.EmbedFS)

	roleCaddy()

	roleCaddyService(o.EmbedFS)

	rolePGCat()

	rolePGCatService(o.EmbedFS)

	return nil
}
