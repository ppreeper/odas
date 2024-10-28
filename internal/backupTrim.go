package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func GetOdooBackups(project string) (backups, addons []string) {
	root_path := "/opt/odoo"
	entries, err := os.ReadDir(filepath.Join(root_path, "backups"))
	if err != nil {
		fmt.Println(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			fname := strings.Split(entry.Name(), "__")
			if len(fname) == 2 {
				backups = append(backups, entry.Name())
			} else if len(fname) == 3 {
				addons = append(addons, entry.Name())
			}
		}
	}
	slices.Sort(backups)
	slices.Sort(addons)
	backups = RemoveDuplicate(backups)
	addons = RemoveDuplicate(addons)
	if project != "" {
		backups = SelectOnly(backups, project)
		addons = SelectOnly(addons, project)
	}
	return
}

// Trim database backups
func (o *ODA) Trim(limit int, all bool) error {
	// # Get all backup files
	backups, addons := GetOdooBackups("")

	// # Group backup files by database name
	bkpFiles := make(map[string][]string)
	for _, k := range backups {
		fname := strings.Split(k, "__")
		dname := strings.TrimSuffix(fname[1], ".tar.zst")
		curFiles := bkpFiles[dname]
		bkpFiles[dname] = append(curFiles, k)
	}
	backupKeys := make([]string, 0, len(bkpFiles))
	for k := range bkpFiles {
		backupKeys = append(backupKeys, k)
	}
	slices.Sort(backupKeys)

	rmbkp := []string{}
	for _, k := range backupKeys {
		if len(bkpFiles[k]) > limit {
			rmbkp = append(rmbkp, bkpFiles[k][:len(bkpFiles[k])-limit]...)
		}
	}

	// # Group addon files by database name
	addonFiles := make(map[string][]string)
	for _, k := range addons {
		fname := strings.Split(k, "__")
		dname := fname[1]
		curFiles := addonFiles[dname]
		addonFiles[dname] = append(curFiles, k)
	}
	addonKeys := make([]string, 0, len(addonFiles))
	for k := range addonFiles {
		addonKeys = append(addonKeys, k)
	}
	slices.Sort(addonKeys)

	rmaddons := []string{}
	for _, k := range addonKeys {
		if len(addonFiles[k]) > limit {
			rmaddons = append(rmaddons, addonFiles[k][:len(addonFiles[k])-limit]...)
		}
	}

	// Join rmbackups and rmaddons
	rmlist := []string{}
	if all {
		rmlist = append(rmlist, rmbkp...)
		rmlist = append(rmlist, rmaddons...)
	} else {
		for _, k := range rmbkp {
			if strings.Contains(k, o.OdooConf.DbName) {
				rmlist = append(rmlist, k)
			}
		}
		for _, k := range rmaddons {
			if strings.Contains(k, o.OdooConf.DbName) {
				rmlist = append(rmlist, k)
			}
		}
	}

	for _, r := range rmlist {
		backupFile := filepath.Join("/", "opt", "odoo", "backups", r)
		// fmt.Println("rm", "-f", backupFile)
		os.Remove(backupFile)
	}

	return nil
}
