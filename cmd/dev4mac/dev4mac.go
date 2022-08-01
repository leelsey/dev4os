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
	cmdIn       = "install"
	//cmdReIn     = "reinstall"
	//cmdRm       = "remove"
	cmdGit     = "git"
	cmdASDF    = checkASDFPath()
	asdfPlugin = "plugin"
	asdfAdd    = "add"
	asdfShim   = "reshim"
	cmdOpt     string
)

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return err != nil
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

func currentUser() string {
	userName, err := user.Current()
	checkError(err)
	return userName.Username
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

func rmFile(filePath string) {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err := os.Remove(filePath)
		checkError(err)
	}
}

func newZProfile() {
	fileContents := "# " + currentUser() + "’s profile\n\n" +
		"# ZSH\n" +
		"export SHELL=zsh\n"
	makeFile(profilePath, fileContents)
}

func newZshRC() {
	fileContents := "#   _________  _   _ ____   ____    __  __    _    ___ _   _\n" +
		"#  |__  / ___|| | | |  _ \\ / ___|  |  \\/  |  / \\  |_ _| \\ | |\n" +
		"#    / /\\___ \\| |_| | |_) | |      | |\\/| | / _ \\  | ||  \\| |\n" +
		"#   / /_ ___) |  _  |  _ <| |___   | |  | |/ ___ \\ | || |\\  |\n" +
		"#  /____|____/|_| |_|_| \\_\\\\____|  |_|  |_/_/   \\_\\___|_| \\_|\n#\n\n"
	makeFile(shrcPath, fileContents)
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
		rmFile(dlA4sPath)
		checkError(err)
	}
	rmFile(dlA4sPath)
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
	ignoreDir := homeDir() + ".config/git/"
	if err := os.MkdirAll(ignoreDir, 0755); err != nil {
		checkError(err)
	}

	ignorePath := ignoreDir + "gitignore_global"
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

	gitIgnore, err := os.OpenFile(ignorePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
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

	fmt.Println(" " + lstDot + "Make \"gitignore_global\" file in " + ignoreDir)
}

func confZshTheme() {
	dlP10kPath := homeDir() + ".p10k.zsh"
	resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Dev4os/main/cmd/dev4os/dev4p10k")
	if err != nil {
		fmt.Println(lstDot + "Dev4os‘s p10k file URL is maybe changed, please check https://github.com/leelsey/Dev4os\n")
		os.Exit(0)
	}
	defer func() {
		err := resp.Body.Close()
		checkError(err)
	}()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	p10kConf, err := os.OpenFile(dlP10kPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	checkError(err)
	defer func() {
		err := p10kConf.Close()
		checkError(err)
	}()
	_, err = p10kConf.Write(rawFile)
	checkError(err)
}

func updateBrew() {
	updateHomebrew := exec.Command(cmdPMS, "update")
	updateBrewCask := exec.Command(cmdPMS, "tap", "homebrew/cask-versions")

	if err := updateHomebrew.Run(); err != nil {
		checkError(err)
	}
	if err := updateBrewCask.Run(); err != nil {
		checkError(err)
	}
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
		rmFile(dlBrewPath)
		checkError(err)
	}
	rmFile(dlBrewPath)

	if checkBrewExists() == false {
		fmt.Println("Brew install failed, please check your system\n")
		os.Exit(0)
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
		sudoPW := exec.Command("sudo", "whoami")
		sudoPW.Env = os.Environ()
		sudoPW.Stdin = os.Stdin
		sudoPW.Stderr = os.Stderr
		whoAmI, err := sudoPW.Output()

		if err != nil {
			fmt.Println(lstDot+"Shell command sudo error: ", err)
			os.Exit(0)
		} else if string(whoAmI) == "root\n" {
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
				lstDot + "Now your username: " + string(whoAmI))
			os.Exit(0)
		}
	}
}

func macEnv() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Setting basic environment..."
	ldBar.FinalMSG = " - Completed environment!\n"
	ldBar.Start()

	confA4s()
	newZProfile()
	newZshRC()

	profileAppend := "# Alias4sh\n" +
		"source ~/.config/alias4sh/aliasrc\n" +
		"# HOMEBREW\n" +
		"eval \"$(" + cmdPMS + " shellenv)\"\n"
	appendFile(profilePath, profileAppend)
	ldBar.Stop()
}

