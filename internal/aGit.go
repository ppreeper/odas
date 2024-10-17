package internal

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func getLastTags(r *git.Repository) (plumbing.Hash, string) {
	iter, err := r.Tags()
	cobra.CheckErr(err)

	lastHash := plumbing.ZeroHash
	lastTag := ""

	for {
		ref, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return plumbing.ZeroHash, ""
		}
		lastHash = ref.Hash()
		lastTag = ref.Name().Short()
	}
	return lastHash, lastTag
}

func gittag() string {
	appPath := ""
	r, err := git.PlainOpen(appPath)
	cobra.CheckErr(err)

	hashCommit, err := r.Head()
	cobra.CheckErr(err)

	_, t := getLastTags(r)

	now := time.Now()
	dateString := now.Format("20060102")

	return fmt.Sprintf("%s (%s-%s)", t, dateString, hashCommit.Hash().String()[:7])
}

func getGitVersions(repo string) []string {
	pattern := regexp.MustCompile(`/[0-9][0-9].[0-9]$`)
	gitout, err := exec.Command("git", "ls-remote", "https://github.com/odoo/"+repo).Output()
	cobra.CheckErr(err)
	odooVersions := []string{}
	gitlines := strings.Split(string(gitout), "\n")
	for _, line := range gitlines {
		if pattern.MatchString(line) {
			fields := strings.Split(line, "\t")
			if len(fields) == 2 {
				vers := strings.TrimLeft(fields[1], "refs/heads/")
				odooVersions = append(odooVersions, vers)
			}
		}
	}
	return odooVersions[len(odooVersions)-4:]
}
