package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"log"
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
	brewPrefix  = cehckBrewPrefix()
	cmdPMS      = checkBrewPath()
	cmdIn       = "install"
	cmdReIn     = "reinstall"
	cmdRm       = "remove"
	cmdASDF     = checkASDFPath()
	asdfPlugin  = "plugin"
	asdfAdd     = "add"
	asdfShim    = "reshim"
	cmdOpt      string
)

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return err != nil
}

func checkNetStatus() bool {
	getTimeout := time.Duration(10000 * time.Millisecond)
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

func cehckBrewPrefix() string {
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
	//if _, err := os.Stat("/opt/homebrew/bin/brew"); !os.IsNotExist(err) {
	if _, err := os.Stat(cmdPMS); !os.IsNotExist(err) {
		return true
		//} else if _, err := os.Stat("/usr/local/bin/brew"); !os.IsNotExist(err) {
		//	return true
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
	user, err := user.Current()
	checkError(err)
	return user.Username
}

func generateFile(filePath, fileContents string) {
	profileFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	checkError(err)
	defer profileFile.Close()
	_, err = profileFile.Write([]byte(fileContents))
	checkError(err)
}

func appendFile(filePath, fileContents string) {
	zshrcFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.FileMode(0600))
	checkError(err)
	defer zshrcFile.Close()
	_, err = zshrcFile.Write([]byte(fileContents))
	checkError(err)
}

func newZProfile() {
	//profileFile, err := os.OpenFile(profilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	//checkError(err)
	//defer profileFile.Close()
	profileInitial := "# " + currentUser() + "’s profile\n\n" +
		"# ZSH\n" +
		"export SHELL=zsh\n"
	generateFile(profilePath, profileInitial)
	//_, err = profileFile.Write([]byte(profileInitial))
	//checkError(err)
}

func newZshRC() {
	//zshrcFile, err := os.OpenFile(shrcPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	//checkError(err)
	//defer zshrcFile.Close()
	zshrcInitial := "#   _________  _   _ ____   ____    __  __    _    ___ _   _\n" +
		"#  |__  / ___|| | | |  _ \\ / ___|  |  \\/  |  / \\  |_ _| \\ | |\n" +
		"#    / /\\___ \\| |_| | |_) | |      | |\\/| | / _ \\  | ||  \\| |\n" +
		"#   / /_ ___) |  _  |  _ <| |___   | |  | |/ ___ \\ | || |\\  |\n" +
		"#  /____|____/|_| |_|_| \\_\\\\____|  |_|  |_/_/   \\_\\___|_| \\_|\n#\n\n"
	generateFile(shrcPath, zshrcInitial)
	//_, err = zshrcFile.Write([]byte(zshrcInitial))
	//checkError(err)
}

func rmFile(filePath string) {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err := os.Remove(filePath)
		checkError(err)
	}
}

