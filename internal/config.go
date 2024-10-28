package internal

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/charmbracelet/huh"
	"github.com/ppreeper/str"
)

func (o *ODA) CaddyfileUpdate(domain string) error {
	hostname, _ := os.Hostname()

	caddyfile := "{{.Hostname}}.{{.Domain}} {\n" +
		"tls internal\n" +
		"reverse_proxy http://{{.Hostname}}:8069\n" +
		"reverse_proxy /websocket http://{{.Hostname}}:8072\n" +
		"reverse_proxy /longpolling/* http://{{.Hostname}}:8072\n" +
		"encode gzip zstd\n" +
		"file_server\n" +
		"log\n" +
		"}\n"
	tmpl, err := template.New("caddyfile").Parse(caddyfile)
	if err != nil {
		return err
	}

	data := struct {
		Hostname string
		Domain   string
	}{
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
	hostname, _ := os.Hostname()

	hosts := str.LJustLen("127.0.1.1", 15) + "{{.Hostname}} {{.Hostname}}.{{.Domain}}" + "\n" +
		str.LJustLen("127.0.0.1", 15) + "localhost" + "\n" +
		str.LJustLen("::1", 15) + "localhost ip6-localhost ip6-loopback" + "\n" +
		str.LJustLen("ff02::1", 15) + "ip6-allnodes" + "\n" +
		str.LJustLen("ff02::2", 15) + "ip6-allrouters" + "\n"
	tmpl, err := template.New("caddyfile").Parse(hosts)
	if err != nil {
		return err
	}

	data := struct {
		Hostname string
		Domain   string
	}{
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
