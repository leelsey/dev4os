package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
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
	cmdPMS      = checkBrewPath()
	pmsIns      = "install"
	//pmsReIn     = "reinstall"
	//pmsRm       = "remove"
	pmsAlt    = "--cask"
	pmsRepo   = "tap"
	cmdGit    = "git"
	cmdASDF   = checkASDFPath()
	p10kPath  = homeDir() + ".config/p10k/"
	p10kCache = homeDir() + ".cache/p10k-" + userName()
	cmdOpt    string
)

func checkError(err error) {
	if err != nil {
		fmt.Println("\n" + lstDot + err.Error())
		os.Exit(0)
	}
}

func checkPermission() string {
	fmt.Printf("   ")
	sudoPW := exec.Command("sudo", "whoami")
	sudoPW.Env = os.Environ()
	sudoPW.Stdin = os.Stdin
	sudoPW.Stderr = os.Stderr
	whoAmI, err := sudoPW.Output()

	if err != nil {
		fmt.Println(lstDot + "Shell command sudo error: " + err.Error())
		os.Exit(0)
	}
	return string(whoAmI)
}

func checkNetStatus() bool {
	getTimeout := 10000 * time.Millisecond
	client := http.Client{
		Timeout: getTimeout,
	}
	_, err := client.Get("https://9.9.9.9")
	if err != nil {
		return false
	}
	return true
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
	checkError(err)
	return homeDirPath + "/"
}

func workingDir() string {
	workingDirPath, err := os.Getwd()
	checkError(err)
	return workingDirPath + "/"
}

func userName() string {
	workingUser, err := user.Current()
	checkError(err)
	return workingUser.Username
}

func makeDir(dirPath string) {
	err := os.MkdirAll(dirPath, 0700)
	checkError(err)
}

func makeFile(filePath, fileContents string) {
	targetFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	checkError(err)
	defer func() {
		err := targetFile.Close()
		checkError(err)
	}()
	_, err = targetFile.Write([]byte(fileContents))
	checkError(err)
}

func appendFile(filePath, fileContents string) {
	targetFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.FileMode(0600))
	checkError(err)
	defer func() {
		err := targetFile.Close()
		checkError(err)
	}()
	_, err = targetFile.Write([]byte(fileContents))
	checkError(err)
}

func removeFile(filePath string) {
	if _, errExist := os.Stat(filePath); !os.IsNotExist(errExist) {
		err := os.Remove(filePath)
		checkError(err)
	}
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

func brewIns(pkg string) {
	if _, errExist := os.Stat(brewPrefix + "Cellar/" + pkg); errors.Is(errExist, os.ErrNotExist) {
		brewIns := exec.Command(cmdPMS, pmsIns, pkg)
		if err := brewIns.Run(); err != nil {
			fmt.Println("\n" + lstDot + "Brew " + pkg + " install error: " + err.Error())
			brewIns.Stderr = os.Stderr
			os.Exit(0)
		}
	}
}

func brewInsCask(pkg, app string) {
	if _, errExist := os.Stat("/Applications/" + app + ".app"); errors.Is(errExist, os.ErrNotExist) {
		brewIns := exec.Command(cmdPMS, pmsIns, pmsAlt, pkg)
		if err := brewIns.Run(); err != nil {
			fmt.Println("\n" + lstDot + "Brew " + app + ".app install (cask) error: " + err.Error())
			brewIns.Stderr = os.Stderr
			os.Exit(0)
		}
	}
}

func asdfIns(plugin, version string) {
	if _, errExist := os.Stat(homeDir() + ".asdf/plugins/" + plugin); errors.Is(errExist, os.ErrNotExist) {
		asdfPlugin := exec.Command(cmdASDF, "plugin", "add", plugin)
		if err := asdfPlugin.Run(); err != nil {
			fmt.Println("\n" + lstDot + "Failed ASDF-VM " + plugin + " add plugin \n" +
				"   Error code: " + err.Error())
			asdfPlugin.Stderr = os.Stderr
			os.Exit(0)
		}
	}

	asdfInstall := exec.Command(cmdASDF, pmsIns, plugin, version)
	if err := asdfInstall.Run(); err != nil {
		fmt.Println("\n" + lstDot + "Failed ASDF-VM " + plugin + " (" + version + ") install\n" +
			"   Error code: " + err.Error())
		asdfInstall.Stderr = os.Stderr
		os.Exit(0)
	}

	asdfGlobal := exec.Command(cmdASDF, "global", plugin, version)
	if err := asdfGlobal.Run(); err != nil {
		fmt.Println("\n" + lstDot + "Failed ASDF-VM " + plugin + " (" + version + ") globalisation\n" +
			"   Error code: " + err.Error())
		asdfGlobal.Stderr = os.Stderr
		os.Exit(0)
	}
}

func confA4s() {
	dlA4sPath := workingDir() + ".dev4mac-alias4sh.sh"
	resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Alias4sh/main/install.sh")
	if err != nil {
		fmt.Println(lstDot + "Alias4sh‘s URL is maybe changed, please check https://github.com/leelsey/Alias4sh\n")
		os.Exit(0)
	}
	defer func() {
		err := resp.Body.Close()
		checkError(err)
	}()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	a4sInstaller, err := os.OpenFile(dlA4sPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0755))
	checkError(err)
	defer func() {
		err := a4sInstaller.Close()
		checkError(err)
	}()
	_, err = a4sInstaller.Write(rawFile)
	checkError(err)

	installA4s := exec.Command("/bin/sh", dlA4sPath)
	if err := installA4s.Run(); err != nil {
		removeFile(dlA4sPath)
		checkError(err)
	}
	removeFile(dlA4sPath)
}

