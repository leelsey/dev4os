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
	"runtime"
	"time"
)

var (
	appVer     = "0.1"
	Git4setV   = "Git4set-0.1"
	lstDot     = " • "
	dlDir      = homeDir() + "Downloads/"
	cmdPMS     = "brew"
	cmdIn      = "install"
	cmdRein    = "reinstall"
	cmdRm      = "remove"
	cmdEcho    = "echo"
	cmdASDF    = "asdf"
	asdfPlugin = "plugin"
	asdfAdd    = "add"
	asdfReshim = "reshim"
	zshrcPath  = homeDir() + ".zshrc"
	prefixPath = brewPrefix()
	cmdOpt     string
	userName   string
	userEmail  string
)

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return err != nil
}

func homeDir() string {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDirPath + "/"
}

func workingDir() string {
	workingDirPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return workingDirPath + "/"
}

func brewPrefix() string {
	switch runtime.GOARCH {
	case "arm64":
		return "/opt/homebrew/"
	}
	return "/usr/local/"
}

func checkOnline() bool {
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

func openZSHRC(zshrcAppend string) {
	zshrcFile, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, os.FileMode(0600))
	checkError(err)
	defer zshrcFile.Close()
	_, err = zshrcFile.Write([]byte(zshrcAppend))
	checkError(err)
}

func confAlias4sh() {
	err := os.MkdirAll(homeDir()+".config/alias4sh", 0755)
	checkError(err)
	alias4shFile, err := os.OpenFile(homeDir()+".config/alias4sh/aliasrc", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	checkError(err)
	defer alias4shFile.Close()
	aliasrcContents := "#             _ _           _  _       _ \n#       /\\   | (_)         | || |     | | \n#      /  \\  | |_  __ _ ___| || |_ ___| |__ \n#     / /\\ \\ | | |/ _` / __|__   _/ __| '_ \\ \n#    / ____ \\| | | (_| \\__ \\  | | \\__ \\ | | | \n#   /_/    \\_\\_|_|\\__,_|___/  |_| |___/_| |_| \n#\n\nalias shrl=\"exec $SHELL\"\nalias zshrl=\"source ~/.zshrc\"\nalias his=\"history\"\nalias hisp=\"history -p\"\nalias hisc=\"echo -n > ~/.zsh_history && history -p  && exec $SHELL -l\"\nalias hiscl=\"rm -f ~/.bash_history && rm -f ~/.node_repl_history && rm -f ~/.python_history\"\nalias grep=\"grep --color=auto\"\nalias egrep=\"egrep --color=auto\"\nalias fgrep=\"fgrep --color=auto\"\nalias diff=\"diff --color=auto\"\nalias ls=\"ls --color=auto\"\nalias l=\"ls -CF\"\nalias ll=\"ls -l\"\nalias la=\"ls -A\"\nalias lla=\"ls -al\"\nalias lld=\"ls -al --group-directories-first\"\nalias lst=\"ls -al | grep -v '^[d|b|c|l|p|s|-]'\"\nalias lr=\"ls -lR\"\nalias tree=\"tree -Csu\"\nalias dir=\"dir --color=auto\"\nalias dird=\"dir -al --group-directories-first\"\nalias vdir=\"vdir --color=auto\"\nalias cls=\"clear\"\nalias ip=\"ipconfig\"\nalias dfh=\"df -h\"\nalias duh=\"du -h\"\nalias cdh=\"cd ~\"\nalias p=\"cd ..\"\nalias f=\"finger\"\nalias j=\"jobs -l\"\nalias d=\"date\"\nalias c=\"cal\"\n#alias curl=\"curl -w '\\n'\"\n#alias rm=\"rm -i\"\n#alias cp=\"cp -i\"\n#alias mv=\"mv -i\"\n#alias mkdir=\"mkdir -p\"\n#alias rmdir=\"rmdir -p\"\n"
	_, err = alias4shFile.Write([]byte(aliasrcContents))
	checkError(err)
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

	fmt.Println(" " + lstDot + "Make \"gitignore_global\" file in " + homeDir() + ".config/git")
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

func installBrew(whichBrew string) {
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
		checkError(err)
	}

	if _, err := os.Stat(dlBrewPath); !os.IsNotExist(err) {
		err := os.Remove(dlBrewPath)
		checkError(err)
	}

	if checkBrew() == true {
		envBrewPrefix := exec.Command("eval", whichBrew)
		envBrewPrefix.Env = append(os.Environ())
		envBrewPrefix.Run()
		updateBrew()
	} else {
		fmt.Println(lstDot + "Brew install failed, please check your system\n")
		os.Exit(0)
	}
}

func checkBrew() bool {
	if _, err := os.Stat("/opt/homebrew/bin/brew"); !os.IsNotExist(err) {
		return true
	} else if _, err := os.Stat("/usr/local/bin/brew"); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func updateBrew() {
	updateHomebrew := exec.Command(cmdPMS, "update")
	updateBrewCask := exec.Command(cmdPMS, "tap", "homebrew/cask-versions")

	updateHomebrew.Run()
	updateBrewCask.Run()
}

func macBrew() {
	if checkBrew() == true {
		ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
		ldBar.Suffix = " Updating homebrew..."
		ldBar.FinalMSG = " - Updated brew!\n"
		ldBar.Start()

		updateBrew()
		ldBar.Stop()
	} else {
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
			ldBar.FinalMSG = " - Installed brew!\n"
			ldBar.Start()

			switch runtime.GOARCH {
			case "arm64":
				whichBrew := "\"$(/opt/homebrew/bin/brew shellenv)\""
				installBrew(whichBrew)
			case "amd64":
				whichBrew := "\"$(/usr/local/bin/brew shellenv)\""
				installBrew(whichBrew)
			default:
				fmt.Println(lstDot + "Sorry, your architecture is not supported\n")
				os.Exit(0)
			}
			ldBar.Stop()
		} else {
			fmt.Println(lstDot + "Incorrect user, please check permission of sudo.\n" +
				lstDot + "It need sudo command of \"root\" user's permission.\n" +
				lstDot + "Now your username: " + string(whoAmI))
			os.Exit(0)
		}
	}
}

func macGit() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing git..."
	ldBar.FinalMSG = " - Installed git!\n"
	ldBar.Start()

	installGit := exec.Command(cmdPMS, cmdIn, "git")
	installGitLfs := exec.Command(cmdPMS, cmdIn, "git-lfs")
	gitLfsInstall := exec.Command("git", "lfs", "install")
	gitBranchMain := exec.Command("git", "config", "--global", "init.defaultBranch", "main")

	installGit.Run()
	installGitLfs.Run()
	gitLfsInstall.Run()
	gitBranchMain.Run()
	ldBar.Stop()
}

