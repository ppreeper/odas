package internal

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func odooService(action string) {
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
		fmt.Fprintf(os.Stderr, "service odoo %s failed: %v\n", action, err)
		return
	}
}

func ServiceStop() {
	odooService("stop")
}

func ServiceStart() {
	odooService("start")
}

func ServiceRestart() {
	odooService("stop")
	time.Sleep(2 * time.Second)
	odooService("start")
}