func macGit() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing git..."
	ldBar.FinalMSG = " - Installed git!\n"
	ldBar.Start()

	brewGit := exec.Command(cmdPMS, cmdIn, cmdGit)
	brewGitLfs := exec.Command(cmdPMS, cmdIn, "git-lfs")
	if err := brewGit.Run(); err != nil {
		checkError(err)
	}
	if err := brewGitLfs.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func macTerminal() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing zsh with useful tools..."
	ldBar.FinalMSG = " - Installed useful tools for terminal!\n"
	ldBar.Start()

	brewNCurses := exec.Command(cmdPMS, cmdIn, "ncurses")
	brewSSL := exec.Command(cmdPMS, cmdIn, "openssl")
	brewZsh := exec.Command(cmdPMS, cmdIn, "zsh")
	brewZshSyntax := exec.Command(cmdPMS, cmdIn, "zsh-syntax-highlighting")
	brewZshAuto := exec.Command(cmdPMS, cmdIn, "zsh-autosuggestions")
	brewZshComp := exec.Command(cmdPMS, cmdIn, "zsh-completions")
	brewTree := exec.Command(cmdPMS, cmdIn, "tree")
	brewZshTheme := exec.Command(cmdPMS, cmdIn, "romkatv/powerlevel10k/powerlevel10k")
	if err := brewNCurses.Run(); err != nil {
		checkError(err)
	}
	if err := brewSSL.Run(); err != nil {
		checkError(err)
	}
	if err := brewZsh.Run(); err != nil {
		checkError(err)
	}
	if err := brewZshSyntax.Run(); err != nil {
		checkError(err)
	}
	if err := brewZshAuto.Run(); err != nil {
		checkError(err)
	}
	if err := brewZshComp.Run(); err != nil {
		checkError(err)
	}
	if err := brewTree.Run(); err != nil {
		checkError(err)
	}
	if err := brewZshTheme.Run(); err != nil {
		checkError(err)
	}

	shrcAppend :=
		"# ZSH SYNTAX HIGHTLIGHTING\n" +
			"source " + brewPrefix + "share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh\n\n" +
			"# ZSH AUTOSUGGESTIONS\n" +
			"source " + brewPrefix + "share/zsh-autosuggestions/zsh-autosuggestions.zsh\n\n" +
			"# POWERLEVEL10K\n" +
			"[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh\n" +
			"source " + brewPrefix + "opt/powerlevel10k/powerlevel10k.zsh-theme\n\n" +
			"# NCURSES\n" +
			"export PATH=\"" + brewPrefix + "opt/ncurses/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/ncurses/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/ncurses/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/ncurses/lib/pkgconfig\"\n\n" +
			"" +
			"# OPENSSL\n" +
			"export PATH=\"" + brewPrefix + "opt/openssl@3/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/openssl@3/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/openssl@3/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/openssl@3/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)
	ldBar.Stop()
}

