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
		return err
	}
	if !create {
		return fmt.Errorf("system initialization aborted")
	}

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

	return nil
}
