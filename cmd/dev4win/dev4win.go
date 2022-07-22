package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sys/windows"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	appVer   = "0.1"
	Git4setV = "Git4set-0.1"
	lstDot   = " • "
	cmdPMS   = "choco"
	cmdIn    = "install"
	cmdRein  = "reinstall"
	cmdRm    = "uninstall"
	cmdYes   = "-y"
	cmdEcho  = "echo"
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
	return homeDirPath + "\\"
}

func workingDir() string {
	workingDirPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return workingDirPath + "\\"
}

func confGit4setWin() {
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
	mvLoc := homeDir() + "Downloads\\" + Git4setV + ".zip"
	err := os.Rename(dlLoc, mvLoc)
	checkError(err)

	fmt.Println(" - Finished to download Git4sh: " + mvLoc + " (Your download directory)\n" +
		"\nPlease extract zip file and run shell script on terminal.\n" +
		lstDot + "Configure global author & ignore: .\\initial-git.bat\n" +
		lstDot + "Only want configure global author: .\\git-conf.bat\n" +
		lstDot + "Only want configure global ignore: .\\git-ignore.bat")
}

func winChoco() {
	if _, err := os.Stat("C:\\ProgramData\\Chocolatey"); !os.IsNotExist(err) {
		ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
		ldBar.Suffix = " Updating chocolatey..."
		ldBar.FinalMSG = " - Updated choco!\n"
		ldBar.Start()

		updateChocolatey := exec.Command(cmdPMS, "upgrade", cmdYes, "all")
		updatingHomebrew, err := updateChocolatey.Output()
		checkError(err)
		fmt.Sprintf(string(updatingHomebrew))
		ldBar.Stop()
	} else {
		ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
		ldBar.Suffix = " Installing chocolatey..."
		ldBar.FinalMSG = " - Installed choco!\n"
		ldBar.Start()

		installChocolatey := exec.Command("PS", "/C", "Set-ExecutionPolicy", "Bypass", "-Scope", "Process", "-Force;", "[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol", "-bor", "3072;iex", "((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))")
		installingChocolatey, err := installChocolatey.Output()
		checkError(err)
		fmt.Sprintf(string(installingChocolatey))
		ldBar.Stop()
	}
}

func winGit() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing git..."
	ldBar.FinalMSG = " - Installed git!\n"
	ldBar.Start()

	installGit := exec.Command(cmdPMS, cmdIn, cmdYes, "git")
	installGitLfs := exec.Command(cmdPMS, cmdIn, cmdYes, "git-lfs")
	gitLfsInstall := exec.Command("git", "lfs", "install")

	installingGit, err := installGit.Output()
	checkError(err)
	installingGitLfs, err := installGitLfs.Output()
	checkError(err)
	gitLfsInstalling, err := gitLfsInstall.Output()
	checkError(err)

	fmt.Sprintf(string(installingGit))
	fmt.Sprintf(string(installingGitLfs))
	fmt.Sprintf(string(gitLfsInstalling))
	ldBar.Stop()
}

func winDependency() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for development work..."
	ldBar.FinalMSG = " - Installed dependencies!\n"
	ldBar.Start()

	installSSL := exec.Command(cmdPMS, cmdIn, cmdYes, "openssl")
	installGnuPG := exec.Command(cmdPMS, cmdIn, cmdYes, "gnupg")
	installcURL := exec.Command(cmdPMS, cmdIn, cmdYes, "curl")
	installWget := exec.Command(cmdPMS, cmdIn, cmdYes, "wget")
	installGzip := exec.Command(cmdPMS, cmdIn, cmdYes, "gzip")
	installBzip2 := exec.Command(cmdPMS, cmdIn, cmdYes, "bzip2")
	installCoreUtils := exec.Command(cmdPMS, cmdIn, cmdYes, "gnuwin32-coreutils.install")
	installRe2C := exec.Command(cmdPMS, cmdIn, cmdYes, "re2c")
	installGhostscript := exec.Command(cmdPMS, cmdIn, cmdYes, "ghostscript")

	installingSSL, err := installSSL.Output()
	checkError(err)
	installingGnuPG, err := installGnuPG.Output()
	checkError(err)
	installingcURL, err := installcURL.Output()
	checkError(err)
	installingWget, err := installWget.Output()
	checkError(err)
	installingGzip, err := installGzip.Output()
	checkError(err)
	installingBzip2, err := installBzip2.Output()
	checkError(err)
	installingCoreUtils, err := installCoreUtils.Output()
	checkError(err)
	installingRe2C, err := installRe2C.Output()
	checkError(err)
	installingGhostscript, err := installGhostscript.Output()
	checkError(err)

	fmt.Sprintf(string(installingSSL))
	fmt.Sprintf(string(installingGnuPG))
	fmt.Sprintf(string(installingcURL))
	fmt.Sprintf(string(installingWget))
	fmt.Sprintf(string(installingGzip))
	fmt.Sprintf(string(installingBzip2))
	fmt.Sprintf(string(installingCoreUtils))
	fmt.Sprintf(string(installingRe2C))
	fmt.Sprintf(string(installingGhostscript))
	ldBar.Stop()
}

