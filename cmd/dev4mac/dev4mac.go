package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	appVer      = "0.1"
	lstDot      = " • "
	shrcPath    = homeDir() + ".zshrc"
	profilePath = homeDir() + ".zprofile"
	arm64Path   = "/opt/homebrew/"
	amd64Path   = "/usr/local/"
	brewPrefix  = checkBrewPrefix()
	cmdRoot     = "sudo"
	cmdPMS      = checkBrewPath()
	cmdGit      = "/usr/bin/git"
	pmsIns      = "install"
	//pmsReIn     = "reinstall"
	//pmsRm       = "remove"
	pmsAlt    = "--cask"
	pmsRepo   = "tap"
	cmdASDF   = checkASDFPath()
	p10kPath  = homeDir() + ".config/p10k/"
	p10kCache = homeDir() + ".cache/p10k-" + userName()
	macLdBar  = spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	//macPgBar =
	clrReset  = "\033[0m"
	clrRed    = "\033[31m"
	clrGreen  = "\033[32m"
	clrYellow = "\033[33m"
	clrBlue   = "\033[34m"
	clrPurple = "\033[35m"
	clrCyan   = "\033[36m"
	clrGrey   = "\033[37m"
	//clrWhite  = "\033[97m"
	//clrBlack  = "\033[30m"
)

func checkError(err error, msg string) {
	defer func() {
		if recv := recover(); recv != nil {
			fmt.Println("\n"+lstDot, recv)
		}
	}()
	if err != nil {
		log.Fatalln(clrRed + "Error >> " + msg + " (" + err.Error() + ")\n")
	}
}

func checkCmdError(err error, msg, pkg string) {
	defer func() {
		if recv := recover(); recv != nil {
			fmt.Println("\n"+lstDot, recv)
		}
	}()
	if err != nil {
		fmt.Println(fmt.Errorf(clrRed + "Error >> " + msg + " " + clrYellow + pkg + clrReset + " (" + err.Error() + ")"))
	}
}

func checkPermission() {
	sudoPW := exec.Command("sudo", "whoami")
	sudoPW.Env = os.Environ()
	sudoPW.Stdin = os.Stdin
	sudoPW.Stderr = os.Stderr
	whoAmI, err := sudoPW.Output()
	fmt.Print("\033[1A\033[K")
	checkError(err, "Failed to get sudo permission")

	if string(whoAmI) != "root\n" {
		log.Fatalln(clrRed + "Error >> " + clrReset + "Incorrect user, please check permission of sudo.\n" +
			lstDot + "It need sudo command of \"" + clrRed + "root" + clrReset + "\" user's permission.\n" +
			lstDot + "Working username: " + string(whoAmI))
	}
}

func checkNetStatus() {
	getTimeout := 10000 * time.Millisecond
	client := http.Client{
		Timeout: getTimeout,
	}

	_, err := client.Get("https://9.9.9.9")
	if err != nil {
		log.Fatalln(errors.New("\n" + lstDot + "Please check your internet connection and try again.\n"))
	}
}

func checkBrewPath() string {
	switch runtime.GOARCH {
	case "amd64":
		return amd64Path + "bin/brew"
	}
	return arm64Path + "bin/brew"
}

func checkBrewPrefix() string {
	switch runtime.GOARCH {
	case "amd64":
		return amd64Path
	}
	return arm64Path
}

func checkASDFPath() string {
	asdfPath := "opt/asdf/libexec/bin/asdf"
	switch runtime.GOARCH {
	case "amd64":
		return amd64Path + asdfPath
	}
	return arm64Path + asdfPath
}

func checkBrewExists() bool {
	if _, err := os.Stat(cmdPMS); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func homeDir() string {
	homeDirPath, err := os.UserHomeDir()
	checkError(err, "Failed to get home directory")
	return homeDirPath + "/"
}

func workingDir() string {
	workingDirPath, err := os.Getwd()
	checkError(err, "Failed to get working directory")
	return workingDirPath + "/"
}

func userName() string {
	workingUser, err := user.Current()
	checkError(err, "Failed to get current user")
	return workingUser.Username
}

func makeDir(dirPath string) {
	err := os.MkdirAll(dirPath, 0700)
	checkError(err, "Failed to make directory")
}

func makeFile(filePath, fileContents string) {
	targetFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	checkError(err, "Failed to get file information to make new file from \""+filePath+"\"")
	defer func() {
		err := targetFile.Close()
		checkError(err, "Failed to finish make file to \""+filePath+"\"")
	}()
	_, err = targetFile.Write([]byte(fileContents))
	checkError(err, "Failed to fill in information to \""+filePath+"\"")
}

func appendFile(filePath, fileContents string) {
	targetFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.FileMode(0600))
	checkError(err, "Failed to get file information to append contents from \""+filePath+"\"")
	defer func() {
		err := targetFile.Close()
		checkError(err, "Failed to finish append ceontents to \""+filePath+"\"")
	}()
	_, err = targetFile.Write([]byte(fileContents))
	checkError(err, "Failed to append contents to \""+filePath+"\"")
}

func removeFile(filePath string) {
	if _, errExist := os.Stat(filePath); !os.IsNotExist(errExist) {
		err := os.Remove(filePath)
		checkError(err, "Failed to remove file \""+filePath+"\"")
	}
}

func downloadFile(filePath, urlPath string) {
	resp, err := http.Get(urlPath)
	checkError(err, "Failed to connect "+urlPath)

	defer func() {
		errBodyClose := resp.Body.Close()
		checkError(errBodyClose, "Failed to download from "+urlPath)
	}()

	rawFile, err := io.ReadAll(resp.Body)
	checkError(err, "Failed to read file information from "+urlPath)

	makeFile(filePath, string(rawFile))
}

func unzipFile(srcPath, destPath string) error {
	reader, err := zip.OpenReader(srcPath)
	checkError(err, "Failed to open zip file")
	defer func() {
		err := reader.Close()
		checkError(err, "Failed to close zip file")
	}()

	errMkDir := os.MkdirAll(destPath, 0755)
	checkError(errMkDir, "Failed to make directory")

	extractFile := func(srcFile *zip.File) error {
		rc, err := srcFile.Open()
		checkError(err, "Failed to open zip file")
		defer func() {
			err := rc.Close()
			checkError(err, "Failed to close zip file")
		}()

		destPath := filepath.Join(destPath, srcFile.Name)

		if !strings.HasPrefix(destPath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", destPath)
		}

		if srcFile.FileInfo().IsDir() {
			errMkDest := os.MkdirAll(destPath, srcFile.Mode())
			checkError(errMkDest, "Failed to make directory")
		} else {
			errMkDest := os.MkdirAll(filepath.Dir(destPath), srcFile.Mode())
			checkError(errMkDest, "Failed to make directory")
			destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcFile.Mode())
			checkError(err, "Failed to open file")
			defer func() {
				err := destFile.Close()
				checkError(err, "Failed to close file")
			}()

			_, err = io.Copy(destFile, rc)
			checkError(err, "Failed to copy zip file")
		}
		return nil
	}

	for _, files := range reader.File {
		err := extractFile(files)
		checkError(err, "Failed to extract zip file")
	}

	return nil
}

