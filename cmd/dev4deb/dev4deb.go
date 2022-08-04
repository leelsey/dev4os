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
	cmdPMS      = "apt"
	pmsIns      = "install"
	//cmdReIns    = "reinstall"
	pmsRm      = "remove"
	pmsYes     = "-y"
	pmsConf    = "config-manager"
	pmsAddRepo = "--add-repo"
	cmdSys     = "systemctl"
	cmdEnable  = "enable"
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

func checkLinuxVer() string {
	checkDebFamily := exec.Command("cat", "/etc/lsb-release")
	if err := checkDebFamily.Run(); err != nil {
		checkError(err)
	}
	return ""
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

func updateapt() {
	updatePMS := exec.Command(cmdPMS, "makecache", "--refresh")
	installEpel := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "epel-release")
	installPlugins := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "apt-plugins-core")
	updateLinux := exec.Command(superUser, cmdPMS, "update", pmsYes)

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
	installFirewall := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "firewalld")
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

	updateapt()
	secureConf()
	ldBar.Stop()
}

func linuxBasic() {
	aptNCurses := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ncurses")
	aptNCursesDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ncurses-devel")
	aptSSL := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "openssl")
	aptSSLDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "openssl-devel")
	aptSSH := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "openssh")

	if err := aptNCurses.Run(); err != nil {
		checkError(err)
	}
	if err := aptNCursesDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptSSL.Run(); err != nil {
		checkError(err)
	}
	if err := aptSSLDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptSSH.Run(); err != nil {
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
		"# HOMEapt\n" +
		"eval \"$(" + cmdPMS + " shellenv)\"\n"
	appendFile(profilePath, profileAppend)
	ldBar.Stop()
}

func linuxGit() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " s git..."
	ldBar.FinalMSG = " - Installed git!\n"
	ldBar.Start()

	aptGit := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, cmdGit)
	aptGitLfs := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "git-lfs")
	if err := aptGit.Run(); err != nil {
		checkError(err)
	}
	if err := aptGitLfs.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func linuxTerminal() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing zsh with useful tools..."
	ldBar.FinalMSG = " - Installed useful tools for terminal!\n"
	ldBar.Start()

	aptZsh := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "zsh")
	aptTree := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "tree")
	aptZshSyntax := exec.Command(cmdGit, gitClone, "https://github.com/zsh-users/zsh-syntax-highlighting.git", "~/.zsh/zsh-syntax-highlighting")
	aptZshAuto := exec.Command(cmdGit, gitClone, "https://github.com/zsh-users/zsh-autosuggestions.git", "~/.zsh/zsh-autosuggestions")
	aptZshComp := exec.Command(cmdGit, gitClone, "https://github.com/zsh-users/zsh-completions.git", "~/.zsh/zsh-completions")
	aptZshTheme := exec.Command(cmdGit, gitClone, "https://github.com/romkatv/powerlevel10k.git", "~/.zsh/powerlevel10k")
	if err := aptZsh.Run(); err != nil {
		checkError(err)
	}
	if err := aptTree.Run(); err != nil {
		checkError(err)
	}
	if err := aptZshSyntax.Run(); err != nil {
		checkError(err)
	}
	if err := aptZshAuto.Run(); err != nil {
		checkError(err)
	}
	if err := aptZshComp.Run(); err != nil {
		checkError(err)
	}
	if err := aptZshTheme.Run(); err != nil {
		checkError(err)
	}

	ldBar.Stop()
}

