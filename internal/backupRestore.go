package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/huh"
)

func (o *ODA) Restore(any, move, full bool) error {
	var backups []string
	var addons []string

	if any {
		backups, addons = GetOdooBackups("")
	} else {
		backups, addons = GetOdooBackups(o.OdooConf.DbName)
	}

	backupOptions := []huh.Option[string]{}
	for _, backup := range backups {
		backupOptions = append(backupOptions, huh.NewOption(backup, backup))
	}
	addonOptions := []huh.Option[string]{}
	addonOptions = append(addonOptions, huh.NewOption("None", "none"))
	for _, addon := range addons {
		addonOptions = append(addonOptions, huh.NewOption(addon, addon))
	}

	var (
		backupFile string
		addonFile  string
		confirm    bool
	)

	huh.NewSelect[string]().
		Title("Odoo Backup File").
		Options(backupOptions...).
		Value(&backupFile).
		Run()

	huh.NewSelect[string]().
		Title("Odoo Addon File").
		Options(addonOptions...).
		Value(&addonFile).
		Run()

	huh.NewConfirm().
		Title("Restore Project?").
		Value(&confirm).
		Run()

	if !confirm {
		fmt.Println("restore cancelled")
		return nil
	}

	if addonFile != "none" {
		fmt.Println("restore from addon file " + addonFile)
		if err := o.restoreAddonsTar(addonFile); err != nil {
			return fmt.Errorf("restore addons tar failed %w", err)
		}
	}

	fmt.Println("restore from backup file " + backupFile)
	if err := o.restoreDBTar(backupFile, move, full); err != nil {
		return fmt.Errorf("restore db tar failed %w", err)
	}

	return nil
}

// restoreAddonsTar Restore Odoo DB addons folders
func (o *ODA) restoreAddonsTar(addonsFile string) error {
	root_dir := filepath.Join("/", "opt", "odoo")
	source := filepath.Join(root_dir, "backups", addonsFile)
	dest := filepath.Join(root_dir, "addons")

	if err := RemoveContents(dest); err != nil {
		return fmt.Errorf("remove contents failed: %w", err)
	}
	cmd := exec.Command("tar",
		"axf", source, "-C", dest, ".",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract addon files failed: %w", err)
	}
	return nil
}

// restoreDBTar Restore Odoo DB from backup
func (o *ODA) restoreDBTar(backupFile string, moveDB bool, neutralize bool) error {
	// Stop Odoo Service
	o.OdooStop()

	cwd := filepath.Join("/", "opt", "odoo")
	source := filepath.Join(cwd, "backups", backupFile)
	port, _ := strconv.Atoi(o.OdooConf.DbPort)

	odb := OdooDB{
		Hostname: o.OdooConf.DbHost,
		Port:     o.OdooConf.DbPort,
		Database: o.OdooConf.DbName,
		Username: o.OdooConf.DbUser,
		Password: o.OdooConf.DbPassword,
		Template: o.OdooConf.DbTemplate,
	}

	// drop target database
	if err := odb.DropDatabase(); err != nil {
		return fmt.Errorf("could not drop postgresql database %s error: %w", o.OdooConf.DbName, err)
	}

	// create new postgresql database
	if err := odb.CreateDatabase(); err != nil {
		return fmt.Errorf("could not create postgresql database %s error: %w", o.OdooConf.DbName, err)
	}

	// restore postgresql database
	if err := odb.RestoreDatabase(source); err != nil {
		return fmt.Errorf("could not restore postgresql database %s error: %w", o.OdooConf.DbName, err)
	}

	// restore data filestore
	// fmt.Println("restore postgresql database")
	data := filepath.Join(cwd, "data")
	if err := RemoveContents(data); err != nil {
		return fmt.Errorf("data files removal failed %w", err)
	}
	filestore := filepath.Join(data, "filestore", o.OdooConf.DbName)
	if err := os.MkdirAll(filestore, 0o755); err != nil {
		return fmt.Errorf("filestore directory creation failed %w", err)
	}
	tarCmd := exec.Command("tar",
		"axf", source, "-C", filestore, "--strip-components=2", "./filestore",
	)
	if err := tarCmd.Run(); err != nil {
		return fmt.Errorf("filestore restore failed %w", err)
	}
	// fmt.Println("restored filestore " + dbname)

	// if not moveDB then reset DBUUID and remove MCode
	if !moveDB {
		fmt.Println("neutralize the database")
		db, err := OpenDatabase(Database{
			Hostname: o.OdooConf.DbHost,
			Port:     port,
			Database: o.OdooConf.DbName,
			Username: o.OdooConf.DbUser,
			Password: o.OdooConf.DbPassword,
		})
		if err != nil {
			return fmt.Errorf("error opening database %w", err)
		}
		defer func() error {
			if err := db.Close(); err != nil {
				return fmt.Errorf("error closing database %w", err)
			}
			return nil
		}()

		db.RemoveEnterpriseCode()
		db.ChangeDBUUID()
		db.UpdateDatabaseExpirationDate()
		db.DisableBankSync()
		db.DisableFetchmail()
		db.DeactivateMailServers()
		db.DeactivateCrons()
		db.ActivateModuleUpdateNotificationCron()
		db.RemoveIRLogging()
		db.DisableProdDeliveryCarriers()
		db.DisableDeliveryCarriers()
		db.DisableIAPAccount()
		db.DisableMailTemplate()
		db.DisablePaymentGeneric()
		db.DeleteWebsiteDomains()
		db.DisableCDN()
		db.DeleteOCNProjectUUID()
		db.UnsetFirebase()
		db.RemoveMapBoxToken()

		// Social Media
		db.RemoveFacebookTokens()
		db.RemoveInstagramTokens()
		db.RemoveLinkedInTokens()
		db.RemoveTwitterTokens()
		db.RemoveYoutubeTokens()

		if neutralize {
			db.ActivateNeutralizationWatermarks()
		}
	}
	return nil
}