func newZProfile() {
	fileContents := "#    ___________  _____   ____  ______ _____ _      ______ \n" +
		"#   |___  /  __ \\|  __ \\ / __ \\|  ____|_   _| |    |  ____|\n" +
		"#      / /| |__) | |__) | |  | | |__    | | | |    | |__   \n" +
		"#     / / |  ___/|  _  /| |  | |  __|   | | | |    |  __|  \n" +
		"#    / /__| |    | | \\ \\| |__| | |     _| |_| |____| |____ \n" +
		"#   /_____|_|    |_|  \\_\\\\____/|_|    |_____|______|______|\n#\n" +
		"#  " + userName() + "’s zsh profile\n\n"
	makeFile(profilePath, fileContents)
}

func newZshRC() {
	fileContents := "#   ______ _____ _    _ _____   _____\n" +
		"#  |___  // ____| |  | |  __ \\ / ____|\n" +
		"#     / /| (___ | |__| | |__) | |\n" +
		"#    / /  \\___ \\|  __  |  _  /| |\n" +
		"#   / /__ ____) | |  | | | \\ \\| |____\n" +
		"#  /_____|_____/|_|  |_|_|  \\_\\\\_____|\n#\n" +
		"#  " + userName() + "’s zsh run commands\n\n"
	makeFile(shrcPath, fileContents)
}

func brewRepository(repo string) {
	brewRepo := exec.Command(cmdPMS, pmsRepo, repo)
	err := brewRepo.Run()
	checkCmdError(err, "Brew failed to add ", repo)
}

func brewUpdate() {
	updateHomebrew := exec.Command(cmdPMS, "update", "--auto-update")
	err := updateHomebrew.Run()
	checkCmdError(err, "Brew failed to", "update repositories")
}

func brewUpgrade() {
	upgradeHomebrew := exec.Command(cmdPMS, "upgrade", "--greedy")
	err := upgradeHomebrew.Run()
	checkCmdError(err, "Brew failed to", "upgrade packages")
}

func brewCleanup() {
	upgradeHomebrew := exec.Command(cmdPMS, "cleanup", "--prune=all", "-nsd")
	err := upgradeHomebrew.Run()
	checkCmdError(err, "Brew failed to", "cleanup old packages")
}

func brewRemoveCache() {
	upgradeHomebrew := exec.Command("rm", "-rf", "\"$(brew --cache)\"")
	err := upgradeHomebrew.Run()
	checkCmdError(err, "Brew failed to", "remove cache")
}

func brewInstall(pkg string) {
	if _, errExist := os.Stat(brewPrefix + "Cellar/" + pkg); errors.Is(errExist, os.ErrNotExist) {
		brewUpdate()

		brewIns := exec.Command(cmdPMS, pmsIns, pkg)
		brewIns.Stderr = os.Stderr
		err := brewIns.Run()
		checkCmdError(err, "Brew failed to install", pkg)
	}
}

func brewCask(pkg, app string) {
	if _, errExist := os.Stat("/Applications/" + app + ".app"); errors.Is(errExist, os.ErrNotExist) {
		brewUpdate()

		brewIns := exec.Command(cmdPMS, pmsIns, pmsAlt, pkg)
		err := brewIns.Run()
		checkCmdError(err, "Brew failed to install cask", pkg)
	}
}

func brewCaskSudo(pkg, app, path string) {
	macLdBar.Stop()

	fmt.Println(lstDot + "Check root permission (sudo) for install " + app)
	checkPermission()
	fmt.Print("\033[1A\033[K")

	macLdBar.Start()

	if _, errExist := os.Stat(path); errors.Is(errExist, os.ErrNotExist) {
		brewUpdate()

		brewIns := exec.Command(cmdPMS, pmsIns, pmsAlt, pkg)
		err := brewIns.Run()
		checkCmdError(err, "Brew failed to install cask", pkg)
	}

}

func asdfInstall(plugin, version string) {
	if _, errExist := os.Stat(homeDir() + ".asdf/plugins/" + plugin); errors.Is(errExist, os.ErrNotExist) {
		asdfPlugin := exec.Command(cmdASDF, "plugin", "add", plugin)
		err := asdfPlugin.Run()
		checkCmdError(err, "ASDF-VM failed to add", plugin)
	}

	asdfIns := exec.Command(cmdASDF, pmsIns, plugin, version)
	errIns := asdfIns.Run()
	checkCmdError(errIns, "ASDF-VM", plugin)

	asdfGlobal := exec.Command(cmdASDF, "global", plugin, version)
	errConf := asdfGlobal.Run()
	checkCmdError(errConf, "ASDF-VM failed to install", plugin)
}

func addJavaHome(tgVer, lnVer string) {
	tgHead := brewPrefix + "opt/openjdk"
	tgTail := " /libexec/openjdk.jdk"
	lnDir := "/Library/Java/JavaVirtualMachines/openjdk"

	if _, errExist := os.Stat(brewPrefix + "Cellar/openjdk" + tgVer); errors.Is(errExist, os.ErrNotExist) {
		lnJavaHome := exec.Command(cmdRoot, "ln", "-sfn", tgHead+tgVer+tgTail, lnDir+lnVer+".jdk")
		lnJavaHome.Stderr = os.Stderr
		err := lnJavaHome.Run()
		checkCmdError(err, "Add failed to java home", "OpenJDK")
	}
}

func confA4s() {
	dlA4sPath := workingDir() + ".dev4mac-alias4sh.sh"

	//resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Alias4sh/main/install.sh")
	//checkError(err, "Alias4sh‘s URL is maybe changed, please check https://github.com/leelsey/Alias4sh")
	//
	//defer func() {
	//	errBodyClose := resp.Body.Close()
	//	checkError(errBodyClose, "Failed to close response body")
	//}()
	//rawFile, err := io.ReadAll(resp.Body)
	//checkError(err, "Failed to read response body")
	//
	//makeFile(dlA4sPath, string(rawFile))

	downloadFile(dlA4sPath, "https://raw.githubusercontent.com/leelsey/Alias4sh/main/install.sh")

	installA4s := exec.Command("/bin/sh", dlA4sPath)
	if err := installA4s.Run(); err != nil {
		removeFile(dlA4sPath)
		checkError(err, "Failed to install Alias4sh")
	}

	removeFile(dlA4sPath)
}