func confG4s() {
	fmt.Println("\nGit global configuration")

	fmt.Println(" 1) Main branch default name changed master -> main")
	setBranchMain := exec.Command(cmdGit, "config", "--global", "init.defaultBranch", "main")
	if err := setBranchMain.Run(); err != nil {
		checkError(err)
	}

	fmt.Println(" 2) Add your information to the global git config")
	consoleReader := bufio.NewScanner(os.Stdin)
	fmt.Printf(" " + lstDot + "User name: ")
	consoleReader.Scan()
	userName := consoleReader.Text()
	fmt.Printf(" " + lstDot + "User email: ")
	consoleReader.Scan()
	userEmail := consoleReader.Text()

	setUserName := exec.Command(cmdGit, "config", "--global", "user.name", userName)
	setUserEmail := exec.Command(cmdGit, "config", "--global", "user.email", userEmail)
	if err := setUserName.Run(); err != nil {
		checkError(err)
	}
	if err := setUserEmail.Run(); err != nil {
		checkError(err)
	}

	fmt.Println(" 3) Setup git global ignore file with directories")
	ignoreDirPath := homeDir() + ".config/git/"
	makeDir(ignoreDirPath)

	ignorePath := ignoreDirPath + "gitignore_global"
	resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Git4set/main/gitignore-sample")
	if err != nil {
		fmt.Println(lstDot + "Git Ignore sample URL is maybe changed, please check https://github.com/leelsey/Git4set\n")
		os.Exit(0)
	}
	defer func() {
		err := resp.Body.Close()
		checkError(err)
	}()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	gitIgnore, err := os.OpenFile(ignorePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	checkError(err)
	defer func() {
		err := gitIgnore.Close()
		checkError(err)
	}()
	_, err = gitIgnore.Write(rawFile)
	checkError(err)

	setExcludesFile := exec.Command(cmdGit, "config", "--global", "core.excludesfile", ignorePath)
	if err := setExcludesFile.Run(); err != nil {
		fmt.Println("error2")
		checkError(err)
	}

	fmt.Println(" " + lstDot + "Make \"gitignore_global\" file in " + ignoreDirPath)
}

func p10kTerm() {
	dlP10kTerm := p10kPath + "p10k-term.zsh"
	respP10kTerm, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devsimple.zsh")
	if err != nil {
		fmt.Println(lstDot + "ZshTheme‘s URL is maybe changed, please check https://github.com/leelsey/ConfStore\n")
		os.Exit(0)
	}
	defer func() {
		err := respP10kTerm.Body.Close()
		checkError(err)
	}()
	rawFileP10kTerm, _ := ioutil.ReadAll(respP10kTerm.Body)

	confP10kTerm, err := os.OpenFile(dlP10kTerm, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	checkError(err)
	defer func() {
		err := confP10kTerm.Close()
		checkError(err)
	}()
	_, err = confP10kTerm.Write(rawFileP10kTerm)
	checkError(err)
}

func p10kiTerm2() {
	dlP10kiTerm2 := p10kPath + "p10k-iterm2.zsh"
	respP10kiTerm2, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devwork.zsh")
	if err != nil {
		fmt.Println(lstDot + "ZshTheme‘s URL is maybe changed, please check https://github.com/leelsey/ConfStore\n")
		os.Exit(0)
	}
	defer func() {
		err := respP10kiTerm2.Body.Close()
		checkError(err)
	}()
	rawFileP10kiTerm2, _ := ioutil.ReadAll(respP10kiTerm2.Body)

	confP10kiTerm2, err := os.OpenFile(dlP10kiTerm2, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	checkError(err)
	defer func() {
		err := confP10kiTerm2.Close()
		checkError(err)
	}()
	_, err = confP10kiTerm2.Write(rawFileP10kiTerm2)
	checkError(err)
}

func p10kTMUX() {
	dlP10kTMUX := p10kPath + "p10k-tmux.zsh"
	respP10kTMUX, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devhelp.zsh")
	if err != nil {
		fmt.Println(lstDot + "ZshTheme‘s URL is maybe changed, please check https://github.com/leelsey/ConfStore\n")
		os.Exit(0)
	}
	defer func() {
		err := respP10kTMUX.Body.Close()
		checkError(err)
	}()
	rawFileP10kTMUX, _ := ioutil.ReadAll(respP10kTMUX.Body)

	confP10kTMUX, err := os.OpenFile(dlP10kTMUX, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	checkError(err)
	defer func() {
		err := confP10kTMUX.Close()
		checkError(err)
	}()
	_, err = confP10kTMUX.Write(rawFileP10kTMUX)
	checkError(err)
}

func p10kEtc() {
	dlP10kEtc := p10kPath + "p10k-etc.zsh"
	respP10kEtc, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-devbegin.zsh")
	if err != nil {
		fmt.Println(lstDot + "ZshTheme‘s URL is maybe changed, please check https://github.com/leelsey/ConfStore\n")
		os.Exit(0)
	}
	defer func() {
		err := respP10kEtc.Body.Close()
		checkError(err)
	}()
	rawFileP10kEtc, _ := ioutil.ReadAll(respP10kEtc.Body)

	confP10kEtc, err := os.OpenFile(dlP10kEtc, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	checkError(err)
	defer func() {
		err := confP10kEtc.Close()
		checkError(err)
	}()
	_, err = confP10kEtc.Write(rawFileP10kEtc)
	checkError(err)
}

func iTerm2Conf() {
	dliTerm2Conf := homeDir() + "Library/Preferences/com.googlecode.iterm2.plist"
	respiTerm2Conf, err := http.Get("https://raw.githubusercontent.com/leelsey/ConfStore/main/iterm2/iTerm2.plist")
	if err != nil {
		fmt.Println(lstDot + "ZshTheme‘s URL is maybe changed, please check https://github.com/leelsey/ConfStore\n")
		os.Exit(0)
	}
	defer func() {
		err := respiTerm2Conf.Body.Close()
		checkError(err)
	}()
	rawFileiTerm2Conf, _ := ioutil.ReadAll(respiTerm2Conf.Body)

	confiTerm2Conf, err := os.OpenFile(dliTerm2Conf, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	checkError(err)
	defer func() {
		err := confiTerm2Conf.Close()
		checkError(err)
	}()
	_, err = confiTerm2Conf.Write(rawFileiTerm2Conf)
	checkError(err)
}

func installBrew() {
	dlBrewPath := workingDir() + ".dev4mac-brew.sh"
	resp, err := http.Get("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh")
	if err != nil {
		fmt.Println(lstDot + "Brew install‘s URL is maybe changed, please check https://github.com/Homebrew/install\n")
		os.Exit(0)
	}
	defer func() {
		err := resp.Body.Close()
		checkError(err)
	}()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	brewInstaller, err := os.OpenFile(dlBrewPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0755))
	checkError(err)
	defer func() {
		err := brewInstaller.Close()
		checkError(err)
	}()
	_, err = brewInstaller.Write(rawFile)
	checkError(err)

	installHomebrew := exec.Command("/bin/bash", "-c", dlBrewPath)
	installHomebrew.Env = append(os.Environ(), "NONINTERACTIVE=1")

	if err := installHomebrew.Run(); err != nil {
		removeFile(dlBrewPath)
		checkError(err)
	}
	removeFile(dlBrewPath)

	if checkBrewExists() == false {
		fmt.Println("Brew install failed, please check your system\n")
		os.Exit(0)
	}
}

func updateBrew() {
	if err := os.Chmod(brewPrefix+"share", 0755); err != nil {
		checkError(err)
	}

	updateHomebrew := exec.Command(cmdPMS, "update")
	updateBrewCore := exec.Command(cmdPMS, pmsRepo, "homebrew/core")
	updateBrewCask := exec.Command(cmdPMS, pmsRepo, "homebrew/cask")
	updateBrewCaskVersions := exec.Command(cmdPMS, pmsRepo, "homebrew/cask-versions")

	if err := updateHomebrew.Run(); err != nil {
		checkError(err)
	}
	if err := updateBrewCore.Run(); err != nil {
		checkError(err)
	}
	if err := updateBrewCask.Run(); err != nil {
		checkError(err)
	}
	if err := updateBrewCaskVersions.Run(); err != nil {
		checkError(err)
	}
}

func macBegin() {
	switch {
	case checkBrewExists() == true:
		ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
		ldBar.Suffix = " Updating homebrew..."
		ldBar.FinalMSG = " - Updated brew!\n"
		ldBar.Start()

		updateBrew()

		ldBar.Stop()
	case checkBrewExists() == false:
		fmt.Println(" - Check root permission (sudo) for install the Homebrew")
		if checkPermission() == "root\n" {
			ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
			ldBar.Suffix = " Installing homebrew..."
			ldBar.FinalMSG = " - Installed and updated brew!\n"
			ldBar.Start()

			installBrew()
			updateBrew()

			ldBar.Stop()
		} else {
			fmt.Println(lstDot + "Incorrect user, please check permission of sudo.\n" +
				lstDot + "It need sudo command of \"root\" user's permission.\n" +
				lstDot + "Now your username: " + checkPermission())
			os.Exit(0)
		}
	}
}

func macEnv() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Setting basic environment..."
	ldBar.FinalMSG = " - Completed environment!\n"
	ldBar.Start()

	newZProfile()
	newZshRC()

	profileAppend := "# HOMEBREW\n" +
		"eval \"$(" + cmdPMS + " shellenv)\"\n"
	appendFile(profilePath, profileAppend)

	ldBar.Stop()
}

func macDependency(runOpt string) {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for basic environment configuration..."
	ldBar.FinalMSG = " - Installed basic dependencies!\n"
	ldBar.Start()

	brewIns("pkg-config")
	brewIns("ca-certificates")
	brewIns("openssl@3")
	brewIns("openssl@1.1")
	brewIns("ncurses")
	brewIns("autoconf")
	brewIns("mpdecimal")
	brewIns("libyaml")
	brewIns("readline")
	brewIns("gdbm")
	brewIns("xz")
	brewIns("sqlite")

	shrcAppend := "# OPENSSL-3\n" +
		"export PATH=\"" + brewPrefix + "opt/openssl@3/bin:$PATH\"\n" +
		"export LDFLAGS=\"-L" + brewPrefix + "opt/openssl@3/lib\"\n" +
		"export CPPFLAGS=\"-I" + brewPrefix + "opt/openssl@3/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/openssl@3/lib/pkgconfig\"\n\n" +
		"# OPENSSL-1.1\n" +
		"export PATH=\"" + brewPrefix + "opt/openssl@1.1/bin:$PATH\"\n" +
		"export LDFLAGS=\"-L" + brewPrefix + "opt/openssl@1.1/lib\"\n" +
		"export CPPFLAGS=\"-I" + brewPrefix + "opt/openssl@1.1/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/openssl@1.1/lib/pkgconfig\"\n\n" +
		"# NCURSES\n" +
		"export PATH=\"" + brewPrefix + "opt/ncurses/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/ncurses/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/ncurses/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/ncurses/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)

	if runOpt == "3" || runOpt == "4" || runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewIns("pcre")
		brewIns("pcre2")
	}

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewIns("krb5")
		brewIns("gnupg")
		brewIns("curl")
		brewIns("wget")
		brewIns("gzip")
		brewIns("libzip")
		brewIns("bzip2")
		brewIns("zlib")
		brewIns("ghc")
		brewIns("ccache")
		brewIns("cabal")
		brewIns("m4")
		brewIns("automake")
		brewIns("libffi")
		brewIns("guile")
		brewIns("gnu-getopt")
		brewIns("coreutils")
		brewIns("bison")
		brewIns("libiconv")
		brewIns("icu4c")
		brewIns("re2c")
		brewIns("gd")
		brewIns("ldns")
		brewIns("html-xml-utils")
		brewIns("xmlto")
		brewIns("gmp")
		brewIns("libsodium")
		brewIns("imagemagick")
		brewIns("ghostscript")

		shrcAppend := "# KRB5\n" +
			"export PATH=\"" + brewPrefix + "opt/krb5/bin:$PATH\"\n" +
			"export PATH=\"" + brewPrefix + "opt/krb5/sbin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/krb5/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/krb5/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/krb5/lib/pkgconfig\"\n\n" +
			"# BZIP2\n" +
			"export PATH=\"" + brewPrefix + "opt/bzip2/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/bzip2/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/bzip2/include\"\n\n" +
			"# ZLIB\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/zlib/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/zlib/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/zlib/lib/pkgconfig\"\n\n" +
			"# GNU GETOPT\n" +
			"export PATH=\"" + brewPrefix + "opt/gnu-getopt/bin:$PATH\"\n\n" +
			"# COREUTILS\n" +
			"export PATH=\"" + brewPrefix + "opt/coreutils/libexec/gnubin:$PATH\"\n\n" +
			"# BISON\n" +
			"export PATH=\"" + brewPrefix + "opt/bison/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/bison/lib\"\n\n" +
			"# LIBICONV\n" +
			"export PATH=\"" + brewPrefix + "opt/libiconv/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/libiconv/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/libiconv/include\"\n\n" +
			"# ICU4C\n" +
			"export PATH=\"" + brewPrefix + "opt/icu4c/bin:$PATH\"\n" +
			"export PATH=\"" + brewPrefix + "opt/icu4c/sbin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/icu4c/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/icu4c/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/icu4c/lib/pkgconfig\"\n\n" +
			"# CURL\n" +
			"export PATH=\"" + brewPrefix + "opt/curl/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/curl/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/curl/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/curl/lib/pkgconfig\"\n\n"
		appendFile(shrcPath, shrcAppend)
	}

	ldBar.Stop()
}

func macLanguage(runOpt string) {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing computer programming language..."
	ldBar.FinalMSG = " - Installed basic languages!\n"
	ldBar.Start()

	brewIns("awk")

	if runOpt == "4" {
		brewIns("openjdk")
		brewIns("nvm")
		shrcAppend := "# JAVA\n" +
			"export PATH=\"" + brewPrefix + "opt/openjdk/bin:$PATH\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/openjdk/include\"\n\n" +
			"# NVM\n" +
			"export NVM_DIR=\"$HOME/.nvm\"\n" +
			"[ -s \"" + brewPrefix + "opt/nvm/nvm.sh\" ] && source \"" + brewPrefix + "opt/nvm/nvm.sh\"\n" +
			"[ -s \"" + brewPrefix + "opt/nvm/etc/bash_completion.d/nvm\" ] && source \"" + brewPrefix + "opt/nvm/etc/bash_completion.d/nvm\"\n\n"
		appendFile(shrcPath, shrcAppend)

		nvmIns := exec.Command("nvm", pmsIns, "--lts")
		if err := nvmIns.Run(); err != nil {
			fmt.Println("\n" + lstDot + "NVM" + " install error: " + err.Error())
			nvmIns.Stderr = os.Stderr
			os.Exit(0)
		}
	} else if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewIns("gawk")
		brewIns("perl")
		brewIns("ruby")
		brewIns("python")
		brewIns("openjdk")
		brewIns("rust")
		brewIns("go")
		brewIns("node")
		brewIns("lua")
		brewIns("php")
		brewIns("groovy")
		brewIns("kotlin")
		brewIns("scala")
		brewIns("clojure")
		brewIns("erlang")
		brewIns("elixir")
		brewIns("typescript")
		brewIns("r")
		brewIns("haskell-stack")
		brewIns("haskell-language-server")
		brewIns("stylish-haskell")

		shrcAppend := "# RUBY\n" +
			"export PATH=\"" + brewPrefix + "opt/ruby/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/ruby/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/ruby/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/ruby/lib/pkgconfig\"\n\n" +
			"# JAVA\n" +
			"export PATH=\"" + brewPrefix + "opt/openjdk/bin:$PATH\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/openjdk/include\"\n\n"
		appendFile(shrcPath, shrcAppend)
	}

	ldBar.Stop()
}

func macServer() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server..."
	ldBar.FinalMSG = " - Installed servers!\n"
	ldBar.Start()

	brewIns("httpd")
	brewIns("tomcat")
	brewIns("nginx")

	ldBar.Stop()
}

