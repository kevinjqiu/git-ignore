package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"os/exec"
	"log"
	"bytes"
	"os/user"
	"strings"
)

var language string
var listLangs bool
var shouldAppend bool
const REMOVE_GITIGNORE_REPO_URL = "https://github.com/github/gitignore.git"

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func cloneGitIgnoreRepo(path string) {
	cloneCommand := exec.Command("git", "clone", REMOVE_GITIGNORE_REPO_URL, path)
	if err := cloneCommand.Run(); err != nil {
		log.Fatal(err)
	}
}

func getGitIgnoreFiles(path string) []string {
	var out bytes.Buffer
	lsCommand := exec.Command("ls", "-1", path)
	lsCommand.Stdout = &out

	if err := lsCommand.Run(); err != nil {
		log.Fatal(err)
	}

	result := make([]string, 0)
	for _, line := range(strings.Split(out.String(), "\n")) {
		if strings.HasSuffix(line, ".gitignore") {
			result = append(result, line)
		}
	}
	return result
}

var RootCmd = &cobra.Command{
	Use:   "git-ignore",
	Short: "Setup .gitignore file",
	Run: func(cmd *cobra.Command, args []string) {
		currentUser, _ := user.Current()
		var localGitignoreRepoPath = currentUser.HomeDir + "/.gitignore-repo"

		if !pathExists(localGitignoreRepoPath) {
			cloneGitIgnoreRepo(localGitignoreRepoPath)
		}

		out := getGitIgnoreFiles(localGitignoreRepoPath)
		for _, line := range(out) {
			if strings.HasSuffix(line, ".gitignore") {
				fmt.Println(line)
			}
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&language, "lang", "", "Generate language specific .gitignore file")
	RootCmd.Flags().BoolVarP(&shouldAppend, "append", "a", false, "Append to the current .gitignore")
	RootCmd.Flags().BoolVarP(&listLangs, "list", "", false, "List available language-specific .gitignore files")
}