func confG4s() {
	fmt.Println("\nGit global configuration")

	fmt.Println(" 1) Main branch default name changed master -> main")
	setBranchMain := exec.Command(cmdGit, "config", "--global", "init.defaultBranch", "main")
	errBranchMain := setBranchMain.Run()
	checkError(errBranchMain, "Failed to change branch default name (master -> main)")

	fmt.Println(" 2) Add your information to the global git config")
	consoleReader := bufio.NewScanner(os.Stdin)
	fmt.Printf(" " + lstDot + "User name: ")
	consoleReader.Scan()
	gitName := consoleReader.Text()
	fmt.Printf(" " + lstDot + "User email: ")
	consoleReader.Scan()
	gitEmail := consoleReader.Text()

	setUserName := exec.Command(cmdGit, "config", "--global", "user.name", gitName)
	errUserName := setUserName.Run()
	checkError(errUserName, "Failed to set git user name")
	setUserEmail := exec.Command(cmdGit, "config", "--global", "user.email", gitEmail)
	errUserEmail := setUserEmail.Run()
	checkError(errUserEmail, "Failed to set git user email")

	fmt.Println(" 3) Setup git global ignore file with directories")

	ignoreDirPath := homeDir() + ".config/git/"
	makeDir(ignoreDirPath)

	ignorePath := ignoreDirPath + "gitignore_global"
	//resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Git4set/main/gitignore-sample")
	//checkError(err, "Failed to download git ignore file, please check https://github.com/leelsey/Git4set\n")
	//
	//defer func() {
	//	err := resp.Body.Close()
	//	checkError(err, "Failed to close git ignore response body")
	//}()
	//rawFile, _ := ioutil.ReadAll(resp.Body)

	//makeFile(ignorePath, string(rawFile))

	downloadFile(ignorePath, "https://raw.githubusercontent.com/leelsey/Git4set/main/gitignore-sample")

	setExcludesFile := exec.Command(cmdGit, "config", "--global", "core.excludesfile", ignorePath)
	errExcludesFile := setExcludesFile.Run()
	checkError(errExcludesFile, "Failed to set git global ignore file")

	fmt.Println(" " + lstDot + "Make \"gitignore_global\" file in " + ignoreDirPath)

	fmt.Println("\n" + lstDot + "Check git global configuration")
	contentGitConf, err := os.ReadFile(homeDir() + ".gitconfig")
	checkError(err, "Failed to get git config file")
	fmt.Println(string(contentGitConf))
}

func p10kTerm() {
	dlP10kTerm := p10kPath + "p10k-term.zsh"

	//respP10kTerm, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devsimple.zsh")
	//checkError(err, "Failed to download p10k-term.zsh, please check https://github.com/leelsey/ConfStore")
	//
	//defer func() {
	//	err := respP10kTerm.Body.Close()
	//	checkError(err, "Failed to close p10k-term.zsh response body")
	//}()
	//rawFileP10kTerm, _ := ioutil.ReadAll(respP10kTerm.Body)
	//
	//confP10kTerm, err := os.OpenFile(dlP10kTerm, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	//checkError(err, "Failed to create p10k-term.zsh")
	//defer func() {
	//	err := confP10kTerm.Close()
	//	checkError(err, "Failed to close p10k-term.zsh")
	//}()
	//_, err = confP10kTerm.Write(rawFileP10kTerm)
	//checkError(err, "Failed to write p10k-term.zsh")

	downloadFile(dlP10kTerm, "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devsimple.zsh")
}

func p10kiTerm2() {
	dlP10kiTerm2 := p10kPath + "p10k-iterm2.zsh"

	//respP10kiTerm2, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devwork.zsh")
	//checkError(err, "Failed to download p10k-iterm2.zsh, please check https://github.com/leelsey/ConfStore")
	//defer func() {
	//	err := respP10kiTerm2.Body.Close()
	//	checkError(err, "Failed to close p10ki-iterm2.zsh response body")
	//}()
	//rawFileP10kiTerm2, _ := ioutil.ReadAll(respP10kiTerm2.Body)
	//
	//confP10kiTerm2, err := os.OpenFile(dlP10kiTerm2, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	//checkError(err, "Failed to create p10ki-iterm2.zsh")
	//defer func() {
	//	err := confP10kiTerm2.Close()
	//	checkError(err, "Failed to close p10ki-iterm2.zsh")
	//}()
	//_, err = confP10kiTerm2.Write(rawFileP10kiTerm2)
	//checkError(err, "Failed to write p10ki-iterm2.zsh")

	downloadFile(dlP10kiTerm2, "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devwork.zsh")
}

func p10kTMUX() {
	dlP10kTMUX := p10kPath + "p10k-tmux.zsh"

	//respP10kTMUX, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devhelp.zsh")
	//checkError(err, "Failed to download p10k-tumx.zsh, please check https://github.com/leelsey/ConfStore")
	//defer func() {
	//	err := respP10kTMUX.Body.Close()
	//	checkError(err, "Failed to close p10k-tmux.zsh response body")
	//}()
	//rawFileP10kTMUX, _ := ioutil.ReadAll(respP10kTMUX.Body)
	//
	//confP10kTMUX, err := os.OpenFile(dlP10kTMUX, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	//checkError(err, "Failed to create p10k-tmux.zsh")
	//defer func() {
	//	err := confP10kTMUX.Close()
	//	checkError(err, "Failed to close p10k-tmux.zsh")
	//}()
	//_, err = confP10kTMUX.Write(rawFileP10kTMUX)
	//checkError(err, "Failed to write p10k-tmux.zsh")

	downloadFile(dlP10kTMUX, "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devhelp.zsh")
}

func p10kEtc() {
	dlP10kEtc := p10kPath + "p10k-etc.zsh"

	//respP10kEtc, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devbegin.zsh")
	//checkError(err, "Failed to download p10k-etc.zsh, please check https://github.com/leelsey/ConfStore")
	//defer func() {
	//	err := respP10kEtc.Body.Close()
	//	checkError(err, "Failed to close p10k-etc.zsh response body")
	//}()
	//rawFileP10kEtc, _ := ioutil.ReadAll(respP10kEtc.Body)
	//
	//confP10kEtc, err := os.OpenFile(dlP10kEtc, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	//checkError(err, "Failed to create p10k-etc.zsh")
	//defer func() {
	//	err := confP10kEtc.Close()
	//	checkError(err, "Failed to close p10k-etc.zsh")
	//}()
	//_, err = confP10kEtc.Write(rawFileP10kEtc)
	//checkError(err, "Failed to write p10k-etc.zsh")

	downloadFile(dlP10kEtc, "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devbegin.zsh")
}