func macDependency() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for development work..."
	ldBar.FinalMSG = " - Installed dependencies!\n"
	ldBar.Start()

	brewKRB5 := exec.Command(cmdPMS, cmdIn, "krb5")
	brewGnuPG := exec.Command(cmdPMS, cmdIn, "gnupg")
	brewcURL := exec.Command(cmdPMS, cmdIn, "curl")
	brewWget := exec.Command(cmdPMS, cmdIn, "wget")
	brewXZ := exec.Command(cmdPMS, cmdIn, "xz")
	brewGzip := exec.Command(cmdPMS, cmdIn, "gzip")
	brewLibzip := exec.Command(cmdPMS, cmdIn, "libzip")
	brewBzip2 := exec.Command(cmdPMS, cmdIn, "bzip2")
	brewZLib := exec.Command(cmdPMS, cmdIn, "zlib")
	brewPkgConfig := exec.Command(cmdPMS, cmdIn, "pkg-config")
	brewReadLine := exec.Command(cmdPMS, cmdIn, "readline")
	brewLibffi := exec.Command(cmdPMS, cmdIn, "libffi")
	brewGuile := exec.Command(cmdPMS, cmdIn, "guile")
	brewGnuGetOpt := exec.Command(cmdPMS, cmdIn, "gnu-getopt")
	brewCoreUtils := exec.Command(cmdPMS, cmdIn, "coreutils")
	brewBison := exec.Command(cmdPMS, cmdIn, "bison")
	brewLibIconv := exec.Command(cmdPMS, cmdIn, "libiconv")
	brewICU4C := exec.Command(cmdPMS, cmdIn, "icu4c")
	brewRe2C := exec.Command(cmdPMS, cmdIn, "re2c")
	brewGD := exec.Command(cmdPMS, cmdIn, "gd")
	brewCaCert := exec.Command(cmdPMS, cmdIn, "ca-certificates")
	brewLDNS := exec.Command(cmdPMS, cmdIn, "ldns")
	brewHTMLXMLUtils := exec.Command(cmdPMS, cmdIn, "html-xml-utils")
	brewXMLto := exec.Command(cmdPMS, cmdIn, "xmlto")
	brewGMP := exec.Command(cmdPMS, cmdIn, "gmp")
	brewLibSodium := exec.Command(cmdPMS, cmdIn, "libsodium")
	brewImageMagick := exec.Command(cmdPMS, cmdIn, "imagemagick")
	brewGhostscript := exec.Command(cmdPMS, cmdIn, "ghostscript")
	if err := brewKRB5.Run(); err != nil {
		checkError(err)
	}
	if err := brewGnuPG.Run(); err != nil {
		checkError(err)
	}
	if err := brewcURL.Run(); err != nil {
		checkError(err)
	}
	if err := brewWget.Run(); err != nil {
		checkError(err)
	}
	if err := brewXZ.Run(); err != nil {
		checkError(err)
	}
	if err := brewGzip.Run(); err != nil {
		checkError(err)
	}
	if err := brewLibzip.Run(); err != nil {
		checkError(err)
	}
	if err := brewBzip2.Run(); err != nil {
		checkError(err)
	}
	if err := brewZLib.Run(); err != nil {
		checkError(err)
	}
	if err := brewPkgConfig.Run(); err != nil {
		checkError(err)
	}
	if err := brewReadLine.Run(); err != nil {
		checkError(err)
	}
	if err := brewLibffi.Run(); err != nil {
		checkError(err)
	}
	if err := brewGuile.Run(); err != nil {
		checkError(err)
	}
	if err := brewGnuGetOpt.Run(); err != nil {
		checkError(err)
	}
	if err := brewCoreUtils.Run(); err != nil {
		checkError(err)
	}
	if err := brewBison.Run(); err != nil {
		checkError(err)
	}
	if err := brewLibIconv.Run(); err != nil {
		checkError(err)
	}
	if err := brewICU4C.Run(); err != nil {
		checkError(err)
	}
	if err := brewRe2C.Run(); err != nil {
		checkError(err)
	}
	if err := brewGD.Run(); err != nil {
		checkError(err)
	}
	if err := brewCaCert.Run(); err != nil {
		checkError(err)
	}
	if err := brewLDNS.Run(); err != nil {
		checkError(err)
	}
	if err := brewHTMLXMLUtils.Run(); err != nil {
		checkError(err)
	}
	if err := brewXMLto.Run(); err != nil {
		checkError(err)
	}
	if err := brewGMP.Run(); err != nil {
		checkError(err)
	}
	if err := brewLibSodium.Run(); err != nil {
		checkError(err)
	}
	if err := brewImageMagick.Run(); err != nil {
		checkError(err)
	}
	if err := brewGhostscript.Run(); err != nil {
		checkError(err)
	}

	shrcAppend := "# KRB5\n" +
		"export PATH=\"" + brewPrefix + "opt/krb5/bin:$PATH\"\n" +
		"export PATH=\"" + brewPrefix + "opt/krb5/sbin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/krb5/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/krb5/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/krb5/lib/pkgconfig\"\n\n" +
		"# CURL\n" +
		"export PATH=\"" + brewPrefix + "opt/curl/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/curl/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/curl/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/curl/lib/pkgconfig\"\n\n" +
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
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/icu4c/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)
	ldBar.Stop()
}