func linuxDependency() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for development work..."
	ldBar.FinalMSG = " - Installed dependencies!\n"
	ldBar.Start()

	aptKRB5 := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "krb5-workstation")
	aptGnuPG := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gnupg")
	aptcURL := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "curl")
	aptWget := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "wget")
	aptXZ := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "xz")
	aptXZDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "xz-devel")
	aptGzip := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "unzip")
	aptUnzip := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gzidp")
	aptLibzip := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libzip")
	aptBzip2 := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "bzip2")
	aptBzip2Dev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "bzip2-devel")
	aptZLib := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "zlib")
	aptZLibDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "zlib-devel")
	aptLibYaml := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libyaml")
	aptPkgConfig := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "pkg-config")
	aptReadLine := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "readline")
	aptReadLineDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "readline-devel")
	aptLibffi := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libffi")
	aptLibffiDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libffi-devel")
	aptLibcURL := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libcurl")
	aptLibcURLDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libcurl-devel")
	aptLibAvif := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libavif")
	aptLibWebP := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libwebp")
	aptLibJpeg := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libjpeg")
	aptLibXpm := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libXpm")
	aptUtilLinux := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "util-linux")
	aptCoreUtils := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "coreutils")
	aptOniguruma := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "oniguruma")
	aptOnigurumaDev := exec.Command(superUser, cmdPMS, "--enablerepo=crb", pmsIns, pmsYes, "oniguruma-devel")
	aptBison := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "bison")
	aptRe2C := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "re2c")
	aptGD := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gd")
	aptGDDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gd-devel")
	aptPerlGD := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "perl-GD")
	aptCaCert := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ca-certificates")
	aptLDNS := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ldns")
	aptXMLto := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "xmlto")
	aptGMP := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gmp")
	aptLibSodium := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libsodium")
	aptImageMagick := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ImageMagick")
	aptGhostscript := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ghostscript")

	if err := aptKRB5.Run(); err != nil {
		checkError(err)
	}
	if err := aptGnuPG.Run(); err != nil {
		checkError(err)
	}
	if err := aptcURL.Run(); err != nil {
		checkError(err)
	}
	if err := aptWget.Run(); err != nil {
		checkError(err)
	}
	if err := aptXZ.Run(); err != nil {
		checkError(err)
	}
	if err := aptXZDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptUnzip.Run(); err != nil {
		checkError(err)
	}
	if err := aptGzip.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibzip.Run(); err != nil {
		checkError(err)
	}
	if err := aptBzip2.Run(); err != nil {
		checkError(err)
	}
	if err := aptBzip2Dev.Run(); err != nil {
		checkError(err)
	}
	if err := aptZLib.Run(); err != nil {
		checkError(err)
	}
	if err := aptZLibDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibYaml.Run(); err != nil {
		checkError(err)
	}
	if err := aptPkgConfig.Run(); err != nil {
		checkError(err)
	}
	if err := aptReadLine.Run(); err != nil {
		checkError(err)
	}
	if err := aptReadLineDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibffi.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibffiDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibcURL.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibcURLDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibAvif.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibWebP.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibJpeg.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibXpm.Run(); err != nil {
		checkError(err)
	}
	if err := aptUtilLinux.Run(); err != nil {
		checkError(err)
	}
	if err := aptCoreUtils.Run(); err != nil {
		checkError(err)
	}
	if err := aptOniguruma.Run(); err != nil {
		checkError(err)
	}
	if err := aptOnigurumaDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptBison.Run(); err != nil {
		checkError(err)
	}
	if err := aptRe2C.Run(); err != nil {
		checkError(err)
	}
	if err := aptGD.Run(); err != nil {
		checkError(err)
	}
	if err := aptGDDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptPerlGD.Run(); err != nil {
		checkError(err)
	}
	if err := aptCaCert.Run(); err != nil {
		checkError(err)
	}
	if err := aptLDNS.Run(); err != nil {
		checkError(err)
	}
	if err := aptXMLto.Run(); err != nil {
		checkError(err)
	}
	if err := aptGMP.Run(); err != nil {
		checkError(err)
	}
	if err := aptLibSodium.Run(); err != nil {
		checkError(err)
	}
	if err := aptImageMagick.Run(); err != nil {
		checkError(err)
	}
	if err := aptGhostscript.Run(); err != nil {
		checkError(err)
	}

	if checkLinuxVer() == "ubuntu" {
		aptLibYamlDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "libyaml-devel")
		aptGDBM := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gdbm")
		aptGDBMDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gdbm-devel")

		if err := aptLibYamlDev.Run(); err != nil {
			checkError(err)
		}
		if err := aptGDBM.Run(); err != nil {
			checkError(err)
		}
		if err := aptGDBMDev.Run(); err != nil {
			checkError(err)
		}
	}

	ldBar.Stop()
}

