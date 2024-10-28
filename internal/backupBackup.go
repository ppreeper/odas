package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (o *ODA) Backup() error {
	// Get the current date and time
	currentTime := time.Now()
	// Format the time as a string
	timeString := currentTime.Format("2006_01_02_15_04_05")
	// 	addons
	o.dumpAddonsTar(timeString)
	// main database and filestore
	o.dumpDBTar(timeString)
	return nil
}

func (o *ODA) getAddons() []string {
	var addons []string
	addons_list := strings.Split(o.OdooConf.AddonsPath, ",")[1:]
	for _, addon := range addons_list {
		if !ExistsIn(o.OdooRepos, filepath.Base(addon)) {
			addons = append(addons, addon)
		}
	}
	return addons
}

func (o *ODA) dumpAddonsTar(bkp_prefix string) error {
	addons := o.getAddons()
	for _, addon := range addons {
		folder := strings.Replace(addon, "/opt/odoo/", "", 1)
		dirlist, _ := os.ReadDir(addon)
		if len(dirlist) != 0 {
			tar_cmd := "tar"
			bkp_file := fmt.Sprintf("%s__%s__%s.tar.zst", bkp_prefix, o.OdooConf.DbName, folder)
			file_path := filepath.Join("/opt/odoo/backups", bkp_file)
			tar_args := []string{"ahcf", file_path, "-C", addon, "."}
			cmd := exec.Command(tar_cmd, tar_args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("addons backup %s failed: %w", bkp_file, err)
			}
			fmt.Println("addons:", file_path)
		}
	}
	return nil
}

func (o *ODA) dumpDBTar(bkp_prefix string) error {
	bkp_file := fmt.Sprintf("%s__%s.tar.zst", bkp_prefix, o.OdooConf.DbName)
	dump_dir := filepath.Join("/opt/odoo/backups", fmt.Sprintf("%s__%s", bkp_prefix, o.OdooConf.DbName))
	file_path := filepath.Join("/opt/odoo/backups", bkp_file)

	// create dump_dir
	if err := os.MkdirAll(dump_dir, 0o755); err != nil {
		return fmt.Errorf("directory already exists: %w", err)
	}

	// postgresql database
	pg_cmd := exec.Command("pg_dump",
		"-h", o.OdooConf.DbHost,
		"-p", o.OdooConf.DbPort,
		"-U", o.OdooConf.DbUser,
		"--no-owner",
		"--file", filepath.Join(dump_dir, "dump.sql"),
		o.OdooConf.DbName,
	)
	pg_cmd.Env = append(pg_cmd.Env, "PGPASSWORD="+o.OdooConf.DbPassword)
	pg_cmd.Stdin = os.Stdin
	pg_cmd.Stdout = os.Stdout
	pg_cmd.Stderr = os.Stderr
	if err := pg_cmd.Run(); err != nil {
		return fmt.Errorf("could not backup postgresql database %s: %w", o.OdooConf.DbName, err)
	}

	// filestore
	filestore := filepath.Join(o.OdooConf.DataDir, "filestore", o.OdooConf.DbName)
	filestore_back := filepath.Join(dump_dir, "filestore")
	if _, err := os.Stat(filestore); err == nil {
		if err := os.Symlink(filestore, filestore_back); err != nil {
			return fmt.Errorf("symlink failed: %w", err)
		}
	}

	// create tar archive
	tar_cmd := exec.Command("tar",
		"achf", file_path, "-C", dump_dir, ".",
	)
	tar_cmd.Stdin = os.Stdin
	tar_cmd.Stdout = os.Stdout
	tar_cmd.Stderr = os.Stderr
	if err := tar_cmd.Run(); err != nil {
		return fmt.Errorf("could not backup database %s: %w", o.OdooConf.DbName, err)
	}

	// cleanup dump_dir
	if err := os.RemoveAll(dump_dir); err != nil {
		return fmt.Errorf("could not cleanup dump_dir: %w", err)
	}

	fmt.Println("odoo:", file_path)

	return nil
}