func iTerm2Conf() {
	dliTerm2Conf := homeDir() + "Library/Preferences/com.googlecode.iterm2.plist"

	//respiTerm2Conf, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/iterm2/iTerm2.plist")
	//checkError(err, "Failed to download iTerm2 configure file, please check https://github.com/leelsey/ConfStore")
	//defer func() {
	//	err := respiTerm2Conf.Body.Close()
	//	checkError(err, "Failed to close iTerm2.plist response body")
	//}()
	//rawFileiTerm2Conf, _ := ioutil.ReadAll(respiTerm2Conf.Body)
	//
	//confiTerm2Conf, err := os.OpenFile(dliTerm2Conf, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	//checkError(err, "Failed to create iTerm2.plist")
	//defer func() {
	//	err := confiTerm2Conf.Close()
	//	checkError(err, "Failed to close iTerm2.plist")
	//}()
	//_, err = confiTerm2Conf.Write(rawFileiTerm2Conf)
	//checkError(err, "Failed to write iTerm2.plist")

	downloadFile(dliTerm2Conf, "https://raw.githubusercontent.com/leelsey/ConfStore/main/iterm2/iTerm2.plist")
}

func installBrew() {
	dlBrewPath := workingDir() + ".dev4mac-brew.sh"

	//resp, err := http.Get("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh")
	//checkError(err, "Failed to download Homebrew install file, please check https://github.com/Homebrew/install")
	//defer func() {
	//	err := resp.Body.Close()
	//	checkError(err, "Failed to close Homebrew install.sh response body")
	//}()
	//rawFile, _ := ioutil.ReadAll(resp.Body)
	//
	//brewInstaller, err := os.OpenFile(dlBrewPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0755))
	//checkError(err, "Failed to create install.sh")
	//defer func() {
	//	err := brewInstaller.Close()
	//	checkError(err, "Failed to close install.sh")
	//}()
	//_, err = brewInstaller.Write(rawFile)
	//checkError(err, "Failed to write install.sh")

	downloadFile(dlBrewPath, "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh")

	installHomebrew := exec.Command("/bin/bash", "-c", dlBrewPath)
	installHomebrew.Env = append(os.Environ(), "NONINTERACTIVE=1")

	if err := installHomebrew.Run(); err != nil {
		removeFile(dlBrewPath)
		checkError(err, "Failed to install Homebrew")
	}
	removeFile(dlBrewPath)

	if checkBrewExists() == false {
		log.Fatalln(clrRed + "Error >>" + clrReset + " Installed brew failed, please check your system\n")
	}
}

func macBegin() {
	if checkBrewExists() == true {
		macLdBar.Suffix = " Updating homebrew... "
		macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "update homebrew!\n"
		macLdBar.Start()
	} else {
		//fmt.Println(lstDot + "Check root permission (sudo) for install the Homebrew")
		fmt.Println(clrYellow + "Check permission " + clrReset + "(sudo) for install Homebrew\n")
		checkPermission()
		fmt.Print("\033[1A\033[K")

		macLdBar.Suffix = " Installing homebrew... "
		macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install and update homebrew!\n"
		macLdBar.Start()

		installBrew()
	}

	err := os.Chmod(brewPrefix+"share", 0755)
	checkError(err, "Failed to change permissions on "+brewPrefix+"share to 755")

	brewUpdate()
	brewRepository("homebrew/core")
	brewRepository("homebrew/cask")
	brewRepository("homebrew/cask-versions")
	brewUpgrade()

	macLdBar.Stop()
}

func macEnv() {
	macLdBar.Suffix = " Setting basic environment... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "setup zsh environment!\n"
	macLdBar.Start()

	newZProfile()
	newZshRC()

	profileAppend := "# HOMEBREW\n" +
		"eval \"$(" + cmdPMS + " shellenv)\"\n"
	appendFile(profilePath, profileAppend)

	macLdBar.Stop()
}

