package servicebuilder

import (
	"bytes"
	"log"
	"strings"

	"github.com/essentier/nomockutil"
	"github.com/go-errors/errors"
)

type gitProject struct {
	projectDir string
	err        error
	stashCount int
}

func (g *gitProject) init() {
	g.runGitCmd("git", "-C", g.projectDir, "init")
}

func (g *gitProject) pull(remoteUrl string, branchName string) {
	if g.err != nil {
		return
	}

	// This will fail if the project is pushed to essentier nomock the first time. That is okay.
	nomockutil.RunCmd("git", "-C", g.projectDir, "pull", "-s", "ours", remoteUrl, branchName)
}

func (g *gitProject) applyStash() {
	if g.stashCount <= 0 {
		return
	}

	g.runGitCmd("git", "-C", g.projectDir, "stash", "apply")
}

func (g *gitProject) checkout(branchName string) {
	g.runGitCmd("git", "-C", g.projectDir, "checkout", branchName)
}

func (g *gitProject) branch(branchName string) {
	g.runGitCmd("git", "-C", g.projectDir, "branch", branchName)
}

func (g *gitProject) addAll() {
	g.runGitCmd("git", "-C", g.projectDir, "add", ".", "-A")
}

func (g *gitProject) push(remoteUrl string, branchName string) {
	g.runGitCmd("git", "-C", g.projectDir, "push", remoteUrl, branchName)
}

func (g *gitProject) commit(message string) {
	if g.err != nil {
		return
	}

	out, _, err := nomockutil.RunCmd("git", "-C", g.projectDir, "commit", "-m", message)
	if err != nil {
		if !strings.Contains(out.String(), "clean") {
			g.err = err
		}
	}
}

func (g *gitProject) stashAll() {
	out := g.runGitCmd("git", "-C", g.projectDir, "stash", "save", "-u")
	if g.err != nil {
		return
	}

	if strings.Contains(out.String(), "HEAD") {
		g.stashCount++
	}
}

func (g *gitProject) getCurrentBranch() string {
	out := g.runGitCmd("git", "-C", g.projectDir, "branch")
	if g.err != nil {
		return ""
	}

	currentBranch := ""
	branchs := strings.Split(out.String(), "\n")
	for _, branch := range branchs {
		if strings.HasPrefix(branch, "*") {
			currentBranch = strings.TrimSpace(strings.TrimPrefix(branch, "*"))
			break
		}
	}
	log.Printf("current branch: [ %v ]", currentBranch)
	if currentBranch == "" {
		g.err = errors.Errorf("Failed to find current git branch.")
	}
	return currentBranch
}

func (g *gitProject) deferredPopStashed() {
	if g.stashCount <= 0 {
		return
	}

	_, _, err := nomockutil.RunCmd("git", "-C", g.projectDir, "stash", "pop")
	if err != nil {
		log.Printf("Error when trying to pop stashed %#v", err)
	} else {
		g.stashCount--
	}
}

func (g *gitProject) deferredDeleteBranch(branchName string) {
	_, _, err := nomockutil.RunCmd("git", "-C", g.projectDir, "branch", "-D", branchName)
	if err != nil {
		log.Printf("Error when trying to delete the nomock branch %#v", err)
	}
}

func (g *gitProject) deferredCheckout(originalBranch string) {
	_, _, err := nomockutil.RunCmd("git", "-C", g.projectDir, "checkout", originalBranch)
	if err != nil {
		log.Printf("Error when trying to checkout original branch %#v", err)
	}
}

func (g *gitProject) runGitCmd(name string, args ...string) *bytes.Buffer {
	if g.err != nil {
		return nil
	}

	out, _, err := nomockutil.RunCmd(name, args...)
	if err != nil {
		g.err = err
	}

	return out
}