func linuxDevToolCLI() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for CLI"
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	aptGawk := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gawk")
	aptTig := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "tig")
	aptJQ := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "jq")
	aptQEMU := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "qemu-kvm")
	aptCCache := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ccache")
	aptMake := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "make")
	aptCMake := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "cmake")
	aptGCC := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gcc")
	aptGCCCpp := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gcc-c++")
	aptAnt := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ant")
	aptMaven := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "maven")
	aptTk := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "tk")
	aptTkDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "tk-devel")
	aptVim := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "vim")
	aptGH := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "gh")

	if err := aptGawk.Run(); err != nil {
		checkError(err)
	}
	if err := aptTig.Run(); err != nil {
		checkError(err)
	}
	if err := aptJQ.Run(); err != nil {
		checkError(err)
	}
	if err := aptQEMU.Run(); err != nil {
		checkError(err)
	}
	if err := aptCCache.Run(); err != nil {
		checkError(err)
	}
	if err := aptMake.Run(); err != nil {
		checkError(err)
	}
	if err := aptCMake.Run(); err != nil {
		checkError(err)
	}
	if err := aptGCC.Run(); err != nil {
		checkError(err)
	}
	if err := aptGCCCpp.Run(); err != nil {
		checkError(err)
	}
	if err := aptAnt.Run(); err != nil {
		checkError(err)
	}
	if err := aptMaven.Run(); err != nil {
		checkError(err)
	}
	if err := aptTk.Run(); err != nil {
		checkError(err)
	}
	if err := aptTkDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptVim.Run(); err != nil {
		checkError(err)
	}
	if err := aptGH.Run(); err != nil {
		checkError(err)
	}

	RMOldDocker := exec.Command(superUser, cmdPMS, pmsRm, pmsYes, "docker", "docker-client", "docker-client-latest", "docker-common", "docker-latest", "docker-latest-logrotate", "docker-logrotate", "docker-engine-selinux", "docker-engine-selinux", "docker-engine")

	if err := RMOldDocker.Run(); err != nil {
		checkError(err)
	}

	if checkLinuxVer() == "ubuntu" {
		aptDirEnv := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "direnv")
		aptWatchman := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "watchman")

		if err := aptDirEnv.Run(); err != nil {
			checkError(err)
		}
		if err := aptWatchman.Run(); err != nil {
			checkError(err)
		}

		aptRepoDocker := exec.Command(superUser, cmdPMS, pmsConf, pmsAddRepo, "https://download.docker.com/linux/ubuntu/docker-ce.repo")
		if err := aptRepoDocker.Run(); err != nil {
			checkError(err)
		}
	} else if checkLinuxVer() == "debian" {
		aptRepoDocker := exec.Command(superUser, cmdPMS, pmsConf, pmsAddRepo, "https://download.docker.com/linux/debian/docker-ce.repo")

		if err := aptRepoDocker.Run(); err != nil {
			checkError(err)
		}
	}

	aptDocker := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "docker-ce", "docker-ce-cli", "containerd.io", "docker-compose-plugin")

	if err := aptDocker.Run(); err != nil {
		checkError(err)
	}

	//shrcAppend := "# DIRENV\n" +
	//	"eval \"$(direnv hook zsh)\"\n\n"
	//appendFile(shrcPath, shrcAppend)
	ldBar.Stop()
}

