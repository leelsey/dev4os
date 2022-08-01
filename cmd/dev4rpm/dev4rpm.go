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
	"time"
)

var (
	appVer      = "0.1"
	lstDot      = " • "
	shrcPath    = homeDir() + ".zshrc"
	profilePath = homeDir() + ".zprofile"
	superUser   = "sudo"
	cmdPMS      = "dnf"
	cmdIns      = "install"
	//cmdReIns    = "reinstall"
	//cmdRm       = "remove"
	cmdYes    = "-y"
	cmdSys    = "systemctl"
	cmdEnable = "enable"
	//cmdDisable = "disable"
	cmdStart   = "start"
	cmdGit     = "git"
	gitClone   = "clone"
	cmdASDF    = homeDir() + ".asdf/"
	asdfPlugin = "plugin"
	asdfAdd    = "add"
	asdfShim   = "reshim"
	chooseCmd  = "Select command: "
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
		fmt.Println(lstDot + "Git Ignore sample‘s URL is maybe changed, please check https://github.com/leelsey/Git4set\n")
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

func updateDNF() {
	updatePMS := exec.Command(cmdPMS, "makecache", "--refresh")
	installEpel := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "epel-release")
	installPlugins := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "dnf-plugins-core")
	updateLinux := exec.Command(superUser, cmdPMS, "update", cmdYes)

	if err := updatePMS.Run(); err != nil {
		checkError(err)
	}
	if err := installEpel.Run(); err != nil {
		checkError(err)
	}
	if err := installPlugins.Run(); err != nil {
		checkError(err)
	}
	if err := updateLinux.Run(); err != nil {
		checkError(err)
	}

}

func secureConf() {
	installFirewall := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "firewalld")
	firewallOn := exec.Command(superUser, cmdSys, cmdEnable, "firewalld")
	firewallStart := exec.Command(superUser, cmdSys, cmdStart, "firewalld")
	if err := installFirewall.Run(); err != nil {
		checkError(err)
	}
	if err := firewallOn.Run(); err != nil {
		checkError(err)
	}
	if err := firewallStart.Run(); err != nil {
		checkError(err)
	}
	fileContents := "net.ipv4.icmp_echo_ignore_all = 1\n"
	appendFile("/etc/sysctl.conf", fileContents)
	sysctlConf := exec.Command(superUser, "sysctl", "-p")
	if err := sysctlConf.Run(); err != nil {
		checkError(err)
	}
}

func linuxBegin() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Updating Linux..."
	ldBar.FinalMSG = " - Updated Linux!\n"
	ldBar.Start()

	updateDNF()
	secureConf()
	ldBar.Stop()
}

func linuxBasic() {
	dnfNCurses := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "ncurses")
	dnfSSL := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "openssl")
	dnfSSH := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "openssh")

	if err := dnfNCurses.Run(); err != nil {
		checkError(err)
	}
	if err := dnfSSL.Run(); err != nil {
		checkError(err)
	}
	if err := dnfSSH.Run(); err != nil {
		checkError(err)
	}
}

func linuxEnv() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Setting basic environment..."
	ldBar.FinalMSG = " - Completed environment!\n"
	ldBar.Start()

	confA4s()
	newZProfile()
	newZshRC()

	profileAppend := "# Alias4sh\n" +
		"source ~/.config/alias4sh/aliasrc\n" +
		"# HOMEdnf\n" +
		"eval \"$(" + cmdPMS + " shellenv)\"\n"
	appendFile(profilePath, profileAppend)
	ldBar.Stop()
}