func macTerminal() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing zsh with useful tools..."
	ldBar.FinalMSG = " - Installed useful tools for terminal!\n"
	ldBar.Start()

	installNCurses := exec.Command(cmdPMS, cmdIn, "ncurses")
	installSSL := exec.Command(cmdPMS, cmdIn, "openssl")
	installZsh := exec.Command(cmdPMS, cmdIn, "zsh")
	installZshSyntax := exec.Command(cmdPMS, cmdIn, "zsh-syntax-highlighting")
	installZshAuto := exec.Command(cmdPMS, cmdIn, "zsh-autosuggestions")
	installZshComp := exec.Command(cmdPMS, cmdIn, "zsh-completions")
	installTree := exec.Command(cmdPMS, cmdIn, "tree")
	installZshTheme := exec.Command(cmdPMS, cmdIn, "romkatv/powerlevel10k/powerlevel10k")

	installNCurses.Run()
	installSSL.Run()
	installZsh.Run()
	installZshSyntax.Run()
	installZshAuto.Run()
	installZshComp.Run()
	installTree.Run()
	installZshTheme.Run()
	confAlias4sh()

	zshrcFile, err := os.OpenFile(zshrcPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	checkError(err)
	defer zshrcFile.Close()

	zshrcInitial := "# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.\n" +
		"[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh\n\n" +
		"# POWERLEVEL10K\n" +
		"source " + prefixPath + "/opt/powerlevel10k/powerlevel10k.zsh-theme\n\n" +
		"#   _________  _   _ ____   ____    __  __    _    ___ _   _\n" +
		"#  |__  / ___|| | | |  _ \\ / ___|  |  \\/  |  / \\  |_ _| \\ | |\n" +
		"#    / /\\___ \\| |_| | |_) | |      | |\\/| | / _ \\  | ||  \\| |\n" +
		"#   / /_ ___) |  _  |  _ <| |___   | |  | |/ ___ \\ | || |\\  |\n" +
		"#  /____|____/|_| |_|_| \\_\\\\____|  |_|  |_/_/   \\_\\___|_| \\_|\n#\n\n" +
		"# ZSH\n" +
		"export SHELL=zsh\n\n" +
		"# Alias4sh\n" +
		"source ~/.config/alias4sh/aliasrc\n\n" +
		"# ZSH SYNTAX HIGHTLIGHTING\n" +
		"source " + prefixPath + "share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh\n\n" +
		"# ZSH AUTOSUGGESTIONS\n" +
		"source " + prefixPath + "share/zsh-autosuggestions/zsh-autosuggestions.zsh\n\n" +
		"# NCURSES\n" +
		"export PATH=\"" + prefixPath + "opt/ncurses/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/ncurses/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/ncurses/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/ncurses/lib/pkgconfig\"\n\n" +
		"" +
		"# OPENSSL\n" +
		"export PATH=\"" + prefixPath + "opt/openssl@3/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/openssl@3/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/openssl@3/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/openssl@3/lib/pkgconfig\"\n\n"
	_, err = zshrcFile.Write([]byte(zshrcInitial))
	checkError(err)
	ldBar.Stop()
}

