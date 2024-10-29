package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func (o *ODA) PSQL() error {
	pgCmd := exec.Command("psql", "-h", o.OdooConf.DbHost, "-p", o.OdooConf.DbPort, "-d", o.OdooConf.DbName, "-U", o.OdooConf.DbUser)
	pgCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", o.OdooConf.DbPassword))
	pgCmd.Stdin = os.Stdin
	pgCmd.Stdout = os.Stdout
	pgCmd.Stderr = os.Stderr
	if err := pgCmd.Run(); err != nil {
		return fmt.Errorf("failed to run psql %w", err)
	}
	return nil
}