func linuxGit() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " s git..."
	ldBar.FinalMSG = " - Installed git!\n"
	ldBar.Start()

	dnfGit := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, cmdGit)
	dnfGitLfs := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "git-lfs")
	if err := dnfGit.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGitLfs.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func linuxTerminal() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing zsh with useful tools..."
	ldBar.FinalMSG = " - Installed useful tools for terminal!\n"
	ldBar.Start()

	dnfZsh := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "zsh")
	dnfTree := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "tree")
	dnfZshSyntax := exec.Command(cmdGit, gitClone, "https://github.com/zsh-users/zsh-syntax-highlighting.git", "~/.zsh/zsh-syntax-highlighting")
	dnfZshAuto := exec.Command(cmdGit, gitClone, "https://github.com/zsh-users/zsh-autosuggestions.git", "~/.zsh/zsh-autosuggestions")
	dnfZshComp := exec.Command(cmdGit, gitClone, "https://github.com/zsh-users/zsh-completions.git", "~/.zsh/zsh-completions")
	dnfZshTheme := exec.Command(cmdGit, gitClone, "https://github.com/romkatv/powerlevel10k.git", "~/.zsh/powerlevel10k")
	if err := dnfZsh.Run(); err != nil {
		checkError(err)
	}
	if err := dnfTree.Run(); err != nil {
		checkError(err)
	}
	if err := dnfZshSyntax.Run(); err != nil {
		checkError(err)
	}
	if err := dnfZshAuto.Run(); err != nil {
		checkError(err)
	}
	if err := dnfZshComp.Run(); err != nil {
		checkError(err)
	}
	if err := dnfZshTheme.Run(); err != nil {
		checkError(err)
	}

	ldBar.Stop()
}

func linuxDependency() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for development work..."
	ldBar.FinalMSG = " - Installed dependencies!\n"
	ldBar.Start()

	dnfKRB5 := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "krb5-workstation")
	dnfGnuPG := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "gnupg")
	dnfcURL := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "curl")
	dnfWget := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "wget")
	dnfXZ := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "xz")
	dnfGzip := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "unzip")
	dnfUnzip := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "gzidp")
	dnfLibzip := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "libzip")
	dnfBzip2 := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "bzip2")
	dnfZLib := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "zlib")
	dnfPkgConfig := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "pkg-config")
	dnfReadLine := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "readline")
	dnfLibffi := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "libffi")
	dnfUtilLinux := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "util-linux")
	dnfCoreUtils := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "coreutils")
	dnfBison := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "bison")
	dnfRe2C := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "re2c")
	dnfGD := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "gd")
	dnfCaCert := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "ca-certificates")
	dnfLDNS := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "ldns")
	dnfXMLto := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "xmlto")
	dnfGMP := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "gmp")
	dnfLibSodium := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "libsodium")
	dnfImageMagick := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "ImageMagick")
	dnfGhostscript := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "ghostscript")
	if err := dnfKRB5.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGnuPG.Run(); err != nil {
		checkError(err)
	}
	if err := dnfcURL.Run(); err != nil {
		checkError(err)
	}
	if err := dnfWget.Run(); err != nil {
		checkError(err)
	}
	if err := dnfXZ.Run(); err != nil {
		checkError(err)
	}
	if err := dnfUnzip.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGzip.Run(); err != nil {
		checkError(err)
	}
	if err := dnfLibzip.Run(); err != nil {
		checkError(err)
	}
	if err := dnfBzip2.Run(); err != nil {
		checkError(err)
	}
	if err := dnfZLib.Run(); err != nil {
		checkError(err)
	}
	if err := dnfPkgConfig.Run(); err != nil {
		checkError(err)
	}
	if err := dnfReadLine.Run(); err != nil {
		checkError(err)
	}
	if err := dnfLibffi.Run(); err != nil {
		checkError(err)
	}
	if err := dnfUtilLinux.Run(); err != nil {
		checkError(err)
	}
	if err := dnfCoreUtils.Run(); err != nil {
		checkError(err)
	}
	if err := dnfBison.Run(); err != nil {
		checkError(err)
	}
	if err := dnfRe2C.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGD.Run(); err != nil {
		checkError(err)
	}
	if err := dnfCaCert.Run(); err != nil {
		checkError(err)
	}
	if err := dnfLDNS.Run(); err != nil {
		checkError(err)
	}
	if err := dnfXMLto.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGMP.Run(); err != nil {
		checkError(err)
	}
	if err := dnfLibSodium.Run(); err != nil {
		checkError(err)
	}
	if err := dnfImageMagick.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGhostscript.Run(); err != nil {
		checkError(err)
	}

	ldBar.Stop()
}