func macDevToolCLI() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for CLI"
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	brewSSH := exec.Command(cmdPMS, cmdIn, "openssh")
	brewGawk := exec.Command(cmdPMS, cmdIn, "gawk")
	brewTig := exec.Command(cmdPMS, cmdIn, "tig")
	brewJQ := exec.Command(cmdPMS, cmdIn, "jq")
	brewDirEnv := exec.Command(cmdPMS, cmdIn, "direnv")
	brewWatchman := exec.Command(cmdPMS, cmdIn, "watchman")
	brewQEMU := exec.Command(cmdPMS, cmdIn, "qemu")
	brewCCache := exec.Command(cmdPMS, cmdIn, "ccache")
	brewMake := exec.Command(cmdPMS, cmdIn, "make")
	brewVim := exec.Command(cmdPMS, cmdIn, "vim")
	brewBat := exec.Command(cmdPMS, cmdIn, "bat")
	brewGH := exec.Command(cmdPMS, cmdIn, "gh")
	if err := brewSSH.Run(); err != nil {
		checkError(err)
	}
	if err := brewGawk.Run(); err != nil {
		checkError(err)
	}
	if err := brewTig.Run(); err != nil {
		checkError(err)
	}
	if err := brewJQ.Run(); err != nil {
		checkError(err)
	}
	if err := brewDirEnv.Run(); err != nil {
		checkError(err)
	}
	if err := brewWatchman.Run(); err != nil {
		checkError(err)
	}
	if err := brewQEMU.Run(); err != nil {
		checkError(err)
	}
	if err := brewCCache.Run(); err != nil {
		checkError(err)
	}
	if err := brewMake.Run(); err != nil {
		checkError(err)
	}
	if err := brewVim.Run(); err != nil {
		checkError(err)
	}
	if err := brewBat.Run(); err != nil {
		checkError(err)
	}
	if err := brewGH.Run(); err != nil {
		checkError(err)
	}

	shrcAppend := "# DIRENV\n" +
		"eval \"$(direnv hook zsh)\"\n\n"
	appendFile(shrcPath, shrcAppend)
	ldBar.Stop()
}

func macASDF() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing ASDF-VM with plugin..."
	ldBar.FinalMSG = " - Installed ASDF-VM, and add basic languages!\n"
	ldBar.Start()

	brewASDF := exec.Command(cmdPMS, cmdIn, "asdf")
	if err := brewASDF.Run(); err != nil {
		checkError(err)
	}

	shrcAppend := "# ASDF VM\n" +
		"source " + brewPrefix + "/opt/asdf/libexec/asdf.sh\n\n"
	appendFile(shrcPath, shrcAppend)

	pluginPath := homeDir() + ".asdf/plugins/"
	if _, err := os.Stat(pluginPath + "perl"); errors.Is(err, os.ErrNotExist) {
		addASDFPerl := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "perl")
		if err := addASDFPerl.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "ruby"); errors.Is(err, os.ErrNotExist) {
		addASDFRuby := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "ruby")
		if err := addASDFRuby.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "python"); errors.Is(err, os.ErrNotExist) {
		addASDFPython := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "python")
		if err := addASDFPython.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "lua"); errors.Is(err, os.ErrNotExist) {
		addASDFLua := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "lua")
		if err := addASDFLua.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "golang"); errors.Is(err, os.ErrNotExist) {
		addASDFGo := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "golang")
		if err := addASDFGo.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "rust"); errors.Is(err, os.ErrNotExist) {
		addASDFRust := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "rust")
		if err := addASDFRust.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "nodejs"); errors.Is(err, os.ErrNotExist) {
		addASDFNode := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "nodejs")
		if err := addASDFNode.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "php"); errors.Is(err, os.ErrNotExist) {
		addASDFPHP := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "php")
		if err := addASDFPHP.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "java"); errors.Is(err, os.ErrNotExist) {
		addASDFJava := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "java")
		if err := addASDFJava.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "groovy"); errors.Is(err, os.ErrNotExist) {
		addASDFGroovy := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "groovy")
		if err := addASDFGroovy.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "kotlin"); errors.Is(err, os.ErrNotExist) {
		addASDFKotlin := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "kotlin")
		if err := addASDFKotlin.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "scala"); errors.Is(err, os.ErrNotExist) {
		addASDFScala := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "scala")
		if err := addASDFScala.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "clojure"); errors.Is(err, os.ErrNotExist) {
		addASDFClojure := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "clojure")
		if err := addASDFClojure.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "erlang"); errors.Is(err, os.ErrNotExist) {
		addASDFErlang := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "erlang")
		if err := addASDFErlang.Run(); err != nil {
			checkError(err)
		}
	}
	if _, err := os.Stat(pluginPath + "elixir"); errors.Is(err, os.ErrNotExist) {
		addASDFElixir := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "elixir")
		if err := addASDFElixir.Run(); err != nil {
			checkError(err)
		}
	}
	asdfReshim := exec.Command(cmdASDF, asdfShim)
	if err := asdfReshim.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func macServer() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server..."
	ldBar.FinalMSG = " - Installed server and database!\n"
	ldBar.Start()

	brewHTTPD := exec.Command(cmdPMS, cmdIn, "httpd")
	brewTomcat := exec.Command(cmdPMS, cmdIn, "tomcat")
	brewSQLite := exec.Command(cmdPMS, cmdIn, "sqlite")
	brewPostgreSQL := exec.Command(cmdPMS, cmdIn, "postgresql")
	brewMySQL := exec.Command(cmdPMS, cmdIn, "mysql")
	brewRedis := exec.Command(cmdPMS, cmdIn, "redis")
	if err := brewHTTPD.Run(); err != nil {
		checkError(err)
	}
	if err := brewTomcat.Run(); err != nil {
		checkError(err)
	}
	if err := brewSQLite.Run(); err != nil {
		checkError(err)
	}
	if err := brewPostgreSQL.Run(); err != nil {
		checkError(err)
	}
	if err := brewMySQL.Run(); err != nil {
		checkError(err)
	}
	if err := brewRedis.Run(); err != nil {
		checkError(err)
	}

	shrcAppend := "# SQLITE3\n" +
		"export PATH=\"" + brewPrefix + "opt/sqlite/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/sqlite/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/sqlite/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/sqlite/lib/pkgconfig\"\n\n"
	appendFile(shrcPath, shrcAppend)
	ldBar.Stop()
}

