package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func (o *ODA) RepoUpdate() error {
	repoDir := filepath.Join("/", "opt", "odoo")
	fmt.Println(repoDir)
	for _, repo := range o.OdooRepos {
		fmt.Fprintln(os.Stderr, "Updating", repo)
		dest := filepath.Join(repoDir, repo)
		if Exists(dest) {
			pull := exec.Command("git", "pull", "--rebase")
			pull.Dir = dest
			if err := pull.Run(); err != nil {
				return fmt.Errorf("git pull on %s %w", repo, err)
			}
		}
	}
	return nil
}

func GetOdooBranchVersion(path string) (branch string, version string) {
	// Open an existing repository or create a new one
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}
	// Get the current branch
	head, err := r.Head()
	if err != nil {
		log.Fatal(err)
	}
	branch = strings.TrimPrefix(head.Name().String(), "refs/heads/")
	branchSplit := strings.Split(branch, "-")
	if len(branchSplit) > 1 {
		version = branchSplit[1]
	} else {
		version = branch
	}
	return branch, version
}

// func getLastTags(r *git.Repository) (plumbing.Hash, string) {
// 	iter, err := r.Tags()
// 	cobra.CheckErr(err)

// 	lastHash := plumbing.ZeroHash
// 	lastTag := ""

// 	for {
// 		ref, err := iter.Next()
// 		if errors.Is(err, io.EOF) {
// 			break
// 		}
// 		if err != nil {
// 			return plumbing.ZeroHash, ""
// 		}
// 		lastHash = ref.Hash()
// 		lastTag = ref.Name().Short()
// 	}
// 	return lastHash, lastTag
// }

// func gittag() string {
// 	appPath := ""
// 	r, err := git.PlainOpen(appPath)
// 	cobra.CheckErr(err)

// 	hashCommit, err := r.Head()
// 	cobra.CheckErr(err)

// 	_, t := getLastTags(r)

// 	now := time.Now()
// 	dateString := now.Format("20060102")

// 	return fmt.Sprintf("%s (%s-%s)", t, dateString, hashCommit.Hash().String()[:7])
// }

// func getGitVersions(repo string) []string {
// 	pattern := regexp.MustCompile(`/[0-9][0-9].[0-9]$`)
// 	gitout, err := exec.Command("git", "ls-remote", "https://github.com/odoo/"+repo).Output()
// 	cobra.CheckErr(err)
// 	odooVersions := []string{}
// 	gitlines := strings.Split(string(gitout), "\n")
// 	for _, line := range gitlines {
// 		if pattern.MatchString(line) {
// 			fields := strings.Split(line, "\t")
// 			if len(fields) == 2 {
// 				vers := strings.TrimLeft(fields[1], "refs/heads/")
// 				odooVersions = append(odooVersions, vers)
// 			}
// 		}
// 	}
// 	return odooVersions[len(odooVersions)-4:]
// }

// func CloneUrlDir(url, baseDir, cloneDir, username, token string) error {
// 	_, err := os.Stat(filepath.Join(baseDir, cloneDir, ".git"))
// 	if os.IsNotExist(err) {
// 		os.MkdirAll(baseDir, 0o755)
// 		_, err = git.PlainClone(filepath.Join(baseDir, cloneDir), false, &git.CloneOptions{
// 			URL:      url,
// 			Progress: os.Stdout,
// 			Auth: &http.BasicAuth{
// 				Username: username,
// 				Password: token,
// 			},
// 		})
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "error cloning repo %v\n", err)
// 		}
// 	}
// 	return nil
// }

// func RepoHeadShortCode(repo string) (string, error) {
// 	repoDir := filepath.Join(viper.GetString("dirs.repo"), repo)

// 	r, err := git.PlainOpen(repoDir)
// 	if err != nil {
// 		return "", err
// 	}

// 	refs, err := r.References()
// 	if err != nil {
// 		return "", err
// 	}
// 	var refList string
// 	if err := refs.ForEach(func(ref *plumbing.Reference) error {
// 		if ref.Type() == plumbing.SymbolicReference {
// 			refList = ref.Target().Short()
// 		}
// 		return nil
// 	}); err != nil {
// 		return "", fmt.Errorf("refs.ForEach on %s %w", repo, err)
// 	}
// 	return refList, nil
// }

// func RepoShortCodes(repo string) ([]string, error) {
// 	repoDir := filepath.Join(viper.GetString("dirs.repo"), repo)

// 	r, err := git.PlainOpen(repoDir)
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	refs, err := r.References()
// 	if err != nil {
// 		return []string{}, err
// 	}
// 	refList := []float64{}
// 	if err := refs.ForEach(func(ref *plumbing.Reference) error {
// 		if ref.Type() != plumbing.SymbolicReference &&
// 			strings.HasPrefix(ref.Name().Short(), "origin") {
// 			refSplit := strings.Split(ref.Name().Short(), "/")
// 			shortName := refSplit[len(refSplit)-1]
// 			if !strings.HasPrefix(shortName, "master") &&
// 				!strings.HasPrefix(shortName, "staging") &&
// 				!strings.HasPrefix(shortName, "saas") &&
// 				!strings.HasPrefix(shortName, "tmp") {
// 				val, _ := strconv.ParseFloat(shortName, 32)
// 				refList = append(refList, val)
// 			}
// 		}
// 		return nil
// 	}); err != nil {
// 		return []string{}, fmt.Errorf("refs.ForEach on %s %w", repo, err)
// 	}
// 	slices.Sort(refList)
// 	slices.Reverse(refList)
// 	shortRefs := []string{}
// 	for _, ref := range refList[0:4] {
// 		shortRefs = append(shortRefs, strconv.FormatFloat(ref, 'f', 1, 64))
// 	}
// 	return shortRefs, nil
// }

// func cloneOdooRepos() {
// 	var version string
// 	var create bool

// 	versions := getGitVersions("odoo")
// 	// versions := OdooVersions
// 	versionOptions := []huh.Option[string]{}
// 	for _, version := range versions {
// 		versionOptions = append(versionOptions, huh.NewOption(version, version))
// 	}

// 	if len(versionOptions) == 0 {
// 		fmt.Fprintln(os.Stderr, "no more branches to clone")
// 		return
// 	}

// 	form := huh.NewForm(
// 		huh.NewGroup(
// 			huh.NewSelect[string]().
// 				Title("Available Odoo Branches").
// 				Options(versionOptions...).
// 				Value(&version),

// 			huh.NewConfirm().
// 				Title("Clone Branch?").
// 				Value(&create),
// 		),
// 	)
// 	if err := form.Run(); err != nil {
// 		fmt.Fprintf(os.Stderr, "odoo version form error %v", err)
// 		return
// 	}

// 	if !create {
// 		return
// 	}

// 	repoDir := filepath.Join("/", "opt", "odoo")

// 	for _, repo := range OdooRepos {
// 		dest := filepath.Join(repoDir, repo)

// 		fmt.Println("git", "clone", "https://github.com/odoo/"+repo, dest)

// 		cloner := exec.Command("git", "clone", "https://github.com/odoo/"+repo, dest)
// 		if err := cloner.Run(); err != nil {
// 			fmt.Fprintf(os.Stderr, "git clone %s %v", repo, err)
// 			return
// 		}
// 	}
// }