func linuxDevToolCLI() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for CLI"
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	dnfGawk := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "gawk")
	dnfTig := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "tig")
	dnfJQ := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "jq")
	//dnfDirEnv := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "direnv")
	//dnfWatchman := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "watchman")
	dnfQEMU := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "qemu-kvm")
	dnfCCache := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "ccache")
	dnfMake := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "make")
	dnfVim := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "vim")
	dnfGH := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "gh")
	if err := dnfGawk.Run(); err != nil {
		checkError(err)
	}
	if err := dnfTig.Run(); err != nil {
		checkError(err)
	}
	if err := dnfJQ.Run(); err != nil {
		checkError(err)
	}
	//if err := dnfDirEnv.Run(); err != nil {
	//	checkError(err)
	//}
	//if err := dnfWatchman.Run(); err != nil {
	//	checkError(err)
	//}
	if err := dnfQEMU.Run(); err != nil {
		checkError(err)
	}
	if err := dnfCCache.Run(); err != nil {
		checkError(err)
	}
	if err := dnfMake.Run(); err != nil {
		checkError(err)
	}
	if err := dnfVim.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGH.Run(); err != nil {
		checkError(err)
	}

	shrcAppend := "# DIRENV\n" +
		"eval \"$(direnv hook zsh)\"\n\n"
	appendFile(shrcPath, shrcAppend)
	ldBar.Stop()
}

func linuxASDF() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing ASDF-VM with plugin..."
	ldBar.FinalMSG = " - Installed ASDF-VM, and add basic languages!\n"
	ldBar.Start()

	dnfASDF := exec.Command(cmdGit, gitClone, "https://github.com/asdf-vm/asdf.git", homeDir()+".asdf", "--branch", "v0.10.2")
	if err := dnfASDF.Run(); err != nil {
		checkError(err)
	}

	shrcAppend := "# DIRENV\n" +
		"source" + homeDir() + "/.asdf/asdf.sh\n" +
		"source " + homeDir() + "/.asdf/completions/asdf.bash\n\n"
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