func macDependency(runOpt string) {
	macLdBar.Suffix = " Installing dependencies... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install dependencies!\n"
	macLdBar.Start()

	brewInstall("pkg-config")
	brewInstall("ca-certificates")
	brewInstall("ncurses")
	brewInstall("openssl@3")
	brewInstall("openssl@1.1")
	brewInstall("readline")
	brewInstall("autoconf")
	brewInstall("automake")
	brewInstall("mpdecimal")
	brewInstall("utf8proc")
	brewInstall("m4")
	brewInstall("gmp")
	brewInstall("mpfr")
	brewInstall("gettext")
	brewInstall("jpeg-turbo")
	brewInstall("libtool")
	brewInstall("libevent")
	brewInstall("libffi")
	brewInstall("libtiff")
	brewInstall("libvmaf")
	brewInstall("libpng")
	brewInstall("libyaml")
	brewInstall("giflib")
	brewInstall("xz")
	brewInstall("gdbm")
	brewInstall("sqlite")
	brewInstall("lz4")
	brewInstall("zstd")
	brewInstall("hiredis")
	brewInstall("berkeley-db")
	brewInstall("asciidoctor")
	brewInstall("freetype")
	brewInstall("pcre")
	brewInstall("pcre2")

	shrcAppend := "# NCURSES\n" +
		"export PATH=\"" + brewPrefix + "opt/ncurses/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/ncurses/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/ncurses/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/ncurses/lib/pkgconfig\"\n\n" +
		"# OPENSSL-3\n" +
		"export PATH=\"" + brewPrefix + "opt/openssl@3/bin:$PATH\"\n" +
		"export LDFLAGS=\"-L" + brewPrefix + "opt/openssl@3/lib\"\n" +
		"export CPPFLAGS=\"-I" + brewPrefix + "opt/openssl@3/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/openssl@3/lib/pkgconfig\"\n\n" +
		"# OPENSSL-1.1\n" +
		"export PATH=\"" + brewPrefix + "opt/openssl@1.1/bin:$PATH\"\n" +
		"export LDFLAGS=\"-L" + brewPrefix + "opt/openssl@1.1/lib\"\n" +
		"export CPPFLAGS=\"-I" + brewPrefix + "opt/openssl@1.1/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/openssl@1.1/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)

	if runOpt != "2" && runOpt != "3" {
		brewInstall("bash")
		brewInstall("zsh")
		brewInstall("ccache")
		brewInstall("perl")
		brewInstall("ruby")
		brewInstall("python@3.10")
		brewInstall("openjdk")
		brewInstall("ghc")
		brewInstall("cabal-install")
	}

	if runOpt == "6" || runOpt == "7" {
		brewInstall("krb5")
		brewInstall("libsodium")
		brewInstall("nettle")
		brewInstall("coreutils")
		brewInstall("ldns")
		brewInstall("gzip")
		brewInstall("bzip2")
		brewInstall("fop")
		brewInstall("jasper")
		brewInstall("little-cms2")
		brewInstall("imath")
		brewInstall("openldap")
		brewInstall("openexr")
		brewInstall("openjpeg")
		brewInstall("jpeg-xl")
		brewInstall("webp")
		brewInstall("rtmpdump")
		brewInstall("aom")
		brewInstall("npth")
		brewInstall("screenresolution")
		brewInstall("gnu-getopt")
		brewInstall("brotli")
		brewInstall("bison")
		brewInstall("tcl-tk")
		brewInstall("gawk") // awk
		brewInstall("swig")
		brewInstall("re2c")
		brewInstall("icu4c")
		brewInstall("bdw-gc")
		brewInstall("guile")
		brewInstall("wxwidgets")
		brewInstall("sphinx-doc")
		brewInstall("docbook")
		brewInstall("docbook-xsl")
		brewInstall("xmlto")
		brewInstall("html-xml-utils")
		brewInstall("shared-mime-info")
		brewInstall("x265")
		brewInstall("zlib")
		brewInstall("glib")
		brewInstall("libgpg-error")
		brewInstall("libgcrypt")
		brewInstall("libunistring")
		brewInstall("libatomic_ops")
		brewInstall("libiconv")
		brewInstall("libidn")
		brewInstall("libidn2")
		brewInstall("libssh2")
		brewInstall("libnghttp2")
		brewInstall("libxml2")
		brewInstall("libtasn1")
		brewInstall("libxslt")
		brewInstall("libavif")
		brewInstall("libzip")
		brewInstall("libde265")
		brewInstall("libheif")
		brewInstall("libksba")
		brewInstall("libusb")
		brewInstall("liblqr")
		brewInstall("libomp")
		brewInstall("libassuan")
		brewInstall("gnutls")
		brewInstall("gd")
		brewInstall("ghostscript")
		brewInstall("imagemagick")
		brewInstall("pinentry")
		brewInstall("p11-kit")
		brewInstall("gnupg")

		shrcAppend := "# KRB5\n" +
			"export PATH=\"" + brewPrefix + "opt/krb5/bin:$PATH\"\n" +
			"export PATH=\"" + brewPrefix + "opt/krb5/sbin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/krb5/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/krb5/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/krb5/lib/pkgconfig\"\n\n" +
			"# COREUTILS\n" +
			"export PATH=\"" + brewPrefix + "opt/coreutils/libexec/gnubin:$PATH\"\n\n" +
			"export PATH=\"" + brewPrefix + "opt/gnu-getopt/bin:$PATH\"\n\n" +
			"# BZIP2\n" +
			"export PATH=\"" + brewPrefix + "opt/bzip2/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/bzip2/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/bzip2/include\"\n\n" +
			"# GNU GETOPT\n" +
			"# BISON\n" +
			"export PATH=\"" + brewPrefix + "opt/bison/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/bison/lib\"\n\n" +
			"# TCL-TK\n" +
			"export PATH=\"" + brewPrefix + "opt/tcl-tk/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/tcl-tk/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/tcl-tk/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/tcl-tk/lib/pkgconfig\"\n\n" +
			"# ICU4C\n" +
			"export PATH=\"" + brewPrefix + "opt/icu4c/bin:$PATH\"\n" +
			"export PATH=\"" + brewPrefix + "opt/icu4c/sbin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/icu4c/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/icu4c/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/icu4c/lib/pkgconfig\"\n\n" +
			"# ZLIB\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/zlib/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/zlib/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/zlib/lib/pkgconfig\"\n\n" +
			"# LIBICONV\n" +
			"export PATH=\"" + brewPrefix + "opt/libiconv/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/libiconv/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/libiconv/include\"\n\n" +
			"# LIBXML2\n" +
			"export PATH=\"" + brewPrefix + "opt/libxml2/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/libxml2/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/libxml2/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/libxml2/lib/pkgconfig\"\n\n" +
			"# LIBXSLT\n" +
			"export PATH=\"" + brewPrefix + "opt/libxslt/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/libxslt/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/libxslt/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/libxslt/lib/pkgconfig\"\n\n"
		appendFile(shrcPath, shrcAppend)
	}

	macLdBar.Stop()
}

func macLanguage(runOpt string) {
	macLdBar.Suffix = " Installing computer programming language... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install languages!\n"
	macLdBar.Start()

	shrcAppend := "# CCACHE\n" +
		"export PATH=\"" + brewPrefix + "opt/ccache/libexec:$PATH\"\n\n" +
		"# RUBY\n" +
		"export PATH=\"" + brewPrefix + "opt/ruby/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/ruby/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/ruby/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/ruby/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)

	if runOpt == "2" || runOpt == "3" {
		shrcAppend := "# JAVA\n" +
			"export PATH=\"" + brewPrefix + "opt/openjdk/bin:$PATH\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/openjdk/include\"\n\n"
		appendFile(shrcPath, shrcAppend)
	} else if runOpt == "4" || runOpt == "5" {
		brewInstall("openjdk@8")
		brewInstall("openjdk@11")
		brewInstall("openjdk@17")
		brewInstall("go")
		brewInstall("php")
		brewInstall("nvm")
		brewInstall("pyenv")
		brewInstall("pyenv-virtualenv")

		shrcAppend := "# NVM\n" +
			"export NVM_DIR=\"$HOME/.nvm\"\n" +
			"[ -s \"" + brewPrefix + "opt/nvm/nvm.sh\" ] && source \"" + brewPrefix + "opt/nvm/nvm.sh\"\n" +
			"[ -s \"" + brewPrefix + "opt/nvm/etc/bash_completion.d/nvm\" ] && source \"" +
			brewPrefix + "opt/nvm/etc/bash_completion.d/nvm\"\n\n" +
			"# PYENV" +
			"export PYENV_ROOT=\"$HOME/.pyenv\"\n" +
			"export PATH=\"$PYENV_ROOT/bin:$PATH\"\n" +
			"eval \"$(pyenv init --path)\"\n" +
			"eval \"$(pyenv init -)\"\n\n"
		appendFile(shrcPath, shrcAppend)

		nvmIns := exec.Command("nvm", pmsIns, "--lts")
		nvmIns.Stderr = os.Stderr
		err := nvmIns.Run()
		checkCmdError(err, "NVM failed to install", "LTS")
	} else if runOpt == "6" || runOpt == "7" {
		brewInstall("openjdk@8")
		brewInstall("openjdk@11")
		brewInstall("openjdk@17")
		brewInstall("go")
		brewInstall("rust")
		brewInstall("node")
		brewInstall("lua")
		brewInstall("php")
		brewInstall("groovy")
		brewInstall("kotlin")
		brewInstall("scala")
		brewInstall("clojure")
		brewInstall("erlang")
		brewInstall("elixir")
		brewInstall("typescript")
		brewInstall("haskell-stack")
		brewInstall("haskell-language-server")
		brewInstall("stylish-haskell")
	}

	if runOpt == "4" || runOpt == "5" || runOpt == "6" || runOpt == "7" {
		checkPermission()

		addJavaHome("", "")
		addJavaHome("@17", "-17")
		addJavaHome("@11", "-11")
		addJavaHome("@8", "-8")
	}

	macLdBar.Stop()
}