func macDependency() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for development work..."
	ldBar.FinalMSG = " - Installed dependencies!\n"
	ldBar.Start()

	installKRB5 := exec.Command(cmdPMS, cmdIn, "krb5")
	installGnuPG := exec.Command(cmdPMS, cmdIn, "gnupg")
	installcURL := exec.Command(cmdPMS, cmdIn, "curl")
	installWget := exec.Command(cmdPMS, cmdIn, "wget")
	installXZ := exec.Command(cmdPMS, cmdIn, "xz")
	installGzip := exec.Command(cmdPMS, cmdIn, "gzip")
	installLibzip := exec.Command(cmdPMS, cmdIn, "libzip")
	installBzip2 := exec.Command(cmdPMS, cmdIn, "bzip2")
	installZLib := exec.Command(cmdPMS, cmdIn, "zlib")
	installPkgConfig := exec.Command(cmdPMS, cmdIn, "pkg-config")
	installReadLine := exec.Command(cmdPMS, cmdIn, "readline")
	installLibffi := exec.Command(cmdPMS, cmdIn, "libffi")
	installGuile := exec.Command(cmdPMS, cmdIn, "guile")
	installGnuGetOpt := exec.Command(cmdPMS, cmdIn, "gnu-getopt")
	installCoreUtils := exec.Command(cmdPMS, cmdIn, "coreutils")
	installBison := exec.Command(cmdPMS, cmdIn, "bison")
	installLibIconv := exec.Command(cmdPMS, cmdIn, "libiconv")
	installICU4C := exec.Command(cmdPMS, cmdIn, "icu4c")
	installRe2C := exec.Command(cmdPMS, cmdIn, "re2c")
	installGD := exec.Command(cmdPMS, cmdIn, "gd")
	installCaCert := exec.Command(cmdPMS, cmdIn, "ca-certificates")
	installLDNS := exec.Command(cmdPMS, cmdIn, "ldns")
	installHTMLXMLUtils := exec.Command(cmdPMS, cmdIn, "html-xml-utils")
	installXMLto := exec.Command(cmdPMS, cmdIn, "xmlto")
	installGMP := exec.Command(cmdPMS, cmdIn, "gmp")
	installLibSodium := exec.Command(cmdPMS, cmdIn, "libsodium")
	installImageMagick := exec.Command(cmdPMS, cmdIn, "imagemagick")
	installGhostscript := exec.Command(cmdPMS, cmdIn, "ghostscript")

	installKRB5.Run()
	installGnuPG.Run()
	installcURL.Run()
	installWget.Run()
	installXZ.Run()
	installGzip.Run()
	installLibzip.Run()
	installBzip2.Run()
	installZLib.Run()
	installPkgConfig.Run()
	installReadLine.Run()
	installLibffi.Run()
	installGuile.Run()
	installGnuGetOpt.Run()
	installCoreUtils.Run()
	installBison.Run()
	installLibIconv.Run()
	installICU4C.Run()
	installRe2C.Run()
	installGD.Run()
	installCaCert.Run()
	installLDNS.Run()
	installHTMLXMLUtils.Run()
	installXMLto.Run()
	installGMP.Run()
	installLibSodium.Run()
	installImageMagick.Run()
	installGhostscript.Run()
	ldBar.Stop()

	zshrcAppend := "# KRB5\n" +
		"export PATH=\"" + prefixPath + "opt/krb5/bin:$PATH\"\n" +
		"export PATH=\"" + prefixPath + "opt/krb5/sbin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/krb5/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/krb5/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/krb5/lib/pkgconfig\"\n\n" +
		"# CURL\n" +
		"export PATH=\"" + prefixPath + "opt/curl/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/curl/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/curl/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/curl/lib/pkgconfig\"\n\n" +
		"# BZIP2\n" +
		"export PATH=\"" + prefixPath + "opt/bzip2/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/bzip2/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/bzip2/include\"\n\n" +
		"# ZLIB\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/zlib/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/zlib/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/zlib/lib/pkgconfig\"\n\n" +
		"# GNU GETOPT\n" +
		"export PATH=\"" + prefixPath + "opt/gnu-getopt/bin:$PATH\"\n\n" +
		"# COREUTILS\n" +
		"export PATH=\"" + prefixPath + "opt/coreutils/libexec/gnubin:$PATH\"\n\n" +
		"# BISON\n" +
		"export PATH=\"" + prefixPath + "opt/bison/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/bison/lib\"\n\n" +
		"# LIBICONV\n" +
		"export PATH=\"" + prefixPath + "opt/libiconv/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/libiconv/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/libiconv/include\"\n\n" +
		"# ICU4C\n" +
		"export PATH=\"" + prefixPath + "opt/icu4c/bin:$PATH\"\n" +
		"export PATH=\"" + prefixPath + "opt/icu4c/sbin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/icu4c/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/icu4c/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/icu4c/lib/pkgconfig\"\n\n"
	openZSHRC(zshrcAppend)
	ldBar.Stop()
}