func linuxServer() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server..."
	ldBar.FinalMSG = " - Installed server and database!\n"
	ldBar.Start()

	dnfHTTPD := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "httpd")
	dnfSQLite := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "sqlite")
	dnfPostgreSQL := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "postgresql")
	dnfMySQL := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "mysql-server")
	dnfRedis := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "redis")
	if err := dnfHTTPD.Run(); err != nil {
		checkError(err)
	}
	if err := dnfSQLite.Run(); err != nil {
		checkError(err)
	}
	if err := dnfPostgreSQL.Run(); err != nil {
		checkError(err)
	}
	if err := dnfMySQL.Run(); err != nil {
		checkError(err)
	}
	if err := dnfRedis.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func linuxLanguage() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing computer programming language..."
	ldBar.FinalMSG = " - Installed basic languages!\n"
	ldBar.Start()

	dnfPerl := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "perl")
	dnfRuby := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "ruby")
	dnfPython := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "python")
	dnfLua := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "lua")
	dnfGo := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "golang")
	dnfRust := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "rust")
	dnfNode := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "nodejs")
	dnfPHP := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "php")
	dnfJDK := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "java")
	dnfScala := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "scala") // Fedora
	dnfMaven := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "maven")
	dnfClojure := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "clojure") // Fedora
	dnfErlang := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "erlang")   // Fedora
	dnfElixir := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "elixir")   // Fedora
	if err := dnfPerl.Run(); err != nil {
		checkError(err)
	}
	if err := dnfRuby.Run(); err != nil {
		checkError(err)
	}
	if err := dnfPython.Run(); err != nil {
		checkError(err)
	}
	if err := dnfLua.Run(); err != nil {
		checkError(err)
	}
	if err := dnfGo.Run(); err != nil {
		checkError(err)
	}
	if err := dnfRust.Run(); err != nil {
		checkError(err)
	}
	if err := dnfNode.Run(); err != nil {
		checkError(err)
	}
	if err := dnfPHP.Run(); err != nil {
		checkError(err)
	}
	if err := dnfJDK.Run(); err != nil {
		checkError(err)
	}
	if err := dnfScala.Run(); err != nil {
		checkError(err)
	}
	if err := dnfMaven.Run(); err != nil {
		checkError(err)
	}
	if err := dnfClojure.Run(); err != nil {
		checkError(err)
	}
	if err := dnfErlang.Run(); err != nil {
		checkError(err)
	}
	if err := dnfElixir.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func linuxUtility() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing advanced utilities for terminal..."
	ldBar.FinalMSG = " - Installed advanced utilities!\n"
	ldBar.Start()

	dnfTmux := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "tmux")
	dnfNeofetch := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "neofetch")
	dnfAsciinema := exec.Command(superUser, cmdPMS, cmdIns, cmdYes, "asciinema")
	getFzf := exec.Command(cmdGit, gitClone, "https://github.com/junegunn/fzf.git", "--depth", "1", homeDir()+".fzf")
	installFzf := exec.Command(homeDir() + ".fzf/install")
	if err := dnfTmux.Run(); err != nil {
		checkError(err)
	}
	if err := dnfNeofetch.Run(); err != nil {
		checkError(err)
	}
	if err := dnfAsciinema.Run(); err != nil {
		checkError(err)
	}
	if err := getFzf.Run(); err != nil {
		checkError(err)
	}
	if err := installFzf.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func linuxEnd() {
	shrcAppend := "\n######## ADD CUSTOM VALUES UNDER HERE ########\n\n\n"
	appendFile(shrcPath, shrcAppend)
}

func main() {
	fmt.Println("\nDev4mac v" + appVer + "\n")
	if checkNetStatus() == true {
		fmt.Println("\nChoose setup type\n" +
			"\t1. Desktop setup\n" +
			"\t2. Server setup\n" +
			"\t3. Minimal setup\n" +
			"\t0. Quit\n")
	beginOpt:
		for {
			fmt.Printf(chooseCmd)
			_, err := fmt.Scanln(&cmdOpt)
			checkError(err)
			if cmdOpt == "1" {
				linuxBegin()
				linuxBasic()
				linuxEnv()
				linuxGit()
				linuxTerminal()
				linuxDependency()
				linuxDevToolCLI()
				linuxASDF()
				linuxServer()
				linuxLanguage()
				linuxUtility()
			} else if cmdOpt == "2" {
				linuxBegin()
				linuxBasic()
				linuxEnv()
				linuxGit()
				linuxDevToolCLI()
				linuxServer()
				linuxLanguage()
			} else if cmdOpt == "3" {
				linuxBegin()
				linuxBasic()
				linuxGit()
			} else if cmdOpt == "0" || cmdOpt == "q" || cmdOpt == "e" || cmdOpt == "quit" || cmdOpt == "exit" {
			} else {
				fmt.Println("Wrong answer. Please choose number 0-3")
				goto beginOpt
			}
			break
		}

		linuxEnd()
		fmt.Println("\nFinished to setup! You can choose 4 options. (Recommend option is 1)\n" +
			"\t1. Setup zsh theme & Configure git global\n" +
			"\t2. Only setup zsh theme that minimal type\n" +
			"\t3. Only configure git global easily\n" +
			"\t0. Nothing, finish Dev4mac (manual setup)\n")
	endOpt:
		for {
			fmt.Printf(chooseCmd)
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
