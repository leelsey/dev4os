package main

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/schollz/progressbar/v3"
	"io"
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

func confGit4set() {
	req, _ := http.NewRequest("GET",
		"https://github.com/leelsey/Git4set/archive/refs/tags/v0.1.zip", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	file, _ := os.OpenFile(Git4setV+".zip", os.O_CREATE|os.O_WRONLY, 0755)
	defer file.Close()
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Git4set...",
	)
	io.Copy(io.MultiWriter(file, bar), resp.Body)

	dlLoc := workingDir() + Git4setV + ".zip"
	mvLoc := homeDir() + "Downloads/" + Git4setV + ".zip"
	err := os.Rename(dlLoc, mvLoc)
	checkError(err)

	fmt.Println(" - Finished to download Git4sh: " + mvLoc + " (Your download directory)\n" +
		"\nPlease extract zip file and run shell script on terminal.\n" +
		lstDot + "Configure global author & ignore: sh ./initial-git.sh\n" +
		lstDot + "Only want configure global author: sh ./git-conf.sh\n" +
		lstDot + "Only want configure global ignore: sh ./git-ignore.sh")
}

func macBrew() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Checking homebrew..."
	ldBar.Start()

	checkBrew := exec.Command("which", cmdPMS)
	checkingBrew, err := checkBrew.Output()
	checkError(err)
	ldBar.Stop()
	if string(checkingBrew) == "/opt/homebrew/bin/brew\n" || string(checkingBrew) == "/usr/local/bin/brew\n" {
		ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
		ldBar.Suffix = " Updating homebrew..."
		ldBar.FinalMSG = " - Updated brew!\n"
		ldBar.Start()

		updateHomebrew := exec.Command(cmdPMS, "update")
		updateBrewCask := exec.Command(cmdPMS, "tap", "homebrew/cask-versions")

		updatingHomebrew, err := updateHomebrew.Output()
		checkError(err)
		updatingBrewCask, err := updateBrewCask.Output()
		checkError(err)

		fmt.Sprintf(string(updatingHomebrew))
		fmt.Sprintf(string(updatingBrewCask))
		ldBar.Stop()
	} else {
		fmt.Println("You need the Homebrew first, and run Dev4mac again. Check detail on this site: https://brew.sh")
		fmt.Println(lstDot + "Enter on terminal: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
		os.Exit(0)
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

	installingGit, err := installGit.Output()
	checkError(err)
	installingGitLfs, err := installGitLfs.Output()
	checkError(err)
	gitLfsInstalling, err := gitLfsInstall.Output()
	checkError(err)
	confGitMain, err := gitBranchMain.Output()
	checkError(err)

	fmt.Sprintf(string(installingGit))
	fmt.Sprintf(string(installingGitLfs))
	fmt.Sprintf(string(gitLfsInstalling))
	fmt.Sprintf(string(confGitMain))
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

	installingNCurses, err := installNCurses.Output()
	checkError(err)
	installingSSL, err := installSSL.Output()
	checkError(err)
	installingZsh, err := installZsh.Output()
	checkError(err)
	installingZshSyntax, err := installZshSyntax.Output()
	checkError(err)
	installingZshAuto, err := installZshAuto.Output()
	checkError(err)
	installingZshComp, err := installZshComp.Output()
	checkError(err)
	installingTree, err := installTree.Output()
	checkError(err)
	installingZshTheme, err := installZshTheme.Output()
	checkError(err)

	fmt.Sprintf(string(installingNCurses))
	fmt.Sprintf(string(installingSSL))
	fmt.Sprintf(string(installingZsh))
	fmt.Sprintf(string(installingZshSyntax))
	fmt.Sprintf(string(installingZshAuto))
	fmt.Sprintf(string(installingZshComp))
	fmt.Sprintf(string(installingTree))
	fmt.Sprintf(string(installingZshTheme))
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
	ldBar.Suffix = " Installing dependencies for development work"
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
	installGB := exec.Command(cmdPMS, cmdIn, "gd")
	installHTMLXMLUtils := exec.Command(cmdPMS, cmdIn, "html-xml-utils")
	installXMLto := exec.Command(cmdPMS, cmdIn, "xmlto")
	installGMP := exec.Command(cmdPMS, cmdIn, "gmp")
	installLibSodium := exec.Command(cmdPMS, cmdIn, "libsodium")
	installImageMagick := exec.Command(cmdPMS, cmdIn, "imagemagick")
	installGhostscript := exec.Command(cmdPMS, cmdIn, "ghostscript")

	installingKRB5, err := installKRB5.Output()
	checkError(err)
	installingGnuPG, err := installGnuPG.Output()
	checkError(err)
	installingcURL, err := installcURL.Output()
	checkError(err)
	installingWget, err := installWget.Output()
	checkError(err)
	installingXZ, err := installXZ.Output()
	checkError(err)
	installingGzip, err := installGzip.Output()
	checkError(err)
	installingLibzip, err := installLibzip.Output()
	checkError(err)
	installingBzip2, err := installBzip2.Output()
	checkError(err)
	installingZLib, err := installZLib.Output()
	checkError(err)
	installingPkgConfig, err := installPkgConfig.Output()
	checkError(err)
	installingReadLine, err := installReadLine.Output()
	checkError(err)
	installingLibffi, err := installLibffi.Output()
	checkError(err)
	installingGuile, err := installGuile.Output()
	checkError(err)
	installingGnuGetOpt, err := installGnuGetOpt.Output()
	checkError(err)
	installingCoreUtils, err := installCoreUtils.Output()
	checkError(err)
	installingBison, err := installBison.Output()
	checkError(err)
	installingLibIconv, err := installLibIconv.Output()
	checkError(err)
	installingICU4C, err := installICU4C.Output()
	checkError(err)
	installingRe2C, err := installRe2C.Output()
	checkError(err)
	installingGB, err := installGB.Output()
	checkError(err)
	installingHTMLXMLUtils, err := installHTMLXMLUtils.Output()
	checkError(err)
	installingXMLto, err := installXMLto.Output()
	checkError(err)
	installingGMP, err := installGMP.Output()
	checkError(err)
	installingLibSodium, err := installLibSodium.Output()
	checkError(err)
	installingImageMagick, err := installImageMagick.Output()
	checkError(err)
	installingGhostscript, err := installGhostscript.Output()
	checkError(err)

	fmt.Sprintf(string(installingKRB5))
	fmt.Sprintf(string(installingGnuPG))
	fmt.Sprintf(string(installingcURL))
	fmt.Sprintf(string(installingWget))
	fmt.Sprintf(string(installingXZ))
	fmt.Sprintf(string(installingGzip))
	fmt.Sprintf(string(installingLibzip))
	fmt.Sprintf(string(installingBzip2))
	fmt.Sprintf(string(installingZLib))
	fmt.Sprintf(string(installingPkgConfig))
	fmt.Sprintf(string(installingReadLine))
	fmt.Sprintf(string(installingLibffi))
	fmt.Sprintf(string(installingGuile))
	fmt.Sprintf(string(installingGnuGetOpt))
	fmt.Sprintf(string(installingCoreUtils))
	fmt.Sprintf(string(installingBison))
	fmt.Sprintf(string(installingLibIconv))
	fmt.Sprintf(string(installingICU4C))
	fmt.Sprintf(string(installingRe2C))
	fmt.Sprintf(string(installingGB))
	fmt.Sprintf(string(installingHTMLXMLUtils))
	fmt.Sprintf(string(installingXMLto))
	fmt.Sprintf(string(installingGMP))
	fmt.Sprintf(string(installingLibSodium))
	fmt.Sprintf(string(installingImageMagick))
	fmt.Sprintf(string(installingGhostscript))

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

	installGawk := exec.Command(cmdPMS, cmdIn, "gawk")
	installTig := exec.Command(cmdPMS, cmdIn, "tig")
	installJQ := exec.Command(cmdPMS, cmdIn, "jq")
	installDirEnv := exec.Command(cmdPMS, cmdIn, "direnv")
	installWatchman := exec.Command(cmdPMS, cmdIn, "watchman")
	installQEMU := exec.Command(cmdPMS, cmdIn, "qemu")
	installCcache := exec.Command(cmdPMS, cmdIn, "ccache")
	installMake := exec.Command(cmdPMS, cmdIn, "make")
	installVim := exec.Command(cmdPMS, cmdIn, "vim")
	installBat := exec.Command(cmdPMS, cmdIn, "bat")
	installGH := exec.Command(cmdPMS, cmdIn, "gh")

	installingGawk, err := installGawk.Output()
	checkError(err)
	installingTig, err := installTig.Output()
	checkError(err)
	installingJQ, err := installJQ.Output()
	checkError(err)
	installingDirEnv, err := installDirEnv.Output()
	checkError(err)
	installingWatchman, err := installWatchman.Output()
	checkError(err)
	installingQEMU, err := installQEMU.Output()
	checkError(err)
	installingCcache, err := installCcache.Output()
	checkError(err)
	installingMake, err := installMake.Output()
	checkError(err)
	installingVim, err := installVim.Output()
	checkError(err)
	installingBat, err := installBat.Output()
	checkError(err)
	installingGH, err := installGH.Output()
	checkError(err)

	fmt.Sprintf(string(installingGawk))
	fmt.Sprintf(string(installingTig))
	fmt.Sprintf(string(installingJQ))
	fmt.Sprintf(string(installingDirEnv))
	fmt.Sprintf(string(installingWatchman))
	fmt.Sprintf(string(installingQEMU))
	fmt.Sprintf(string(installingCcache))
	fmt.Sprintf(string(installingMake))
	fmt.Sprintf(string(installingVim))
	fmt.Sprintf(string(installingBat))
	fmt.Sprintf(string(installingGH))

	zshrcAppend := "# DIRENV\n" +
		"eval \"$(direnv hook zsh)\"\n\n"
	openZSHRC(zshrcAppend)
	ldBar.Stop()
}

func macASDF() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing ASDF-VM with plugin "
	ldBar.FinalMSG = " - Installed ASDF-VM, and add basic languages!\n"
	ldBar.Start()

	installASDF := exec.Command(cmdPMS, cmdIn, cmdASDF)
	installingASDF, err := installASDF.Output()
	checkError(err)
	fmt.Sprintf(string(installingASDF))

	zshrcAppend := "# ASDF VM\n" +
		"source " + prefixPath + "opt/asdf/libexec/asdf.sh\n\n"
	openZSHRC(zshrcAppend)

	if _, err := os.Stat(homeDir() + ".asdf/plugins/perl"); errors.Is(err, os.ErrNotExist) {
		addASDFPerl := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "perl")
		addingASDFPerl, err := addASDFPerl.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFPerl))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/ruby"); errors.Is(err, os.ErrNotExist) {
		addASDFRuby := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "ruby")
		addingASDFRuby, err := addASDFRuby.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFRuby))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/python"); errors.Is(err, os.ErrNotExist) {
		addASDFPython := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "python")
		addingASDFPython, err := addASDFPython.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFPython))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/lua"); errors.Is(err, os.ErrNotExist) {
		addASDFLua := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "lua")
		addingASDFLua, err := addASDFLua.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFLua))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/golang"); errors.Is(err, os.ErrNotExist) {
		addASDFGo := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "golang")
		addingASDFGo, err := addASDFGo.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFGo))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/rust"); errors.Is(err, os.ErrNotExist) {
		addASDFRust := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "rust")
		addingASDFRust, err := addASDFRust.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFRust))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/nodejs"); errors.Is(err, os.ErrNotExist) {
		addASDFNode := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "nodejs")
		addingASDFNode, err := addASDFNode.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFNode))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/php"); errors.Is(err, os.ErrNotExist) {
		addASDFPHP := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "php")
		addingASDFPHP, err := addASDFPHP.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFPHP))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/java"); errors.Is(err, os.ErrNotExist) {
		addASDFJava := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "java")
		addingASDFJava, err := addASDFJava.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFJava))

	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/groovy"); errors.Is(err, os.ErrNotExist) {
		addASDFGroovy := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "groovy")
		addingASDFGroovy, err := addASDFGroovy.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFGroovy))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/kotlin"); errors.Is(err, os.ErrNotExist) {
		addASDFKotlin := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "kotlin")
		addingASDFKotlin, err := addASDFKotlin.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFKotlin))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/scala"); errors.Is(err, os.ErrNotExist) {
		addASDFScala := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "scala")
		addingASDFScala, err := addASDFScala.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFScala))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/clojure"); errors.Is(err, os.ErrNotExist) {
		addASDFClojure := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "clojure")
		addingASDFClojure, err := addASDFClojure.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFClojure))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/erlang"); errors.Is(err, os.ErrNotExist) {
		addASDFErlang := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "erlang")
		addingASDFErlang, err := addASDFErlang.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFErlang))
	}
	if _, err := os.Stat(homeDir() + ".asdf/plugins/elixir"); errors.Is(err, os.ErrNotExist) {
		addASDFElixir := exec.Command(cmdASDF, asdfPlugin, asdfAdd, "elixir")
		addingASDFElixir, err := addASDFElixir.Output()
		checkError(err)
		fmt.Sprintf(string(addingASDFElixir))
	}
	ldBar.Stop()
}