func macDatabase() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for database..."
	ldBar.FinalMSG = " - Installed databases!\n"
	ldBar.Start()

	brewIns("sqlite-analyzer")
	brewIns("postgresql")
	brewIns("mysql")
	brewIns("redis")
	brewIns("mongodb-community")
	brewIns("mongodb")

	shrcAppend := "# SQLITE3\n" +
		"export PATH=\"" + brewPrefix + "opt/sqlite/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/sqlite/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/sqlite/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/sqlite/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)

	ldBar.Stop()
}

func macDevVM() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing management developer tool version with plugin..."
	ldBar.FinalMSG = " - Installed ASDF-VM, and add basic languages!\n"
	ldBar.Start()

	brewIns("asdf")

	asdfIns("perl", "latest")
	//asdfIns("ruby", "latest")   // error
	//asdfIns("python", "latest") // error
	asdfIns("java", "openjdk-11.0.2") // JDK LTS 11
	asdfIns("java", "openjdk-17.0.2") // JDK LTS 17
	asdfIns("rust", "latest")
	asdfIns("golang", "latest")
	asdfIns("nodejs", "latest")
	asdfIns("lua", "latest")
	//asdfIns("php", "latest") // error
	asdfIns("groovy", "latest")
	asdfIns("kotlin", "latest")
	asdfIns("scala", "latest")
	asdfIns("clojure", "latest")
	//asdfIns("erlang", "latest") // error
	asdfIns("elixir", "latest")
	asdfIns("haskell", "latest")
	asdfIns("gleam", "latest")
	//asdfIns("r", "latest") // error

	shrcAppend := "# ASDF VM\n" +
		"source " + brewPrefix + "opt/asdf/libexec/asdf.sh\n" +
		"source " + homeDir() + ".asdf/plugins/java/set-java-home.zsh\n" +
		"java_macos_integration_enable = yes\n\n"
	appendFile(shrcPath, shrcAppend)

	asdfReshim := exec.Command(cmdASDF, "reshim")
	if err := asdfReshim.Run(); err != nil {
		fmt.Println("\n" + lstDot + "ASDF reshim error: " + err.Error())
		os.Exit(0)
	}

	ldBar.Stop()
}

