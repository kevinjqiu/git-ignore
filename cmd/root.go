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
	"io/ioutil"
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

func cloneGitIgnoreRepoIfNecessary(path string) {
	if !pathExists(path) {
		cloneCommand := exec.Command("git", "clone", REMOVE_GITIGNORE_REPO_URL, path)
		if err := cloneCommand.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func getGitIgnoreLangs(path string) []string {
	var out bytes.Buffer
	lsCommand := exec.Command("ls", "-1", path)
	lsCommand.Stdout = &out

	if err := lsCommand.Run(); err != nil {
		log.Fatal(err)
	}

	result := make([]string, 0)
	for _, line := range(strings.Split(out.String(), "\n")) {
		if strings.HasSuffix(line, ".gitignore") {
			parts := strings.Split(line, ".")
			result = append(result, parts[0])
		}
	}
	return result
}

func getGitIgnoreRepoPath() string {
	currentUser, _ := user.Current()
	return currentUser.HomeDir + "/.gitignore-repo"
}

var RootCmd = &cobra.Command{
	Use:   "git-ignore",
	Short: "Setup .gitignore file",
	Run: func(cmd *cobra.Command, args []string) {
		if listLangs == true {
			var localGitignoreRepoPath = getGitIgnoreRepoPath()
			cloneGitIgnoreRepoIfNecessary(localGitignoreRepoPath)

			langs := getGitIgnoreLangs(localGitignoreRepoPath)
			for _, lang := range(langs) {
				fmt.Println(strings.ToLower(lang))
			}
		} else if language != "" {
			var localGitignoreRepoPath = getGitIgnoreRepoPath()
			cloneGitIgnoreRepoIfNecessary(localGitignoreRepoPath)

			langs := getGitIgnoreLangs(localGitignoreRepoPath)
			for _, lang := range(langs) {
				if strings.ToLower(lang) == strings.ToLower(language) {
					templateFilePath := fmt.Sprintf("%s/%s.gitignore", localGitignoreRepoPath, lang)
					cwd, err := os.Getwd()
					if err != nil {
						os.Stderr.WriteString(fmt.Sprintf("Cannot get CWD: %s", err))
					}

					gitIgnoreFilePath := cwd + "/.gitignore"
					templateFileBody, err := ioutil.ReadFile(templateFilePath)
					if err != nil {
						os.Stderr.WriteString(fmt.Sprintf("Cannot read file %s: %s", templateFilePath, err))
						return
					}

					mode := os.O_WRONLY
					if shouldAppend == true {
						mode = mode | os.O_APPEND
					}
					f, err := os.OpenFile(gitIgnoreFilePath, mode, 0644)
					defer f.Close()
					if err != nil {
						os.Stderr.WriteString(fmt.Sprintf("Cannot open file %s for writing: %s", gitIgnoreFilePath, err))
						return
					}

					if _, err := f.Write(templateFileBody); err != nil {
						os.Stderr.WriteString(fmt.Sprintf("Error while writing %s: %s", gitIgnoreFilePath, err))
						return
					}
					return
				}
			}
			os.Stderr.WriteString(fmt.Sprintf("No template for %s", language))
		} else {
			os.Stderr.WriteString("No action provided.\nPlease see git ignore -h for help\n")
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
	RootCmd.Flags().StringVarP(&language, "lang", "g", "", "Generate language specific .gitignore file")
	RootCmd.Flags().BoolVarP(&shouldAppend, "append", "a", false, "Append to the current .gitignore")
	RootCmd.Flags().BoolVarP(&listLangs, "list", "", false, "List available language-specific .gitignore files")
}
