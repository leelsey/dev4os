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
	dlPath   = homeDir() + "Downloads\\"
	pSh      = "powershell"
	cmdPMS   = "C:\\ProgramData\\chocolatey\\choco.exe"
	cmdIn    = "install"
	cmdRein  = "reinstall"
	cmdRm    = "uninstall"
	cmdY     = "-y"
	cmdGit   = "C:\\'Program Files'\\git\\bin\\git.exe"
	cmdOpt   string
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

func checkAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func restartWin() {
	fmt.Println("Restarting now ...")
	if err := exec.Command(pSh, "shutdown", "/r", "/t", "0").Run(); err != nil {
		fmt.Println(" - Failed to restart Windows")
		winEnd()
	}
}

func runElevated() {
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

func confG4s() {
	dlG4s := dlPath + Git4setV + ".zip"
	req, _ := http.NewRequest("GET", "https://github.com/leelsey/Git4set/archive/refs/tags/v0.1.zip", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	file, _ := os.OpenFile(dlG4s, os.O_CREATE|os.O_WRONLY, 0755)
	defer file.Close()
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Git4set...",
	)
	io.Copy(io.MultiWriter(file, bar), resp.Body)

	fmt.Println(" - Finished to download Git4sh: " + dlG4s + " (Your download directory)\n" +
		"\nPlease extract zip file and run script on terminal.\n" +
		lstDot + "Configure global author & ignore: initial-git\n" +
		lstDot + "Only want configure global author: git-conf\n" +
		lstDot + "Only want configure global ignore: git-ignore")
}

func winChoco() {
	if _, err := os.Stat("C:\\ProgramData\\Chocolatey"); !os.IsNotExist(err) {
		ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
		ldBar.Suffix = " Updating chocolatey..."
		ldBar.FinalMSG = " - Updated choco!\n"
		ldBar.Start()

		updateChocolatey := exec.Command(cmdPMS, "upgrade", cmdY, "all")
		updatingHomebrew, err := updateChocolatey.Output()
		checkError(err)
		fmt.Sprintf(string(updatingHomebrew))
		ldBar.Stop()
	} else {
		ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
		ldBar.Suffix = " Installing chocolatey..."
		ldBar.FinalMSG = " - Installed choco!\n"
		ldBar.Start()

		installChocolatey := exec.Command(pSh, "Set-ExecutionPolicy", "Bypass", "-Scope", "Process", "-Force;", "[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol", "-bor", "3072;", "iex", "((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))")
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

	installGit := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "git")
	installGitLfs := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "git-lfs")
	gitLfsInstall := exec.Command(pSh, cmdGit, "lfs", "install")
	gitBranchMain := exec.Command(pSh, cmdGit, "config", "--global", "init.defaultBranch", "main")

	installingGit, err := installGit.Output()
	checkError(err)
	installingGitLfs, err := installGitLfs.Output()
	checkError(err)
	gitLfsInstalling, err := gitLfsInstall.Output()
	checkError(err)
	confGitBranchMain, err := gitBranchMain.Output()
	checkError(err)

	fmt.Sprintf(string(installingGit))
	fmt.Sprintf(string(installingGitLfs))
	fmt.Sprintf(string(gitLfsInstalling))
	fmt.Sprintf(string(confGitBranchMain))
	ldBar.Stop()
}

func winDependency() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for development work..."
	ldBar.FinalMSG = " - Installed dependencies!\n"
	ldBar.Start()

	installSSL := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "openssl")
	installGnuPG := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "gnupg")
	installcURL := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "curl")
	installWget := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "wget")
	installGzip := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "gzip")
	installBzip2 := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "bzip2")
	installCoreUtils := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "gnuwin32-coreutils.install")
	installRe2C := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "re2c")
	installGhostscript := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "ghostscript")

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

	installGawk := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "gawk")
	installJQ := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "jq")
	installWatchman := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "watchman")
	installQEMU := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "qemu")
	installCcache := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "ccache")
	installMake := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "make")
	installVim := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "vim")
	installBat := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "bat")
	installGH := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "gh")
	installPS := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "powershell")
	installCygwin := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "cygwin")
	installVS := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "visualstudio2022community")

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
	installingPS, err := installPS.Output()
	checkError(err)
	installingCygwin, err := installCygwin.Output()
	checkError(err)
	installingVS, err := installVS.Output()
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
	fmt.Sprintf(string(installingPS))
	fmt.Sprintf(string(installingCygwin))
	fmt.Sprintf(string(installingVS))
	ldBar.Stop()
}

func winServer() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server..."
	ldBar.FinalMSG = " - Installed server and database!\n"
	ldBar.Start()

	installHTTPD := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "apache-httpd")
	installTomcat := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "tomcat")
	installSQLite := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "sqlite")
	installPostgreSQL := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "postgresql")
	installMySQL := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "mysql")

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

	installGCC := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "mingw")
	installLLVM := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "llvm")
	installNuGet := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "nuget.commandline")
	installPerl := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "strawberryperl")
	installRuby := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "ruby")
	installPython := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "python")
	installLua := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "lua")
	installGo := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "go")
	installRust := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "rust")
	installNode := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "nodejs")
	installPHP := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "php")
	installJDK := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "openjdk")
	installGroovy := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "groovy")
	installScala := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "scala")
	installClojure := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "clojure")
	installErlang := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "erlang")
	installElixir := exec.Command(pSh, cmdPMS, cmdIn, cmdY, "elixir")

	installingGCC, err := installGCC.Output()
	checkError(err)
	installingLLVM, err := installLLVM.Output()
	checkError(err)
	installingNuGet, err := installNuGet.Output()
	checkError(err)
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

	fmt.Sprintf(string(installingGCC))
	fmt.Sprintf(string(installingLLVM))
	fmt.Sprintf(string(installingNuGet))
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

func winWLS() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing Windows Subsystem for Linux with Ubuntu..."
	ldBar.FinalMSG = " - Installed WSL with Ubuntu!\n"
	ldBar.Start()

	setWSL := exec.Command(pSh, "wsl", "--install")
	settingWSL, err := setWSL.Output()
	checkError(err)
	fmt.Sprintf(string(settingWSL))
	ldBar.Stop()
}

func winEnd() {
	fmt.Println("\n----------Finished!----------\n" +
		"Please RESTART your terminal and OS!\n" +
		lstDot + "Restart the terminal (CMD or PowerShell) for the changes to take effect.\n" +
		lstDot + "WSL has been setup. Restart OS for the changes to take effect.")
}

func main() {
	if !checkAdmin() {
		runElevated()
	}
	if checkAdmin() {
		fmt.Println("\nDev4win v" + appVer + "\n")
		winChoco()
		winGit()
		winDependency()
		winDevToolCLI()
		winServer()
		winLanguage()
		winWLS()
		fmt.Println("\nFinished to setup! You can choose 4 options. (Recommend option is 1)\n" +
			"\t1. Restart OS after download Git4set\n" +
			"\t2. Restart Windows operating system\n" +
			"\t3. Download easily configure global git (Git4set)\n" +
			"\t0. Nothing, finish Dev4win")
		for {
			fmt.Printf("\nSelect command: ")
			fmt.Scanln(&cmdOpt)
			if cmdOpt == "1" {
				confG4s()
				restartWin()
				break
			} else if cmdOpt == "2" {
				restartWin()
				break
			} else if cmdOpt == "3" {
				confG4s()
				winEnd()
				break
			} else if cmdOpt == "0" {
				winEnd()
				break
			} else {
				fmt.Println("Wrong answer. Please choose between 1,2,3,0.")
			}
		}
	}
}