func macDevToolCLI() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for CLI"
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	installSSH := exec.Command(cmdPMS, cmdIn, "openssh")
	installGawk := exec.Command(cmdPMS, cmdIn, "gawk")
	installTig := exec.Command(cmdPMS, cmdIn, "tig")
	installJQ := exec.Command(cmdPMS, cmdIn, "jq")
	installDirEnv := exec.Command(cmdPMS, cmdIn, "direnv")
	installWatchman := exec.Command(cmdPMS, cmdIn, "watchman")
	installQEMU := exec.Command(cmdPMS, cmdIn, "qemu")
	installCCache := exec.Command(cmdPMS, cmdIn, "ccache")
	installMake := exec.Command(cmdPMS, cmdIn, "make")
	installVim := exec.Command(cmdPMS, cmdIn, "vim")
	installBat := exec.Command(cmdPMS, cmdIn, "bat")
	installGH := exec.Command(cmdPMS, cmdIn, "gh")

	installSSH.Run()
	installGawk.Run()
	installTig.Run()
	installJQ.Run()
	installDirEnv.Run()
	installWatchman.Run()
	installQEMU.Run()
	installCCache.Run()
	installMake.Run()
	installVim.Run()
	installBat.Run()
	installGH.Run()

	zshrcAppend := "# DIRENV\n" +
		"eval \"$(direnv hook zsh)\"\n\n"
	openZSHRC(zshrcAppend)
	ldBar.Stop()
}