func macServer(runOpt string) {
	macLdBar.Suffix = " Installing developing tools for server... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install servers!\n"
	macLdBar.Start()

	if runOpt == "3" {
	} else {
		brewInstall("httpd")
	}
	brewInstall("httpd")
	brewInstall("tomcat")
	brewInstall("nginx")

	macLdBar.Stop()
}

func macDatabase(runOpt string) {
	macLdBar.Suffix = " Installing developing tools for database... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install databases!\n"
	macLdBar.FinalMSG = " - Installed databases for !\n"
	macLdBar.Start()

	if runOpt == "3" {
		brewInstall("sqlite-analyzer")
		brewInstall("mysql")
	} else {
		brewInstall("sqlite-analyzer")
		brewInstall("postgresql")
		brewInstall("mysql")
		brewInstall("redis")
		brewRepository("mongodb/brew")
		brewInstall("mongodb-community")
	}

	shrcAppend := "# SQLITE3\n" +
		"export PATH=\"" + brewPrefix + "opt/sqlite/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/sqlite/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/sqlite/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/sqlite/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)

	macLdBar.Stop()
}

func macDevVM() {
	macLdBar.Suffix = " Installing developer tools version management tool with plugin... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install ASDF-VM with languages!\n"
	macLdBar.Start()

	brewInstall("asdf")

	asdfInstall("perl", "latest")
	//asdfInstall("ruby", "latest")   // error
	//asdfInstall("python", "latest") // error
	asdfInstall("java", "openjdk-11.0.2") // JDK LTS 11
	asdfInstall("java", "openjdk-17.0.2") // JDK LTS 17
	asdfInstall("rust", "latest")
	asdfInstall("golang", "latest")
	asdfInstall("nodejs", "latest")
	asdfInstall("lua", "latest")
	//asdfInstall("php", "latest") // error
	asdfInstall("groovy", "latest")
	asdfInstall("kotlin", "latest")
	asdfInstall("scala", "latest")
	asdfInstall("clojure", "latest")
	//asdfInstall("erlang", "latest") // error
	asdfInstall("elixir", "latest")
	asdfInstall("haskell", "latest")
	asdfInstall("gleam", "latest")

	shrcAppend := "# ASDF VM\n" +
		"source " + brewPrefix + "opt/asdf/libexec/asdf.sh\n" +
		"source " + homeDir() + ".asdf/plugins/java/set-java-home.zsh\n" +
		"java_macos_integration_enable = yes\n\n"
	appendFile(shrcPath, shrcAppend)

	asdfReshim := exec.Command(cmdASDF, "reshim")
	err := asdfReshim.Run()
	checkCmdError(err, "ASDF failed to", "reshim")

	macLdBar.Stop()
}

func macTerminal(runOpt string) {
	macLdBar.Suffix = " Installing zsh with useful tools... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install and configure for terminal!\n"
	macLdBar.Start()

	brewInstall("zsh-completions")
	brewInstall("zsh-syntax-highlighting")
	brewInstall("zsh-autosuggestions")
	brewInstall("z")
	brewInstall("tree")
	brewInstall("romkatv/powerlevel10k/powerlevel10k")

	makeFile(homeDir()+".z", "")
	makeDir(p10kPath)
	makeDir(p10kCache)

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInstall("fzf")
		brewInstall("tmux")
		brewInstall("tmuxinator")
		brewInstall("neofetch")

		iTerm2Conf()
	}

	p10kTerm()

	if runOpt == "2" || runOpt == "3" || runOpt == "4" {
		profileAppend := "# ZSH\n" +
			"export SHELL=zsh\n\n" +
			"# POWERLEVEL10K\n" +
			"source " + brewPrefix + "opt/powerlevel10k/powerlevel10k.zsh-theme\n" +
			"if [[ -r \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\" ]]; then\n" +
			"  source \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\"\n" +
			"fi\n" +
			"[[ ! -f " + p10kPath + "p10k-terminal.zsh ]] || source " + p10kPath + "p10k-terminal.zsh\n\n"
		appendFile(profilePath, profileAppend)
	} else if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		p10kiTerm2()
		p10kTMUX()
		p10kEtc()

		profileAppend := "# ZSH\n" +
			"export SHELL=zsh\n\n" +
			"# POWERLEVEL10K\n" +
			"source " + brewPrefix + "opt/powerlevel10k/powerlevel10k.zsh-theme\n" +
			"if [[ -r \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\" ]]; then\n" +
			"  source \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\"\n" +
			"fi\n" +
			"if [[ -d /Applications/iTerm.app ]]; then\n" +
			"  if [[ $TERM_PROGRAM = \"Apple_Terminal\" ]]; then\n" +
			"    [[ ! -f " + p10kPath + "p10k-term.zsh ]] || source " + p10kPath + "p10k-term.zsh\n" +
			"  elif [[ $TERM_PROGRAM = \"iTerm.app\" ]]; then\n" +
			"    echo ''; neofetch --bold off\n" +
			"    [[ ! -f " + p10kPath + "p10k-iterm2.zsh ]] || source " + p10kPath + "p10k-iterm2.zsh\n" +
			"  elif [[ $TERM_PROGRAM = \"tmux\" ]]; then\n" +
			"    echo ''; neofetch --bold off\n" +
			"    [[ ! -f " + p10kPath + "p10k-tmux.zsh ]] || source " + p10kPath + "p10k-tmux.zsh\n" +
			"  else\n" +
			"    [[ ! -f " + p10kPath + "p10k-etc.zsh ]] || source " + p10kPath + "p10k-etc.zsh\n" +
			"  fi\n" +
			"else\n" +
			"  [[ ! -f " + p10kPath + "p10k-term.zsh ]] || source " + p10kPath + "p10k-term.zsh\n" +
			"fi\n\n"
		appendFile(profilePath, profileAppend)
	}

	profileAppend := "# ZSH-COMPLETIONS\n" +
		"if type brew &>/dev/null; then\n" +
		"  mv () { command mv \"$@\" ; }\n" +
		"  FPATH=" + brewPrefix + "share/zsh-completions:$FPATH\n" +
		"  autoload -Uz compinit\n" +
		"  compinit\n" +
		"fi\n\n" +
		"# ZSH SYNTAX HIGHLIGHTING\n" +
		"source " + brewPrefix + "share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh\n\n" +
		"# ZSH AUTOSUGGESTIONS\n" +
		"source " + brewPrefix + "share/zsh-autosuggestions/zsh-autosuggestions.zsh\n\n" +
		"# Z\n" +
		"source " + brewPrefix + "etc/profile.d/z.sh\n\n" +
		"# Edit\n" +
		"export EDITOR=/usr/bin/vi\n" +
		"edit () { $EDITOR \"$@\" }\n" +
		"#vi () { $EDITOR \"$@\" }\n\n"
	appendFile(profilePath, profileAppend)

	confA4s()

	macLdBar.Stop()
}