func confA4s() {
	//err := os.MkdirAll(homeDir()+".config/alias4sh", 0755)
	//checkError(err)
	//alias4shFile, err := os.OpenFile(homeDir()+".config/alias4sh/aliasrc", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	//checkError(err)
	//defer alias4shFile.Close()
	//aliasrcContents := "#             _ _           _  _       _ \n#       /\\   | (_)         | || |     | | \n#      /  \\  | |_  __ _ ___| || |_ ___| |__ \n#     / /\\ \\ | | |/ _` / __|__   _/ __| '_ \\ \n#    / ____ \\| | | (_| \\__ \\  | | \\__ \\ | | | \n#   /_/    \\_\\_|_|\\__,_|___/  |_| |___/_| |_| \n#\n\nalias shrl=\"exec $SHELL\"\nalias zshrl=\"source ~/.zprofile ~/.zshrc\"\nalias his=\"history\"\nalias hisp=\"history -p\"\nalias hisc=\"echo -n > ~/.zsh_history && history -p  && exec $SHELL -l\"\nalias hiscl=\"rm -f ~/.bash_history && rm -f ~/.node_repl_history && rm -f ~/.python_history\"\nalias grep=\"grep --color=auto\"\nalias egrep=\"egrep --color=auto\"\nalias fgrep=\"fgrep --color=auto\"\nalias diff=\"diff --color=auto\"\nalias ls=\"ls --color=auto\"\nalias l=\"ls -CF\"\nalias ll=\"ls -l\"\nalias la=\"ls -A\"\nalias lla=\"ls -al\"\nalias lld=\"ls -al --group-directories-first\"\nalias lst=\"ls -al | grep -v '^[d|b|c|l|p|s|-]'\"\nalias lr=\"ls -lR\"\nalias tree=\"tree -Csu\"\nalias dir=\"dir --color=auto\"\nalias dird=\"dir -al --group-directories-first\"\nalias vdir=\"vdir --color=auto\"\nalias cls=\"clear\"\nalias ip=\"ipconfig\"\nalias dfh=\"df -h\"\nalias duh=\"du -h\"\nalias cdh=\"cd ~\"\nalias p=\"cd ..\"\nalias f=\"finger\"\nalias j=\"jobs -l\"\nalias d=\"date\"\nalias c=\"cal\"\n#alias curl=\"curl -w '\\n'\"\n#alias rm=\"rm -i\"\n#alias cp=\"cp -i\"\n#alias mv=\"mv -i\"\n#alias mkdir=\"mkdir -p\"\n#alias rmdir=\"rmdir -p\""
	//_, err = alias4shFile.Write([]byte(aliasrcContents))
	//checkError(err)

	dlA4sPath := workingDir() + ".dev4mac-alias4sh.sh"
	resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Alias4sh/main/install.sh")
	if err != nil {
		fmt.Println(lstDot + "Brew install URL is maybe changed, please check https://github.com/leelsey/Alias4sh\n")
		os.Exit(0)
	}
	defer resp.Body.Close()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	a4sInstaller, err := os.OpenFile(dlA4sPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0755))
	checkError(err)
	defer a4sInstaller.Close()
	_, err = a4sInstaller.Write([]byte(rawFile))
	checkError(err)

	installA4s := exec.Command("sh", dlA4sPath)
	if err := installA4s.Run(); err != nil {
		rmFile(dlA4sPath)
		checkError(err)
	}
	rmFile(dlA4sPath)
}