func winDevToolCLI() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for CLI..."
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	installGawk := exec.Command(cmdPMS, cmdIn, cmdYes, "gawk")
	installJQ := exec.Command(cmdPMS, cmdIn, cmdYes, "jq")
	installWatchman := exec.Command(cmdPMS, cmdIn, cmdYes, "watchman")
	installQEMU := exec.Command(cmdPMS, cmdIn, cmdYes, "qemu")
	installCcache := exec.Command(cmdPMS, cmdIn, cmdYes, "ccache")
	installMake := exec.Command(cmdPMS, cmdIn, cmdYes, "make")
	installVim := exec.Command(cmdPMS, cmdIn, cmdYes, "vim")
	installBat := exec.Command(cmdPMS, cmdIn, cmdYes, "bat")
	installGH := exec.Command(cmdPMS, cmdIn, cmdYes, "gh")

	installingGawk, err := installGawk.Output()
	checkError(err)
	installingJQ, err := installJQ.Output()
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
	fmt.Sprintf(string(installingJQ))
	fmt.Sprintf(string(installingWatchman))
	fmt.Sprintf(string(installingQEMU))
	fmt.Sprintf(string(installingCcache))
	fmt.Sprintf(string(installingMake))
	fmt.Sprintf(string(installingVim))
	fmt.Sprintf(string(installingBat))
	fmt.Sprintf(string(installingGH))
	ldBar.Stop()
}

func winServer() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server..."
	ldBar.FinalMSG = " - Installed server and database!\n"
	ldBar.Start()

	installHTTPD := exec.Command(cmdPMS, cmdIn, cmdYes, "apache-httpd")
	installTomcat := exec.Command(cmdPMS, cmdIn, cmdYes, "tomcat")
	installSQLite := exec.Command(cmdPMS, cmdIn, cmdYes, "sqlite")
	installPostgreSQL := exec.Command(cmdPMS, cmdIn, cmdYes, "postgresql")
	installMySQL := exec.Command(cmdPMS, cmdIn, cmdYes, "mysql")

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

	fmt.Sprintf(string(installingHTTPD))
	fmt.Sprintf(string(installingTomcat))
	fmt.Sprintf(string(installingSQLite))
	fmt.Sprintf(string(installingPostgreSQL))
	fmt.Sprintf(string(installingMySQL))
	ldBar.Stop()
}

func winLanguage() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing computer programming language..."
	ldBar.FinalMSG = " - Installed basic languages!\n"
	ldBar.Start()

	installPerl := exec.Command(cmdPMS, cmdIn, cmdYes, "strawberryperl")
	installRuby := exec.Command(cmdPMS, cmdIn, cmdYes, "ruby")
	installPython := exec.Command(cmdPMS, cmdIn, cmdYes, "python")
	installLua := exec.Command(cmdPMS, cmdIn, cmdYes, "lua")
	installGo := exec.Command(cmdPMS, cmdIn, cmdYes, "go")
	installRust := exec.Command(cmdPMS, cmdIn, cmdYes, "rust")
	installNode := exec.Command(cmdPMS, cmdIn, cmdYes, "nodejs")
	installPHP := exec.Command(cmdPMS, cmdIn, cmdYes, "php")
	installJDK := exec.Command(cmdPMS, cmdIn, cmdYes, "openjdk")
	installGroovy := exec.Command(cmdPMS, cmdIn, cmdYes, "groovy")
	installScala := exec.Command(cmdPMS, cmdIn, cmdYes, "scala")
	installClojure := exec.Command(cmdPMS, cmdIn, cmdYes, "clojure")
	installErlang := exec.Command(cmdPMS, cmdIn, cmdYes, "erlang")
	installElixir := exec.Command(cmdPMS, cmdIn, cmdYes, "elixir")

	installingPerl, err := installPerl.Output()
	checkError(err)
	installingRuby, err := installRuby.Output()
	checkError(err)
	installingPython, err := installPython.Output()
	checkError(err)
	installingLua, err := installLua.Output()
	checkError(err)
	installingGo, err := installGo.Output()
	checkError(err)
	installingRust, err := installRust.Output()
	checkError(err)
	installingNode, err := installNode.Output()
	checkError(err)
	installingPHP, err := installPHP.Output()
	checkError(err)
	installingJDK, err := installJDK.Output()
	checkError(err)
	installingGroovy, err := installGroovy.Output()
	checkError(err)
	installingScala, err := installScala.Output()
	checkError(err)
	installingClojure, err := installClojure.Output()
	checkError(err)
	installingErlang, err := installErlang.Output()
	checkError(err)
	installingElixir, err := installElixir.Output()
	checkError(err)

	fmt.Sprintf(string(installingPerl))
	fmt.Sprintf(string(installingRuby))
	fmt.Sprintf(string(installingPython))
	fmt.Sprintf(string(installingLua))
	fmt.Sprintf(string(installingGo))
	fmt.Sprintf(string(installingRust))
	fmt.Sprintf(string(installingNode))
	fmt.Sprintf(string(installingPHP))
	fmt.Sprintf(string(installingJDK))
	fmt.Sprintf(string(installingGroovy))
	fmt.Sprintf(string(installingScala))
	fmt.Sprintf(string(installingClojure))
	fmt.Sprintf(string(installingErlang))
	fmt.Sprintf(string(installingElixir))
	ldBar.Stop()
}

func winEnd() {
	fmt.Println("\n----------Finished!----------\n" +
		"Please RESTART your terminal!\n" +
		lstDot + "Restart the Terminal by yourself.")
}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")
	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)
	var showCmd int32 = 1
	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	checkError(err)
}

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func main() {
	if !amAdmin() {
		runMeElevated()
	}
	if amAdmin() {
		fmt.Println("\nDev4win v" + appVer + "\n")
		winChoco()
		winGit()
		winDependency()
		winDevToolCLI()
		winServer()
		winLanguage()
		fmt.Printf("\nPress any key to finish, " +
			"or press (i) if you want configure global git... ")
		var setCMD string
		fmt.Scanln(&setCMD)
		if setCMD == "i" || setCMD == "I" {
			confGit4setWin()
		}
		winEnd()
	}
}