func macLanguage() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing computer programming language..."
	ldBar.FinalMSG = " - Installed basic languages!\n"
	ldBar.Start()

	brewPerl := exec.Command(cmdPMS, cmdIn, "perl")
	brewRuby := exec.Command(cmdPMS, cmdIn, "ruby")
	brewPython := exec.Command(cmdPMS, cmdIn, "python@3.10")
	fixPython := exec.Command(cmdPMS, "link", "--overwrite", "python@3.10")
	brewLua := exec.Command(cmdPMS, cmdIn, "lua")
	brewGo := exec.Command(cmdPMS, cmdIn, "go")
	brewRust := exec.Command(cmdPMS, cmdIn, "rust")
	brewNode := exec.Command(cmdPMS, cmdIn, "node")
	brewTS := exec.Command(cmdPMS, cmdIn, "typescript")
	brewPHP := exec.Command(cmdPMS, cmdIn, "php")
	brewJDK := exec.Command(cmdPMS, cmdIn, "openjdk")
	brewGroovy := exec.Command(cmdPMS, cmdIn, "groovy")
	brewKotlin := exec.Command(cmdPMS, cmdIn, "kotlin")
	brewScala := exec.Command(cmdPMS, cmdIn, "scala")
	brewMaven := exec.Command(cmdPMS, cmdIn, "maven")
	brewGradle := exec.Command(cmdPMS, cmdIn, "gradle")
	brewClojure := exec.Command(cmdPMS, cmdIn, "clojure")
	brewErlang := exec.Command(cmdPMS, cmdIn, "erlang")
	brewElixir := exec.Command(cmdPMS, cmdIn, "elixir")
	if err := brewPerl.Run(); err != nil {
		checkError(err)
	}
	if err := brewRuby.Run(); err != nil {
		checkError(err)
	}
	if err := brewPython.Run(); err != nil {
		checkError(err)
	}
	if err := fixPython.Run(); err != nil {
		checkError(err)
	}
	if err := brewLua.Run(); err != nil {
		checkError(err)
	}
	if err := brewGo.Run(); err != nil {
		checkError(err)
	}
	if err := brewRust.Run(); err != nil {
		checkError(err)
	}
	if err := brewNode.Run(); err != nil {
		checkError(err)
	}
	if err := brewTS.Run(); err != nil {
		checkError(err)
	}
	if err := brewPHP.Run(); err != nil {
		checkError(err)
	}
	if err := brewJDK.Run(); err != nil {
		checkError(err)
	}
	if err := brewGroovy.Run(); err != nil {
		checkError(err)
	}
	if err := brewKotlin.Run(); err != nil {
		checkError(err)
	}
	if err := brewScala.Run(); err != nil {
		checkError(err)
	}
	if err := brewMaven.Run(); err != nil {
		checkError(err)
	}
	if err := brewGradle.Run(); err != nil {
		checkError(err)
	}
	if err := brewClojure.Run(); err != nil {
		checkError(err)
	}
	if err := brewErlang.Run(); err != nil {
		checkError(err)
	}
	if err := brewElixir.Run(); err != nil {
		checkError(err)
	}

	shrcAppend := "# JAVA\n" +
		"export PATH=\"" + brewPrefix + "opt/openjdk/bin:$PATH\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/openjdk/include\"\n\n" +
		"# RUBY\n" +
		"export PATH=\"" + brewPrefix + "opt/ruby/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/ruby/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/ruby/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/ruby/lib/pkgconfig\"\n\n" +
		"# PYTHON\n" +
		"# brew link --overwrite python@[version]\n\n" +
		"# ANDROID STUDIO\n" +
		"export ANDROID_HOME=$HOME/Library/Android/sdk\n" +
		"export PATH=$PATH:$ANDROID_HOME/emulator\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools/bin\n" +
		"export PATH=$PATH:$ANDROID_HOME/platform-tools\n\n"
	appendFile(shrcPath, shrcAppend)
	ldBar.Stop()
}

