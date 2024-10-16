package internal

import (
	"errors"
	"fmt"
	"io"
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