func macCLIApp(runOpt string) {
	macLdBar.Suffix = " Installing CLI applications... "
	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install CLI applications!\n"
	macLdBar.Start()

	brewInstall("unzip")
	brewInstall("diffutils")
	brewInstall("transmission-cli")

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInstall("curl")
		brewInstall("wget")
		brewInstall("openssh")
		brewInstall("mosh")
		brewInstall("inetutils")
		brewInstall("git")
		brewInstall("git-lfs")
		brewInstall("gh")
		brewInstall("tldr")
		brewInstall("diffr")
		brewInstall("bat")
		brewInstall("tig")
		brewInstall("watchman")
		brewInstall("direnv")
		brewInstall("jupyterlab")

		shrcAppend := "# CURL\n" +
			"export PATH=\"" + brewPrefix + "opt/curl/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/curl/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/curl/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/curl/lib/pkgconfig\"\n\n" +
			"# DIRENV\n" +
			"eval \"$(direnv hook zsh)\"\n\n"
		appendFile(shrcPath, shrcAppend)
	}

	if runOpt == "6" || runOpt == "7" {
		brewInstall("make")
		brewInstall("cmake")
		brewInstall("ninja")
		brewInstall("maven")
		brewInstall("gradle")
		brewInstall("rustup-init")
		brewInstall("htop")
		brewInstall("qemu")
		brewInstall("vim")
		brewInstall("neovim")
		brewInstall("httpie")
		brewInstall("curlie")
		brewInstall("jq")
		brewInstall("yq")
		brewInstall("dasel")
		brewInstall("asciinema")
	}

	if runOpt == "7" {
		brewInstall("tor")
		brewInstall("torsocks")
		brewInstall("nmap")
		brewInstall("radare2")
		brewInstall("sleuthkit")
		brewInstall("autopsy")
		brewInstall("virustotal-cli")
	}

	macLdBar.Stop()
}

func macGUIApp(runOpt string) {
	macLdBar.Suffix = " Installing GUI applications... "
	macLdBar.Start()

	if runOpt != "7" {
		brewCask("appcleaner", "AppCleaner")
	} else if runOpt == "7" {
		brewCask("sensei", "Sensei")
	}

	brewCask("keka", "Keka")
	brewCask("iina", "IINA")
	brewCask("transmission", "Transmission")
	brewCask("rectangle", "Rectangle")
	brewCask("google-chrome", "Google Chrome")
	brewCask("firefox", "Firefox")
	brewCask("tor-browser", "Tor Browser")
	brewCask("spotify", "Spotify")
	brewCask("signal", "Signal")
	brewCask("discord", "Discord")
	brewCask("slack", "Slack")
	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewCask("jetbrains-space", "JetBrains Space")
	}

	if runOpt == "3" || runOpt == "6" || runOpt == "7" {
		brewCask("dropbox", "Dropbox")
		brewCask("dropbox-capture", "Dropbox Capture")
		brewCask("sketch", "Sketch")
		brewCask("zeplin", "Zeplin")
		brewCask("blender", "Blender")
		brewCask("obs", "OBS")
		brewCaskSudo("loopback", "Loopback", "/Applications/Loopback.app")
		if runOpt == "3" {
			brewCaskSudo("blackhole-64ch", "BlackHole (64ch)", "/Library/Audio/Plug-Ins/HAL/BlackHoleXch.driver")
		}
	}

	if runOpt == "3" || runOpt == "4" {
		brewCask("visual-studio-code", "Visual Studio Code")
		brewCask("atom", "Atom")
		brewCask("eclipse-ide", "Eclipse")
		brewCask("intellij-idea-ce", "IntelliJ IDEA CE")
		brewCask("android-studio", "Android Studio")
		brewCask("fork", "Fork")
	} else if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewCask("iterm2", "iTerm")
		brewCask("visual-studio-code", "Visual Studio Code")
		brewCask("atom", "Atom")
		brewCask("intellij-idea", "IntelliJ IDEA")
		brewCask("tableplus", "TablePlus")
		brewCask("proxyman", "Proxyman")
		brewCask("postman", "Postman")
		brewCask("paw", "Paw")
		brewCask("boop", "Boop")
		brewCask("github", "Github")
		brewCask("fork", "Fork")
		brewCaskSudo("vmware-fusion", "VMware Fusion", "/Applications/VMware Fusion.app")
		brewCask("docker", "Docker")
		brewCask("firefox-developer-edition", "Firefox Developer Edition")
		brewCask("staruml", "StarUML")
		brewCask("vnc-viewer", "VNC Viewer")
		brewCask("forklift", "ForkLift")
	}

	shrcAppend := "# ANDROID STUDIO\n" +
		"export ANDROID_HOME=$HOME/Library/Android/sdk\n" +
		"export PATH=$PATH:$ANDROID_HOME/emulator\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools/bin\n" +
		"export PATH=$PATH:$ANDROID_HOME/platform-tools\n\n"
	appendFile(shrcPath, shrcAppend)

	if runOpt == "7" {
		brewCask("burp-suite", "Burp Suite Community Edition")
		brewCask("burp-suite-professional", "Burp Suite Professional")
		brewCaskSudo("wireshark", "Wireshark", "/Applications/Wireshark.app")
		brewCaskSudo("zenmap", "Zenmap", "/Applications/Zenmap.app")
		brewCask("imazing", "iMazing")
		brewCask("apparency", "Apparency")
		brewCask("suspicious-package", "Suspicious Package")
		brewCask("cutter", "Cutter")
		// Ghidra
		// Hopper Disassembler
	}

	macLdBar.FinalMSG = " - " + clrGreen + "Succeed " + clrReset + "install GUI applications!\n"
	macLdBar.Stop()
}