func macTerminal(runOpt string) {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing zsh with useful tools..."
	ldBar.FinalMSG = " - Installed useful tools for terminal!\n"
	ldBar.Start()

	brewIns("zsh-completions")
	brewIns("zsh-syntax-highlighting")
	brewIns("zsh-autosuggestions")
	brewIns("z")
	brewIns("tree")
	brewIns("romkatv/powerlevel10k/powerlevel10k")

	makeFile(homeDir()+".z", "")
	makeDir(p10kPath)
	makeDir(p10kCache)

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewIns("zsh")
		brewIns("fzf")
		brewIns("tmux")
		brewIns("tmuxinator")
		brewIns("neofetch")

		iTerm2Conf()
	}

	p10kTerm()

	if runOpt == "2" || runOpt == "3" || runOpt == "4" {
		profileAppend := "# POWERLEVEL10K\n" +
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

	ldBar.Stop()
}

func macCLIApp(runOpt string) {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for CLI"
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	brewIns("diffutils")

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewIns("openssh")
		brewIns("mosh")
		brewIns("inetutils")
		brewIns("git")
		brewIns("git-lfs")
		brewIns("gh")
		brewIns("tldr")
		brewIns("diffr")
		brewIns("bat")
		brewIns("tig")
		brewIns("watchman")
		brewIns("direnv")
		brewIns("jupyterlab")

		profileAppend := "# DIRENV\n" +
			"eval \"$(direnv hook zsh)\"\n\n"
		appendFile(profilePath, profileAppend)
	}

	if runOpt == "6" || runOpt == "7" {
		brewIns("make")
		brewIns("cmake")
		brewIns("ninja")
		brewIns("maven")
		brewIns("gradle")
		brewIns("htop")
		brewIns("qemu")
		brewIns("vim")
		brewIns("neovim")
		brewIns("httpie")
		brewIns("curlie")
		brewIns("jq")
		brewIns("yq")
		brewIns("dasel")
		brewIns("asciinema")
	}

	if runOpt == "7" {
		brewIns("tor")
		brewIns("torsocks")
		brewIns("nmap")
		brewIns("radare2")
		brewIns("sleuthkit")
		brewIns("autopsy")
	}

	ldBar.Stop()
}