func confG4s() {
	fmt.Println("\nGit global configuration")

	fmt.Println(" 1) Main branch default name changed master -> main")
	setBranchMain := exec.Command("git", "config", "--global", "init.defaultBranch", "main")
	setBranchMain.Run()

	fmt.Println(" 2) Add your information to the global git config")
	consoleReader := bufio.NewScanner(os.Stdin)
	fmt.Printf(" " + lstDot + "User name: ")
	consoleReader.Scan()
	userName := consoleReader.Text()
	fmt.Printf(" " + lstDot + "User email: ")
	consoleReader.Scan()
	userEmail := consoleReader.Text()

	unsetUserName := exec.Command("git", "config", "--unset", "--global", "user.name")
	unsetUserEmail := exec.Command("git", "config", "--unset", "--global", "user.email")
	setUserName := exec.Command("git", "config", "--global", "user.name", userName)
	setUserEmail := exec.Command("git", "config", "--global", "user.email", userEmail)
	unsetUserName.Run()
	unsetUserEmail.Run()
	setUserName.Run()
	setUserEmail.Run()

	fmt.Println(" 3) Setup git global ignore file with directories")
	ignoreDir := homeDir() + ".config/git/"
	if err := os.MkdirAll(ignoreDir, 0755); err != nil {
		log.Fatal(err)
	}

	ignorePath := ignoreDir + "gitignore_global"
	resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Git4set/main/gitignore-sample")
	if err != nil {
		fmt.Println(lstDot + "Git Ignore sample URL is maybe changed, please check https://github.com/leelsey/Git4set\n")
		os.Exit(0)
	}
	defer resp.Body.Close()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	gitIgnore, err := os.OpenFile(ignorePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	checkError(err)
	defer gitIgnore.Close()
	_, err = gitIgnore.Write([]byte(rawFile))
	checkError(err)

	unsetExcludesFile := exec.Command("git", "config", "--unset", "--global", "core.excludesfile")
	setExcludesFile := exec.Command("git", "config", "--unset", "core.excludesfile", ignorePath)
	unsetExcludesFile.Run()
	setExcludesFile.Run()

	fmt.Println(" " + lstDot + "Make \"gitignore_global\" file in " + ignoreDir)
}

func confZshTheme() {
	dlP10kPath := homeDir() + ".p10k.zsh"
	resp, err := http.Get("https://raw.githubusercontent.com/leelsey/Dev4os/main/cmd/dev4os/dev4p10k")
	if err != nil {
		fmt.Println(lstDot + "Dev4os's p10k file URL is maybe changed, please check https://github.com/leelsey/Dev4os\n")
		os.Exit(0)
	}
	defer resp.Body.Close()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	p10kConf, err := os.OpenFile(dlP10kPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	checkError(err)
	defer p10kConf.Close()
	_, err = p10kConf.Write([]byte(rawFile))
	checkError(err)
}

func updateBrew() {
	updateHomebrew := exec.Command(cmdPMS, "update")
	updateBrewCask := exec.Command(cmdPMS, "tap", "homebrew/cask-versions")

	updateHomebrew.Run()
	updateBrewCask.Run()
}

func installBrew() {
	dlBrewPath := workingDir() + ".dev4mac-brew.sh"
	resp, err := http.Get("https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh")
	if err != nil {
		fmt.Println(lstDot + "Brew install URL is maybe changed, please check https://github.com/Homebrew/install\n")
		os.Exit(0)
	}
	defer resp.Body.Close()
	rawFile, _ := ioutil.ReadAll(resp.Body)

	brewInstaller, err := os.OpenFile(dlBrewPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0755))
	checkError(err)
	defer brewInstaller.Close()
	_, err = brewInstaller.Write([]byte(rawFile))
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

	brewGit := exec.Command(cmdPMS, cmdIn, "git")
	brewGitLfs := exec.Command(cmdPMS, cmdIn, "git-lfs")
	brewGit.Run()
	brewGitLfs.Run()
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
	brewNCurses.Run()
	brewSSL.Run()
	brewZsh.Run()
	brewZshSyntax.Run()
	brewZshAuto.Run()
	brewZshComp.Run()
	brewTree.Run()
	brewZshTheme.Run()

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
	brewKRB5.Run()
	brewGnuPG.Run()
	brewcURL.Run()
	brewWget.Run()
	brewXZ.Run()
	brewGzip.Run()
	brewLibzip.Run()
	brewBzip2.Run()
	brewZLib.Run()
	brewPkgConfig.Run()
	brewReadLine.Run()
	brewLibffi.Run()
	brewGuile.Run()
	brewGnuGetOpt.Run()
	brewCoreUtils.Run()
	brewBison.Run()
	brewLibIconv.Run()
	brewICU4C.Run()
	brewRe2C.Run()
	brewGD.Run()
	brewCaCert.Run()
	brewLDNS.Run()
	brewHTMLXMLUtils.Run()
	brewXMLto.Run()
	brewGMP.Run()
	brewLibSodium.Run()
	brewImageMagick.Run()
	brewGhostscript.Run()
	ldBar.Stop()

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
	brewSSH.Run()
	brewGawk.Run()
	brewTig.Run()
	brewJQ.Run()
	brewDirEnv.Run()
	brewWatchman.Run()
	brewQEMU.Run()
	brewCCache.Run()
	brewMake.Run()
	brewVim.Run()
	brewBat.Run()
	brewGH.Run()

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
	brewASDF.Run()

	shrcAppend := "# ASDF VM\n" +
		"source " + brewPrefix + "/opt/asdf/libexec/asdf.sh\n\n"
	appendFile(shrcPath, shrcAppend)

	pluginPath := homeDir() + ".asdf/plugins/"
	if _, err := os.Stat(pluginPath + "perl"); errors.Is(err, os.ErrNotExist) {
		addASDFPerl := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "perl")
		addASDFPerl.Run()
	}
	if _, err := os.Stat(pluginPath + "ruby"); errors.Is(err, os.ErrNotExist) {
		addASDFRuby := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "ruby")
		addASDFRuby.Run()
	}
	if _, err := os.Stat(pluginPath + "python"); errors.Is(err, os.ErrNotExist) {
		addASDFPython := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "python")
		addASDFPython.Run()
	}
	if _, err := os.Stat(pluginPath + "lua"); errors.Is(err, os.ErrNotExist) {
		addASDFLua := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "lua")
		addASDFLua.Run()
	}
	if _, err := os.Stat(pluginPath + "golang"); errors.Is(err, os.ErrNotExist) {
		addASDFGo := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "golang")
		addASDFGo.Run()
	}
	if _, err := os.Stat(pluginPath + "rust"); errors.Is(err, os.ErrNotExist) {
		addASDFRust := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "rust")
		addASDFRust.Run()
	}
	if _, err := os.Stat(pluginPath + "nodejs"); errors.Is(err, os.ErrNotExist) {
		addASDFNode := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "nodejs")
		addASDFNode.Run()
	}
	if _, err := os.Stat(pluginPath + "php"); errors.Is(err, os.ErrNotExist) {
		addASDFPHP := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "php")
		addASDFPHP.Run()
	}
	if _, err := os.Stat(pluginPath + "java"); errors.Is(err, os.ErrNotExist) {
		addASDFJava := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "java")
		addASDFJava.Run()
	}
	if _, err := os.Stat(pluginPath + "groovy"); errors.Is(err, os.ErrNotExist) {
		addASDFGroovy := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "groovy")
		addASDFGroovy.Run()
	}
	if _, err := os.Stat(pluginPath + "kotlin"); errors.Is(err, os.ErrNotExist) {
		addASDFKotlin := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "kotlin")
		addASDFKotlin.Run()
	}
	if _, err := os.Stat(pluginPath + "scala"); errors.Is(err, os.ErrNotExist) {
		addASDFScala := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "scala")
		addASDFScala.Run()
	}
	if _, err := os.Stat(pluginPath + "clojure"); errors.Is(err, os.ErrNotExist) {
		addASDFClojure := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "clojure")
		addASDFClojure.Run()
	}
	if _, err := os.Stat(pluginPath + "erlang"); errors.Is(err, os.ErrNotExist) {
		addASDFErlang := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "erlang")
		addASDFErlang.Run()
	}
	if _, err := os.Stat(pluginPath + "elixir"); errors.Is(err, os.ErrNotExist) {
		addASDFElixir := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "elixir")
		addASDFElixir.Run()
	}
	asdfReshim := exec.Command(cmdASDF, asdfShim)
	asdfReshim.Run()
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
	brewHTTPD.Run()
	brewTomcat.Run()
	brewSQLite.Run()
	brewPostgreSQL.Run()
	brewMySQL.Run()
	brewRedis.Run()

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
	brewPerl.Run()
	brewRuby.Run()
	brewPython.Run()
	fixPython.Run()
	brewLua.Run()
	brewGo.Run()
	brewRust.Run()
	brewNode.Run()
	brewTS.Run()
	brewPHP.Run()
	brewJDK.Run()
	brewGroovy.Run()
	brewKotlin.Run()
	brewScala.Run()
	brewMaven.Run()
	brewGradle.Run()
	brewClojure.Run()
	brewErlang.Run()
	brewElixir.Run()

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
	brewTmux.Run()
	brewTmuxinator.Run()
	brewFzf.Run()
	brewNeofetch.Run()
	brewAsciinema.Run()
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
			fmt.Scanln(&cmdOpt)
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