func macEnd() {
	brewCleanup()
	brewRemoveCache()

	shrcAppend := "\n######## ADD CUSTOM VALUES UNDER HERE ########\n\n\n"
	appendFile(shrcPath, shrcAppend)
}

func main() {
	var (
		beginOpt string
		endOpt   string
	)

	fmt.Println(clrPurple + "\nDev4mac " + clrGrey + "v" + appVer + clrReset)
	checkNetStatus()

	fmt.Println("\nChoose an installation option. " + clrGrey + "(Recommend option is 5)\n" + clrReset +
		"For a detailed explanation of each option and a list of installations, " +
		"read the README if need manual or explain: https://github.com/leelsey/Dev4os.\n" +
		"\t1. Minimal\n" +
		"\t2. Basic\n" +
		"\t3. Creator\n" +
		"\t4. Beginner\n" +
		"\t5. Developer\n" +
		"\t6. Professional\n" +
		"\t7. Specialist\n" +
		"\t0. Exit\n")
startOpt:
	for {
		fmt.Printf("Select command: ")
		_, errBeginOpt := fmt.Scanln(&beginOpt)
		if errBeginOpt != nil {
			beginOpt = "Null"
		}

		if beginOpt == "1" {
			fmt.Println(lstDot + "Select option " + clrBlue + "1\n" + clrReset + lstDot + clrBlue + "Minimal" +
				clrReset + ": setup homebrew with configure shell.")
			macBegin()
			macEnv()
		} else if beginOpt == "2" {
			fmt.Println(lstDot + "Select option " + clrBlue + "2\n" + clrReset + lstDot + clrBlue + "Basic" +
				clrReset + ": setup Homebrew with configure Shell, then install Dependencies, Languages and " +
				"Terminal/CLI applications with set basic preferences.")
			macBegin()
			macEnv()
			macDependency(beginOpt)
			macLanguage(beginOpt)
			macTerminal(beginOpt)
			macCLIApp(beginOpt)
		} else if beginOpt == "3" {
			fmt.Println(lstDot + "Select option " + clrBlue + "3\n" + clrReset + lstDot + clrBlue + "Creator" +
				clrReset + ": setup Homebrew with configure Shell, then install Dependencies, Languages and " +
				"Terminal/CLI/GUI applications with set basic preferences.")
			macBegin()
			macEnv()
			macDependency(beginOpt)
			macLanguage(beginOpt)
			macTerminal(beginOpt)
			macCLIApp(beginOpt)
			macGUIApp(beginOpt)
		} else if beginOpt == "4" {
			fmt.Println(lstDot + "Select option " + clrBlue + "4\n" + clrReset + lstDot + clrBlue + "Beginner" +
				clrReset + ": setup Homebrew with configure Shell, then install Dependencies, Languages and " +
				"Terminal/CLI/GUI applications with set basic preferences.")
			macBegin()
			macEnv()
			macDependency(beginOpt)
			macLanguage(beginOpt)
			macServer(beginOpt)
			macDatabase(beginOpt)
			macTerminal(beginOpt)
			macCLIApp(beginOpt)
			macGUIApp(beginOpt)
		} else if beginOpt == "5" {
			fmt.Println(lstDot + "Select option " + clrBlue + "5\n" + clrReset + lstDot + clrBlue + "Developer" +
				clrReset + ": setup Homebrew with configure Shell, then install Dependencies, Languages, Server" +
				", Database and Terminal/CLI/GUI applications with set basic preferences.")
			macBegin()
			macEnv()
			macDependency(beginOpt)
			macLanguage(beginOpt)
			macServer(beginOpt)
			macDatabase(beginOpt)
			macTerminal(beginOpt)
			macCLIApp(beginOpt)
			macGUIApp(beginOpt)
		} else if beginOpt == "6" {
			fmt.Println(lstDot + "Select option " + clrBlue + "6\n" + clrReset + lstDot + clrBlue + "Professional" +
				clrReset + ": setup homebrew with configure shell, then install Dependencies, Danguages, Server" +
				", Database, management DevTools and Terminal/CLI/GUI applications with set basic preferences.")
			macBegin()
			macEnv()
			macDependency(beginOpt)
			macLanguage(beginOpt)
			macServer(beginOpt)
			macDatabase(beginOpt)
			macDevVM()
			macTerminal(beginOpt)
			macCLIApp(beginOpt)
			macGUIApp(beginOpt)
		} else if beginOpt == "7" {
			fmt.Println(lstDot + "Select option " + clrBlue + "7\n" + clrReset + lstDot + clrBlue + "Specialist" +
				clrReset + ": setup Homebrew with configure Shell, then install Dependencies, Languages, Server" +
				", Database, management DevTools and Terminal/CLI/GUI applications with set basic preferences.")
			macBegin()
			macEnv()
			macDependency(beginOpt)
			macLanguage(beginOpt)
			macServer(beginOpt)
			macDatabase(beginOpt)
			macDevVM()
			macTerminal(beginOpt)
			macCLIApp(beginOpt)
			macGUIApp(beginOpt)
		} else if beginOpt == "0" || beginOpt == "q" || beginOpt == "e" || beginOpt == "quit" || beginOpt == "exit" {
			fmt.Println(lstDot + "Exited Dev4mac.")
			os.Exit(0)
		} else {
			fmt.Println(fmt.Errorf(lstDot + clrYellow + beginOpt + clrReset +
				" is invalid option. Please choose number " + clrRed + "0-7" + clrReset + "."))
			goto startOpt
		}
		break
	}
	macEnd()

	fmt.Printf(clrCyan + "\nFinished to setup!" + clrReset +
		"\nEnter [Y] to set git global configuration, or enter any key to exit. ")
	_, errEndOpt := fmt.Scanln(&endOpt)
	if errEndOpt != nil {
		endOpt = "Enter"
	}

	if endOpt == "y" || endOpt == "Y" || endOpt == "yes" || endOpt == "Yes" || endOpt == "YES" {
		fmt.Println(endOpt + " pressed, so running Git4sh.")
		confG4s()
	} else {
		fmt.Println(endOpt + " pressed, so finishing Dev4mac.")
	}

	fmt.Println(clrGrey + "\n----------Finished!----------\n" + clrReset +
		"Please" + clrRed + " RESTART " + clrReset + "your terminal!\n" +
		lstDot + "Enter this on terminal: source ~/.zprofile && source ~/.zshrc\n" +
		lstDot + "Or restart the Terminal.app by yourself.\n")
}