func macGUIApp(runOpt string) {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for GUI"
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	if runOpt != "7" {
		brewInsCask("appcleaner", "AppCleaner")
	} else if runOpt == "7" {
		brewInsCask("sensei", "Sensei")
	}

	brewInsCask("keka", "Keka")
	brewInsCask("iina", "IINA")
	brewInsCask("transmission", "Transmission")
	brewInsCask("rectangle", "Rectangle")
	brewInsCask("google-chrome", "Google Chrome")
	brewInsCask("firefox", "Firefox")
	brewInsCask("tor-browser", "Tor Browser")
	brewInsCask("spotify", "Spotify")
	brewInsCask("signal", "Signal")
	brewInsCask("slack", "Slack")
	brewInsCask("discord", "Discord")

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInsCask("jetbrains-space", "JetBrains Space")
	}

	if runOpt == "3" || runOpt == "6" || runOpt == "7" {
		brewInsCask("dropbox", "Dropbox")
		brewInsCask("dropbox-capture", "Dropbox Capture")
		brewInsCask("sketch", "Sketch")
		brewInsCask("zeplin", "Zeplin")
		brewInsCask("blender", "Blender")
		brewInsCask("obs", "OBS")
	}

	if runOpt == "3" || runOpt == "4" {
		brewInsCask("visual-studio-code", "Visual Studio Code")
		brewInsCask("atom", "Atom")
		brewInsCask("eclipse-ide", "Eclipse")
		brewInsCask("intellij-idea-ce", "IntelliJ IDEA CE")
		brewInsCask("android-studio", "Android Studio")
		brewInsCask("fork", "Fork")
	} else if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInsCask("iterm2", "iTerm")
		brewInsCask("visual-studio-code", "Visual Studio Code")
		brewInsCask("atom", "Atom")
		brewInsCask("intellij-idea", "IntelliJ IDEA")
		brewInsCask("tableplus", "TablePlus")
		brewInsCask("proxyman", "Proxyman")
		brewInsCask("postman", "Postman")
		brewInsCask("paw", "Paw")
		brewInsCask("github", "Github")
		brewInsCask("fork", "Fork")
		brewInsCask("boop", "Boop")
		brewInsCask("firefox-developer-edition", "Firefox Developer Edition")
		brewInsCask("staruml", "StarUML")
		brewInsCask("docker", "Docker")
	}

	shrcAppend := "# ANDROID STUDIO\n" +
		"export ANDROID_HOME=$HOME/Library/Android/sdk\n" +
		"export PATH=$PATH:$ANDROID_HOME/emulator\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools/bin\n" +
		"export PATH=$PATH:$ANDROID_HOME/platform-tools\n\n"
	appendFile(shrcPath, shrcAppend)

	if runOpt == "6" {
		brewInsCask("vnc-viewer", "VNC Viewer")
	} else if runOpt == "7" {
		brewInsCask("vnc-viewer", "VNC Viewer")
		brewInsCask("burp-suite", "Burp Suite Community Edition")
		brewInsCask("burp-suite-professional", "Burp Suite Professional")
		brewInsCask("imazing", "iMazing")
		brewInsCask("apparency", "Apparency")
		brewInsCask("suspicious-package", "Suspicious Package")
		brewInsCask("cutter", "Cutter")
		// Gihdra
	}

	ldBar.Stop()
}