func macASDF() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing ASDF-VM with plugin..."
	ldBar.FinalMSG = " - Installed ASDF-VM, and add basic languages!\n"
	ldBar.Start()

	installASDF := exec.Command(cmdPMS, cmdIn, cmdASDF)
	installASDF.Run()

	zshrcAppend := "# ASDF VM\n" +
		"source " + prefixPath + "opt/asdf/libexec/asdf.sh\n\n"
	openZSHRC(zshrcAppend)

	if _, err := os.Stat(homeDir() + ".asdf/plugins/perl"); errors.Is(err, os.ErrNotExist) {
		addASDFPerl := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "perl")
		addASDFPerl.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/ruby"); errors.Is(err, os.ErrNotExist) {
		addASDFRuby := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "ruby")
		addASDFRuby.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/python"); errors.Is(err, os.ErrNotExist) {
		addASDFPython := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "python")
		addASDFPython.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/lua"); errors.Is(err, os.ErrNotExist) {
		addASDFLua := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "lua")
		addASDFLua.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/golang"); errors.Is(err, os.ErrNotExist) {
		addASDFGo := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "golang")
		addingASDFGo, err := addASDFGo.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFGo))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/rust"); errors.Is(err, os.ErrNotExist) {
		addASDFRust := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "rust")
		addASDFRust.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/nodejs"); errors.Is(err, os.ErrNotExist) {
		addASDFNode := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "nodejs")
		addASDFNode.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/php"); errors.Is(err, os.ErrNotExist) {
		addASDFPHP := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "php")
		addASDFPHP.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/java"); errors.Is(err, os.ErrNotExist) {
		addASDFJava := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "java")
		addASDFJava.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/groovy"); errors.Is(err, os.ErrNotExist) {
		addASDFGroovy := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "groovy")
		addASDFGroovy.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/kotlin"); errors.Is(err, os.ErrNotExist) {
		addASDFKotlin := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "kotlin")
		addASDFKotlin.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/scala"); errors.Is(err, os.ErrNotExist) {
		addASDFScala := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "scala")
		addASDFScala.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/clojure"); errors.Is(err, os.ErrNotExist) {
		addASDFClojure := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "clojure")
		addASDFClojure.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/erlang"); errors.Is(err, os.ErrNotExist) {
		addASDFErlang := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "erlang")
		addASDFErlang.Run()
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/elixir"); errors.Is(err, os.ErrNotExist) {
		addASDFElixir := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "elixir")
		addASDFElixir.Run()
	}
	ldBar.Stop()
}

func macServer() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server..."
	ldBar.FinalMSG = " - Installed server and database!\n"
	ldBar.Start()

	installHTTPD := exec.Command(cmdPMS, cmdIn, "httpd")
	installTomcat := exec.Command(cmdPMS, cmdIn, "tomcat")
	installSQLite := exec.Command(cmdPMS, cmdIn, "sqlite")
	installPostgreSQL := exec.Command(cmdPMS, cmdIn, "postgresql")
	installMySQL := exec.Command(cmdPMS, cmdIn, "mysql")
	installRedis := exec.Command(cmdPMS, cmdIn, "redis")

	installHTTPD.Run()
	installTomcat.Run()
	installSQLite.Run()
	installPostgreSQL.Run()
	installMySQL.Run()
	installRedis.Run()

	zshrcAppend := "# SQLITE3\n" +
		"export PATH=\"" + prefixPath + "opt/sqlite/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/sqlite/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/sqlite/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/sqlite/lib/pkgconfig\"\n\n"
	openZSHRC(zshrcAppend)
	ldBar.Stop()
}

