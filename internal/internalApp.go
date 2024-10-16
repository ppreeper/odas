package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func moduleList(modules ...string) string {
	mods := []string{}
	for _, mod := range modules {
		mm := strings.Split(mod, ",")
		mods = append(mods, mm...)
	}
	return strings.Join(removeDuplicate(mods), ",")
}

func InstanceAppInstallUpgrade(install bool, modules ...string) error {
	iu := "-u"

	if install {
		iu = "-i"
	}

	cmd := exec.Command("odoo/odoo-bin",
		"--no-http", "--stop-after-init",
		"-c", "/opt/odoo/conf/odoo.conf",
		iu, moduleList(modules...),
	)
	cmd.Dir = "/opt/odoo"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error installing/upgrading modules %w", err)
	}

	return nil
}

func Scaffold(module string) error {
	cmd := exec.Command("odoo/odoo-bin",
		"scaffold", module, "/opt/odoo/addons",
	)
	cmd.Dir = "/opt/odoo"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error scaffolding module %w", err)
	}
	return nil
}