func macGUIAppPlus(runOpt string) {
	fmt.Println(" - Check root permission (sudo) for install the GUI App")

	if checkPermission() == "root\n" {
		ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
		ldBar.Suffix = " Installing advanced tools for GUI"
		ldBar.FinalMSG = " - Installed developer utilities!\n"
		ldBar.Start()

		if runOpt == "3" || runOpt == "6" || runOpt == "7" {
			brewInsCask("loopback", "Loopback")
		}

		if runOpt == "6" || runOpt == "7" {
			brewInsCask("vmware-fusion", "VMware Fusion")
		}

		if runOpt == "7" {
			brewInsCask("wireshark", "Wireshark")
			brewInsCask("zenmap", "Zenmap")
		}

		ldBar.Stop()
	}
}

func macEnd() {
	shrcAppend := "\n######## ADD CUSTOM VALUES UNDER HERE ########\n\n\n"
	appendFile(shrcPath, shrcAppend)
}

func main() {
	fmt.Println("\nDev4mac v" + appVer + "\n")
	if checkNetStatus() == true {
		fmt.Println("\nChoose an installation option. (Recommend option is 5)\n" +
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
			_, err := fmt.Scanln(&cmdOpt)
			checkError(err)
			if cmdOpt == "1" {
				macBegin()
				macEnv()
			} else if cmdOpt == "2" {
				macBegin()
				macEnv()
				macDependency(cmdOpt)
				macLanguage(cmdOpt)
				macTerminal(cmdOpt)
				macCLIApp(cmdOpt)
			} else if cmdOpt == "3" {
				macBegin()
				macEnv()
				macDependency(cmdOpt)
				macLanguage(cmdOpt)
				macTerminal(cmdOpt)
				macCLIApp(cmdOpt)
				macGUIApp(cmdOpt)
				macGUIAppPlus(cmdOpt)
			} else if cmdOpt == "4" {
				macBegin()
				macEnv()
				macDependency(cmdOpt)
				macLanguage(cmdOpt)
				macTerminal(cmdOpt)
				macCLIApp(cmdOpt)
				macGUIApp(cmdOpt)
			} else if cmdOpt == "5" {
				macBegin()
				macEnv()
				macDependency(cmdOpt)
				macLanguage(cmdOpt)
				macServer()
				macDatabase()
				macTerminal(cmdOpt)
				macCLIApp(cmdOpt)
				macGUIApp(cmdOpt)
			} else if cmdOpt == "6" {
				macBegin()
				macEnv()
				macDependency(cmdOpt)
				macLanguage(cmdOpt)
				macServer()
				macDatabase()
				macDevVM()
				macTerminal(cmdOpt)
				macCLIApp(cmdOpt)
				macGUIApp(cmdOpt)
				macGUIAppPlus(cmdOpt)
			} else if cmdOpt == "7" {
				macBegin()
				macEnv()
				macDependency(cmdOpt)
				macLanguage(cmdOpt)
				macServer()
				macDatabase()
				macDevVM()
				macTerminal(cmdOpt)
				macCLIApp(cmdOpt)
				macGUIApp(cmdOpt)
				macGUIAppPlus(cmdOpt)
			} else if cmdOpt == "0" || cmdOpt == "q" || cmdOpt == "e" || cmdOpt == "quit" || cmdOpt == "exit" {
				os.Exit(0)
			} else {
				fmt.Println("Wrong answer. Please choose number 0-7")
				goto startOpt
			}
			break
		}
		macEnd()

		fmt.Printf("\nFinished to setup!\nEnter [Y] to set git global configuration, or enter [N] key to exit. ")
	finishOpt:
		for {
			_, err := fmt.Scan(&cmdOpt)
			checkError(err)
			if cmdOpt == "y" || cmdOpt == "Y" || cmdOpt == "yes" || cmdOpt == "Yes" || cmdOpt == "YES" {
				confG4s()
			} else if cmdOpt == "n" || cmdOpt == "N" || cmdOpt == "no" || cmdOpt == "No" || cmdOpt == "NO" {
				break
			} else {
				fmt.Printf("Wrong answer. Please enter [Y] or [N]. ")
				goto finishOpt
			}
		}
		fmt.Println("\n----------Finished!----------\n" +
			"Please RESTART your terminal!\n" +
			lstDot + "Enter this on terminal: source ~/.zprofile && source ~/.zshrc\n" +
			lstDot + "Or restart the Terminal.app by yourself.\n")
	} else {
		fmt.Println(lstDot + "Please check your internet connection and try again.\n")
	}
}