func linuxASDF() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing ASDF-VM with plugin..."
	ldBar.FinalMSG = " - Installed ASDF-VM, and add basic languages!\n"
	ldBar.Start()

	aptASDF := exec.Command(cmdGit, gitClone, "https://github.com/asdf-vm/asdf.git", homeDir()+".asdf", "--branch", "v0.10.2")
	if err := aptASDF.Run(); err != nil {
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

	aptHTTPD := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "httpd")
	aptSQLite := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "sqlite")
	aptSQLiteDev := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "sqlite-devel")
	aptPostgreSQL := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "postgresql")
	aptRedis := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "redis")

	if err := aptHTTPD.Run(); err != nil {
		checkError(err)
	}
	if err := aptSQLite.Run(); err != nil {
		checkError(err)
	}
	if err := aptSQLiteDev.Run(); err != nil {
		checkError(err)
	}
	if err := aptPostgreSQL.Run(); err != nil {
		checkError(err)
	}
	if err := aptRedis.Run(); err != nil {
		checkError(err)
	}

	if checkLinuxVer() == "ubuntu" {
		aptMySQL := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "community-mysql-server")

		if err := aptMySQL.Run(); err != nil {
			checkError(err)
		}
	} else if checkLinuxVer() == "debian" {
		aptMySQL := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "mysql-server")

		if err := aptMySQL.Run(); err != nil {
			checkError(err)
		}
	}
	ldBar.Stop()
}

func linuxLanguage() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing computer programming language..."
	ldBar.FinalMSG = " - Installed basic languages!\n"
	ldBar.Start()

	aptPerl := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "perl")
	aptRuby := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "ruby")
	aptPython := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "python")
	aptLua := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "lua")
	aptGo := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "golang")
	aptRust := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "rust")
	aptNode := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "nodejs")
	aptPHP := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "php")
	aptJDK := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "java")

	if err := aptPerl.Run(); err != nil {
		checkError(err)
	}
	if err := aptRuby.Run(); err != nil {
		checkError(err)
	}
	if err := aptPython.Run(); err != nil {
		checkError(err)
	}
	if err := aptLua.Run(); err != nil {
		checkError(err)
	}
	if err := aptGo.Run(); err != nil {
		checkError(err)
	}
	if err := aptRust.Run(); err != nil {
		checkError(err)
	}
	if err := aptNode.Run(); err != nil {
		checkError(err)
	}
	if err := aptPHP.Run(); err != nil {
		checkError(err)
	}
	if err := aptJDK.Run(); err != nil {
		checkError(err)
	}

	if checkLinuxVer() == "ubuntu" {
		aptScala := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "scala")
		aptClojure := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "clojure")
		aptErlang := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "erlang")
		aptElixir := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "elixir")

		if err := aptScala.Run(); err != nil {
			checkError(err)
		}
		if err := aptClojure.Run(); err != nil {
			checkError(err)
		}
		if err := aptErlang.Run(); err != nil {
			checkError(err)
		}
		if err := aptElixir.Run(); err != nil {
			checkError(err)
		}
	}
	ldBar.Stop()
}

func linuxUtility() {
	ldBar := spinner.New(spinner.CharSets[16], 50*time.Millisecond)
	ldBar.Suffix = " Installing advanced utilities for terminal..."
	ldBar.FinalMSG = " - Installed advanced utilities!\n"
	ldBar.Start()

	aptTmux := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "tmux")
	aptNeofetch := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "neofetch")
	aptAsciinema := exec.Command(superUser, cmdPMS, pmsIns, pmsYes, "asciinema")
	getFzf := exec.Command(cmdGit, gitClone, "https://github.com/junegunn/fzf.git", "--depth", "1", homeDir()+".fzf")
	installFzf := exec.Command(homeDir() + ".fzf/install")
	if err := aptTmux.Run(); err != nil {
		checkError(err)
	}
	if err := aptNeofetch.Run(); err != nil {
		checkError(err)
	}
	if err := aptAsciinema.Run(); err != nil {
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