func macServer() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server"
	ldBar.FinalMSG = " - Installed server and database!\n"
	ldBar.Start()

	installHTTPD := exec.Command(cmdPMS, cmdIn, "httpd")
	installTomcat := exec.Command(cmdPMS, cmdIn, "tomcat")
	installSQLite := exec.Command(cmdPMS, cmdIn, "sqlite")
	installPostgreSQL := exec.Command(cmdPMS, cmdIn, "postgresql")
	installMySQL := exec.Command(cmdPMS, cmdIn, "mysql")
	installRedis := exec.Command(cmdPMS, cmdIn, "redis")

	installingHTTPD, err := installHTTPD.Output()
	checkError(err)
	installingTomcat, err := installTomcat.Output()
	checkError(err)
	installingSQLite, err := installSQLite.Output()
	checkError(err)
	installingPostgreSQL, err := installPostgreSQL.Output()
	checkError(err)
	installingMySQL, err := installMySQL.Output()
	checkError(err)
	installingRedis, err := installRedis.Output()
	checkError(err)

	fmt.Sprintf(string(installingHTTPD))
	fmt.Sprintf(string(installingTomcat))
	fmt.Sprintf(string(installingSQLite))
	fmt.Sprintf(string(installingPostgreSQL))
	fmt.Sprintf(string(installingMySQL))
	fmt.Sprintf(string(installingRedis))

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
	ldBar.Suffix = " Installing computer programming language"
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

	installingPerl, err := installPerl.Output()
	checkError(err)
	installingRuby, err := installRuby.Output()
	checkError(err)
	installingPython, err := installPython.Output()
	checkError(err)
	fixingPython, err := fixPython.Output()
	checkError(err)
	installingLua, err := installLua.Output()
	checkError(err)
	installingGo, err := installGo.Output()
	checkError(err)
	installingRust, err := installRust.Output()
	checkError(err)
	installingNode, err := installNode.Output()
	checkError(err)
	installingTS, err := installTS.Output()
	checkError(err)
	installingPHP, err := installPHP.Output()
	checkError(err)
	installingJDK, err := installJDK.Output()
	checkError(err)
	installingGroovy, err := installGroovy.Output()
	checkError(err)
	installingKotlin, err := installKotlin.Output()
	checkError(err)
	installingScala, err := installScala.Output()
	checkError(err)
	installingClojure, err := installClojure.Output()
	checkError(err)
	installingMaven, err := installMaven.Output()
	checkError(err)
	installingGradle, err := installGradle.Output()
	checkError(err)
	installingErlang, err := installErlang.Output()
	checkError(err)
	installingElixir, err := installElixir.Output()
	checkError(err)

	fmt.Sprintf(string(installingPerl))
	fmt.Sprintf(string(installingRuby))
	fmt.Sprintf(string(installingPython))
	fmt.Sprintf(string(fixingPython))
	fmt.Sprintf(string(installingLua))
	fmt.Sprintf(string(installingGo))
	fmt.Sprintf(string(installingRust))
	fmt.Sprintf(string(installingNode))
	fmt.Sprintf(string(installingTS))
	fmt.Sprintf(string(installingPHP))
	fmt.Sprintf(string(installingJDK))
	fmt.Sprintf(string(installingGroovy))
	fmt.Sprintf(string(installingKotlin))
	fmt.Sprintf(string(installingScala))
	fmt.Sprintf(string(installingClojure))
	fmt.Sprintf(string(installingMaven))
	fmt.Sprintf(string(installingGradle))
	fmt.Sprintf(string(installingErlang))
	fmt.Sprintf(string(installingElixir))

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

	installingTmux, err := installTmux.Output()
	checkError(err)
	installingTmuxinator, err := installTmuxinator.Output()
	checkError(err)
	installingFzf, err := installFzf.Output()
	checkError(err)
	installingNeofetch, err := installNeofetch.Output()
	checkError(err)
	installingAsciinema, err := installAsciinema.Output()
	checkError(err)

	fmt.Sprintf(string(installingTmux))
	fmt.Sprintf(string(installingTmuxinator))
	fmt.Sprintf(string(installingFzf))
	fmt.Sprintf(string(installingNeofetch))
	fmt.Sprintf(string(installingAsciinema))
	ldBar.Stop()
}

func macEnd() {
	zshrcAppend := "\n######## ADD CUSTOM VALUES UNDER HERE ########\n\n\n"
	openZSHRC(zshrcAppend)
	fmt.Println("\n----------Finished!----------\n" +
		"Please RESTART your terminal!\n" +
		lstDot + "Enter this on terminal: source ~/.zshrc\n" +
		lstDot + "Or restart the Terminal.app by yourself.")
}

func main() {
	fmt.Println("\nDev4mac v" + appVer + "\n")
	macBrew()
	macGit()
	macTerminal()
	macDependency()
	macDevToolCLI()
	macASDF()
	macServer()
	macLanguage()
	macUtility()
	fmt.Printf("\nPress any key to finish, " +
		"or press (i) if you want configure global git... ")
	var setCMD string
	fmt.Scanln(&setCMD)
	if setCMD == "i" || setCMD == "I" {
		confGit4set()
	}
	macEnd()
}
