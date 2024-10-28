package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func odooService(action string) error {
	fmt.Println(action + " odoo service")
	cmd := exec.Command("sudo",
		"systemctl",
		action,
		"odoo.service",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("service odoo %s failed: %w", action, err)
	}
	return nil
}

func (o *ODA) OdooStart() error {
	return odooService("start")
}

func (o *ODA) OdooStop() error {
	return odooService("start")
}

func (o *ODA) OdooRestart() error {
	return odooService("restart")
}

func (o *ODA) Logs() error {
	command := exec.Command("sudo",
		"journalctl",
		"-u",
		"odoo.service",
		"-f",
	)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("error getting logs %w", err)
	}
	return nil
}
