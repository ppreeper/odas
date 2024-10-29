package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ppreeper/passhash"
)

func (o *ODA) AdminPassword() error {
	var password1, password2 string
	huh.NewInput().
		Title("Please enter  the admin password:").
		Prompt(">").
		EchoMode(huh.EchoModePassword).
		Value(&password1).
		Run()
	huh.NewInput().
		Title("Please verify the admin password:").
		Prompt(">").
		EchoMode(huh.EchoModePassword).
		Value(&password2).
		Run()
	if password1 == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if password1 != password2 {
		return fmt.Errorf("passwords entered do not match")
	}
	var confirm bool
	huh.NewConfirm().
		Title("Are you sure you want to change the admin password?").
		Affirmative("yes").
		Negative("no").
		Value(&confirm).
		Run()
	if !confirm {
		return fmt.Errorf("password change cancelled")
	}

	dbport, err := strconv.Atoi(o.OdooConf.DbPort)
	if err != nil {
		return fmt.Errorf("error getting port %w", err)
	}

	db, err := OpenDatabase(Database{
		Hostname: o.OdooConf.DbHost,
		Port:     dbport,
		Username: o.OdooConf.DbUser,
		Password: o.OdooConf.DbPassword,
		Database: o.OdooConf.DbName,
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

	// Write password to database
	passkey, err := passwordHash(password1)
	if err != nil {
		fmt.Println("password hashing error", err)
	}
	_, err = db.Exec("update res_users set password=$1 where id=2;", passkey)
	if err != nil {
		return fmt.Errorf("error updating password %w", err)
	}

	fmt.Println("admin password changed")
	return nil
}

func passwordHash(password string) (string, error) {
	passkey, err := passhash.MakePassword(password, 0, "")
	if err != nil {
		return "", fmt.Errorf("password hashing error %w", err)
	}
	return strings.TrimSpace(string(passkey)), nil
}

func (o *ODA) AdminUsername() error {
	var user1, user2 string
	huh.NewInput().
		Title("Please enter  the new admin username:").
		Prompt(">").
		Value(&user1).
		Run()
	huh.NewInput().
		Title("Please verify the new admin username:").
		Prompt(">").
		Value(&user2).
		Run()

	if user1 == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if user1 != user2 {
		return fmt.Errorf("usernames entered do not match")
	}

	dbport, err := strconv.Atoi(o.OdooConf.DbPort)
	if err != nil {
		return fmt.Errorf("error getting port %w", err)
	}

	db, err := OpenDatabase(Database{
		Hostname: o.OdooConf.DbHost,
		Port:     dbport,
		Username: o.OdooConf.DbUser,
		Password: o.OdooConf.DbPassword,
		Database: o.OdooConf.DbName,
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

	// Write username to database
	_, err = db.Exec("update res_users set login=$1 where id=2;",
		strings.TrimSpace(string(user1)))
	if err != nil {
		return fmt.Errorf("error updating username %w", err)
	}

	fmt.Println("Admin username changed to", user1)
	return nil
}

func (o *ODA) UpdateUser() error {
	dbport, err := strconv.Atoi(o.OdooConf.DbPort)
	if err != nil {
		return fmt.Errorf("error getting port %w", err)
	}

	db, err := OpenDatabase(Database{
		Hostname: o.OdooConf.DbHost,
		Port:     dbport,
		Username: o.OdooConf.DbUser,
		Password: o.OdooConf.DbPassword,
		Database: o.OdooConf.DbName,
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

	getUsersQuery := "select id,company_id,partner_id,login from res_users where active=true order by login;"
	type User struct {
		ID        int    `db:"id"`
		CompanyID int    `db:"company_id"`
		PartnerID int    `db:"partner_id"`
		Login     string `db:"login"`
	}
	users := []User{}
	stmt, err := db.Preparex(getUsersQuery)
	if err != nil {
		fmt.Println("error preparing query", err)
	}
	err = stmt.Select(&users)
	if err != nil {
		fmt.Println("error getting users", err)
	}
	if len(users) == 0 {
		fmt.Println("No active users found")
	}

	usernames := []huh.Option[int]{}
	for _, user := range users {
		usernames = append(usernames, huh.NewOption(user.Login, user.ID))
	}

	var confirm bool
	var userid int
	formUser := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Odoo User").
				Options(usernames...).
				Value(&userid),
			huh.NewConfirm().
				Title("Update User").
				Value(&confirm),
		),
	)
	if err := formUser.Run(); err != nil {
		return fmt.Errorf("error updating user %w", err)
	}
	if !confirm {
		fmt.Println("update user cancelled")
		return nil
	}

	var userSelected User
	for _, user := range users {
		if user.ID == userid {
			userSelected = user
		}
	}

	user1 := userSelected.Login
	user2 := userSelected.Login
	huh.NewInput().
		Title("Please enter the new username:").
		Prompt(">").
		Value(&user1).
		Run()
	huh.NewInput().
		Title("Please verify the new username:").
		Prompt(">").
		Value(&user2).
		Run()
	if user1 == "" && user2 == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if user1 != user2 {
		return fmt.Errorf("usernames entered do not match")
	} else {
		userSelected.Login = user1
	}

	var password1, password2 string
	huh.NewInput().
		Title("Please enter the new password:").
		Prompt(">").
		EchoMode(huh.EchoModePassword).
		Value(&password1).
		Run()
	huh.NewInput().
		Title("Please verify the new password:").
		Prompt(">").
		EchoMode(huh.EchoModePassword).
		Value(&password2).
		Run()
	if password1 == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if password1 != password2 {
		return fmt.Errorf("passwords entered do not match")
	}

	huh.NewConfirm().
		Title("Are you sure you want to update user: " + userSelected.Login).
		Affirmative("yes").
		Negative("no").
		Value(&confirm).
		Run()
	if !confirm {
		return fmt.Errorf("user update cancelled")
	}

	passkey, err := passwordHash(password1)
	if err != nil {
		fmt.Println("password hashing error", err)
	}

	updateStmt, err := db.Prepare("update res_users set login=$1, password=$2 where id=$3;")
	if err != nil {
		fmt.Println("error preparing update statement", err)
	}
	_, err = updateStmt.Exec(userSelected.Login, passkey, userSelected.ID)
	if err != nil {
		fmt.Println("error updating user", err)
	}
	fmt.Println("update user", userSelected.Login, "successful")

	return nil
}