func macUtility() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing advanced utilities for terminal..."
	ldBar.FinalMSG = " - Installed advanced utilities!\n"
	ldBar.Start()

	brewTmux := exec.Command(cmdPMS, cmdIn, "tmux")
	brewTmuxinator := exec.Command(cmdPMS, cmdIn, "tmuxinator")
	brewFzf := exec.Command(cmdPMS, cmdIn, "fzf")
	brewNeofetch := exec.Command(cmdPMS, cmdIn, "neofetch")
	brewAsciinema := exec.Command(cmdPMS, cmdIn, "asciinema")
	if err := brewTmux.Run(); err != nil {
		checkError(err)
	}
	if err := brewTmuxinator.Run(); err != nil {
		checkError(err)
	}
	if err := brewFzf.Run(); err != nil {
		checkError(err)
	}
	if err := brewNeofetch.Run(); err != nil {
		checkError(err)
	}
	if err := brewAsciinema.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func macEnd() {
	shrcAppend := "\n######## ADD CUSTOM VALUES UNDER HERE ########\n\n\n"
	appendFile(shrcPath, shrcAppend)
}

func main() {
	fmt.Println("\nDev4mac v" + appVer + "\n")
	if checkNetStatus() == true {
		macBegin()
		macEnv()
		macGit()
		macTerminal()
		macDependency()
		macDevToolCLI()
		macASDF()
		macServer()
		macLanguage()
		macUtility()
		macEnd()
		fmt.Println("\nFinished to setup! You can choose 4 options. (Recommend option is 1)\n" +
			"\t1. Setup zsh theme & Configure git global\n" +
			"\t2. Only setup zsh theme that minimal type\n" +
			"\t3. Only configure git global easily\n" +
			"\t0. Nothing, finish Dev4mac (manual setup)\n")
	endOpt:
		for {
			fmt.Printf("Select command: ")
			_, err := fmt.Scanln(&cmdOpt)
			checkError(err)
			if cmdOpt == "1" {
				confZshTheme()
				confG4s()
			} else if cmdOpt == "2" {
				confZshTheme()
			} else if cmdOpt == "3" {
				confG4s()
			} else if cmdOpt == "0" || cmdOpt == "q" || cmdOpt == "e" || cmdOpt == "quit" || cmdOpt == "exit" {
			} else {
				fmt.Println("Wrong answer. Please choose number 0-3")
				goto endOpt
			}
			break
		}
		fmt.Println("\n----------Finished!----------\n" +
			"Please RESTART your terminal!\n" +
			lstDot + "Enter this on terminal: source ~/.zprofile && source ~/.zshrc\n" +
			lstDot + "Or restart the Terminal.app by yourself.\n")
	} else {
		fmt.Println(lstDot + "Please check your internet connection and try again.\n")
	}
}
