package internal

import (
	"os"
	"strings"

	"github.com/dimiro1/banner"
	"github.com/ppreeper/str"
)

func cText(color, msg string) string {
	return color + msg + "{{ .AnsiColor.Default }}"
}

func (o *ODA) Welcome() error {
	tRed := "{{ .AnsiColor.BrightRed }}"
	tMagenta := "{{ .AnsiColor.Magenta }}"

	fqdn, _, _ := GetFQDN()

	virtType := " " + GetInstanceType()

	osName, osVersion, _ := GetOSVersionName()
	osversionstring := strings.TrimSpace(osName + " " + osVersion)
	var odooVersion string
	if Exists("/opt/odoo/odoo") {
		_, odooVersion = GetOdooBranchVersion("/opt/odoo/odoo")
	}

	exampleCommands := []struct {
		Cmd  string
		Help string
	}{
		{Cmd: "odoo-bin shell", Help: "Open an Odoo shell(odoo/odoo-bin shell -c /opt/odoo/conf/odoo.conf)[IPython]"},
		{Cmd: "odas update", Help: "Update modules in the database"},
		{Cmd: "odas restart", Help: "Restart Odoo.sh services"},
		{Cmd: "odas psql", Help: "Open PostgreSQL shell"},
		{Cmd: "odoosh-storage", Help: "Check the storage usage of your instance's container filesystem (ncdu /home/odoo/data/filestore/qoc-innovations-artes-main-11599063)"},
		{Cmd: "odas logs", Help: "Navigate in your instance's odoo.log file"},
	}

	welcomeTemplate := cText(tMagenta, `{{ .Title "ODAS" "rectangles" 0 }}`) + "\n" + cText(tMagenta, o.Version) + "\n\n"
	welcomeTemplate += "You are connected to your " + cText(tRed, "<production>") + cText(tRed, virtType) + " instance " + cText(tRed, fqdn) + "\n" +
		"running " + cText(tRed, "Odoo "+odooVersion) + " on " + cText(tRed, osversionstring) + "\n\n"
	welcomeTemplate += "Overview of useful commands:\n\n"

	cmdLen := 0
	for _, cmd := range exampleCommands {
		if len(cmd.Cmd) > cmdLen {
			cmdLen = len(cmd.Cmd)
		}
	}

	for _, cmd := range exampleCommands {
		welcomeTemplate += str.RJustLen("$ ", 3) + cText(tMagenta, str.LJustLen(cmd.Cmd, cmdLen+2)) + cmd.Help + "\n"
	}
	welcomeTemplate += "\n"

	isEnabled := true
	isColorEnabled := true
	banner.InitString(os.Stdout, isEnabled, isColorEnabled, welcomeTemplate)

	return nil
}
