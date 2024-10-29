package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
)

func GetOSVersion() string {
	osRelease, err := os.Open("/etc/os-release")
	if err != nil {
		fmt.Println("Error loading os-release file", err)
		return ""
	}
	defer func() {
		if err := osRelease.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(osRelease)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`^VERSION_ID="(.+)"$`)
		if re.MatchString(line) {
			match := re.FindStringSubmatch(line)
			return match[1]
		}
	}
	return ""
}

func GetOSVersionName() (osName string, osVersionID string, osVersionCodename string) {
	osRelease, err := os.Open("/etc/os-release")
	if err != nil {
		fmt.Println("Error loading os-release file", err)
		return
	}
	defer func() {
		if err := osRelease.Close(); err != nil {
			panic(err)
		}
	}()

	reName := regexp.MustCompile(`^NAME="(.+)"$`)
	reVersionID := regexp.MustCompile(`^VERSION_ID="(.+)"$`)
	reVersionCodename := regexp.MustCompile(`^VERSION_CODENAME=(.+)$`)

	scanner := bufio.NewScanner(osRelease)
	for scanner.Scan() {
		line := scanner.Text()
		if reName.MatchString(line) {
			match := reName.FindStringSubmatch(line)
			osName = match[1]
		}
		if reVersionID.MatchString(line) {
			match := reVersionID.FindStringSubmatch(line)
			osVersionID = match[1]
		}
		if reVersionCodename.MatchString(line) {
			match := reVersionCodename.FindStringSubmatch(line)
			osVersionCodename = match[1]
		}
	}

	return strings.TrimSpace(osName), strings.TrimSpace(osVersionID), strings.TrimSpace(osVersionCodename)
}

func GetFQDN() (fqdn, hostname, domain string) {
	hostname, _ = os.Hostname()
	fqdn, _ = net.LookupCNAME(hostname)
	ff := strings.Split(fqdn, ".")
	if len(ff) == 1 {
		return fqdn + ".local", hostname, "local"
	}
	if ff[len(ff)-1] == "" {
		ff = ff[:len(ff)-1]
	}
	domain = strings.Join(ff[1:], ".")
	if domain == "" {
		domain = "local"
	}
	return hostname + "." + domain, hostname, domain
}

func AreYouSure(prompt string) bool {
	var confirm1, confirm2 bool
	huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to %s?", prompt)).
		Value(&confirm1).
		Run()
	if !confirm1 {
		return false
	}
	huh.NewConfirm().
		Title(fmt.Sprintf("Are you really sure you want to %s?", prompt)).
		Value(&confirm2).
		Run()
	if !confirm1 || !confirm2 {
		return false
	}
	return true
}

// removeDuplicate Remove Duplicate Values from Slice
func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func SelectOnly(sliceList []string, value string) []string {
	list := []string{}
	for _, item := range sliceList {
		if strings.Contains(item, value) {
			list = append(list, item)
		}
	}
	return list
}

// existsIn searches list for value
func ExistsIn[T comparable](sliceList []T, value T) bool {
	for _, item := range sliceList {
		if value == item {
			return true
		}
	}
	return false
}