func macLanguage() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing computer programming language..."
	ldBar.FinalMSG = " - Installed basic languages!\n"
	ldBar.Start()

	installPerl := exec.Command(cmdPMS, cmdIn, "perl")
	installRuby := exec.Command(cmdPMS, cmdIn, "ruby")
	installPython := exec.Command(cmdPMS, cmdIn, "python@3.10")
	fixPython := exec.Command(cmdPMS, "link", "--overwrite", "python@3.10")
	installLua := exec.Command(cmdPMS, cmdIn, "lua")
	installGo := exec.Command(cmdPMS, cmdIn, "go")
	installRust := exec.Command(cmdPMS, cmdIn, "rust")
	installNode := exec.Command(cmdPMS, cmdIn, "node")
	installTS := exec.Command(cmdPMS, cmdIn, "typescript")
	installPHP := exec.Command(cmdPMS, cmdIn, "php")
	installJDK := exec.Command(cmdPMS, cmdIn, "openjdk")
	installGroovy := exec.Command(cmdPMS, cmdIn, "groovy")
	installKotlin := exec.Command(cmdPMS, cmdIn, "kotlin")
	installScala := exec.Command(cmdPMS, cmdIn, "scala")
	installMaven := exec.Command(cmdPMS, cmdIn, "maven")
	installGradle := exec.Command(cmdPMS, cmdIn, "gradle")
	installClojure := exec.Command(cmdPMS, cmdIn, "clojure")
	installErlang := exec.Command(cmdPMS, cmdIn, "erlang")
	installElixir := exec.Command(cmdPMS, cmdIn, "elixir")

	installPerl.Run()
	installRuby.Run()
	installPython.Run()
	fixPython.Run()
	installLua.Run()
	installGo.Run()
	installRust.Run()
	installNode.Run()
	installTS.Run()
	installPHP.Run()
	installJDK.Run()
	installGroovy.Run()
	installKotlin.Run()
	installScala.Run()
	installMaven.Run()
	installGradle.Run()
	installClojure.Run()
	installErlang.Run()
	installElixir.Run()

	zshrcAppend := "# JAVA\n" +
		"export PATH=\"" + prefixPath + "opt/openjdk/bin:$PATH\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/openjdk/include\"\n\n" +
		"# RUBY\n" +
		"export PATH=\"" + prefixPath + "opt/ruby/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + prefixPath + "opt/ruby/lib\"\n" +
		"export CPPFLAGS=\"" + prefixPath + "opt/ruby/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + prefixPath + "opt/ruby/lib/pkgconfig\"\n\n" +
		"# PYTHON\n" +
		"# brew link --overwrite python@[version]\n\n" +
		"# ANDROID STUDIO\n" +
		"export ANDROID_HOME=$HOME/Library/Android/sdk\n" +
		"export PATH=$PATH:$ANDROID_HOME/emulator\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools/bin\n" +
		"export PATH=$PATH:$ANDROID_HOME/platform-tools\n\n"
	openZSHRC(zshrcAppend)
	ldBar.Stop()
}

func macUtility() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing advanced utilities for terminal..."
	ldBar.FinalMSG = " - Installed advanced utilities!\n"
	ldBar.Start()

	installTmux := exec.Command(cmdPMS, cmdIn, "tmux")
	installTmuxinator := exec.Command(cmdPMS, cmdIn, "tmuxinator")
	installFzf := exec.Command(cmdPMS, cmdIn, "fzf")
	installNeofetch := exec.Command(cmdPMS, cmdIn, "neofetch")
	installAsciinema := exec.Command(cmdPMS, cmdIn, "asciinema")

	installTmux.Run()
	installTmuxinator.Run()
	installFzf.Run()
	installNeofetch.Run()
	installAsciinema.Run()

	ldBar.Stop()
}

func macEnd() {
	zshrcAppend := "\n######## ADD CUSTOM VALUES UNDER HERE ########\n\n\n"
	openZSHRC(zshrcAppend)
	fmt.Println("\n----------Finished!----------\n" +
		"Please RESTART your terminal!\n" +
		lstDot + "Enter this on terminal: source ~/.zshrc\n" +
		lstDot + "Or restart the Terminal.app by yourself.\n")
}

func main() {
	fmt.Println("\nDev4mac v" + appVer + "\n")
	if checkOnline() == true {
		macBrew()
		macGit()
		macTerminal()
		macDependency()
		macDevToolCLI()
		macASDF()
		macServer()
		macLanguage()
		macUtility()
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
			macEnd()
			break
		}
	} else {
		fmt.Println(lstDot + "Please check your internet connection and try again.\n")
	}
}
