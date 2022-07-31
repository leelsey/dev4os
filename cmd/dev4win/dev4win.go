package main

import (
	"bufio"
	"fmt"
	"github.com/briandowns/spinner"
	"golang.org/x/sys/windows"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	appVer = "0.1"
	lstDot = " • "
	pSh    = "powershell"
	cmdPMS = "C:\\ProgramData\\chocolatey\\choco.exe"
	cmdIn  = "install"
	//cmdRein  = "reinstall"
	//cmdRm    = "uninstall"
	cmdY   = "-y"
	cmdGit = "C:\\'Program Files'\\git\\bin\\git.exe"
	cmdOpt string
)

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
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
	}
	os.Exit(0)
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
	fmt.Println("\nGit global configuration")

	fmt.Println(" 1) Main branch default name changed master -> main")
	setBranchMain := exec.Command("git", "config", "--global", "init.defaultBranch", "main")
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

	setUserName := exec.Command("git", "config", "--global", "user.name", userName)
	setUserEmail := exec.Command("git", "config", "--global", "user.email", userEmail)
	if err := setUserName.Run(); err != nil {
		checkError(err)
	}
	if err := setUserEmail.Run(); err != nil {
		checkError(err)
	}

	fmt.Println(" 3) Setup git global ignore file with directories")
	ignoreDir := homeDir() + ".config\\git\\"
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

	setExcludesFile := exec.Command("git", "config", "--global", "core.excludesfile", ignorePath)
	if err := setExcludesFile.Run(); err != nil {
		fmt.Println("error2")
		checkError(err)
	}

	fmt.Println(" " + lstDot + "Make \"gitignore_global\" file in " + ignoreDir)
}

func updateChoco() {
	updateChocolatey := exec.Command(cmdPMS, "upgrade", cmdY, "all")
	if err := updateChocolatey.Run(); err != nil {
		checkError(err)
	}
}

func installChoco() {
	installChocolatey := exec.Command(pSh, `Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))`)
	if err := installChocolatey.Run(); err != nil {
		checkError(err)
	}
}

func winBegin() {
	if _, err := os.Stat("C:\\ProgramData\\Chocolatey"); !os.IsNotExist(err) {
		ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
		ldBar.Suffix = " Updating chocolatey..."
		ldBar.FinalMSG = " - Updated choco!\n"
		ldBar.Start()

		updateChoco()
		ldBar.Stop()
	} else {
		ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
		ldBar.Suffix = " Installing chocolatey..."
		ldBar.FinalMSG = " - Installed choco!\n"
		ldBar.Start()

		installChoco()
		updateChoco()
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

	if err := installGit.Run(); err != nil {
		checkError(err)
	}
	if err := installGitLfs.Run(); err != nil {
		checkError(err)
	}
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

	if err := installSSL.Run(); err != nil {
		checkError(err)
	}
	if err := installGnuPG.Run(); err != nil {
		checkError(err)
	}
	if err := installcURL.Run(); err != nil {
		checkError(err)
	}
	if err := installWget.Run(); err != nil {
		checkError(err)
	}
	if err := installGzip.Run(); err != nil {
		checkError(err)
	}
	if err := installBzip2.Run(); err != nil {
		checkError(err)
	}
	if err := installCoreUtils.Run(); err != nil {
		checkError(err)
	}
	if err := installRe2C.Run(); err != nil {
		checkError(err)
	}
	if err := installGhostscript.Run(); err != nil {
		checkError(err)
	}
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

	if err := installGawk.Run(); err != nil {
		checkError(err)
	}
	if err := installJQ.Run(); err != nil {
		checkError(err)
	}
	if err := installWatchman.Run(); err != nil {
		checkError(err)
	}
	if err := installQEMU.Run(); err != nil {
		checkError(err)
	}
	if err := installCcache.Run(); err != nil {
		checkError(err)
	}
	if err := installMake.Run(); err != nil {
		checkError(err)
	}
	if err := installVim.Run(); err != nil {
		checkError(err)
	}
	if err := installBat.Run(); err != nil {
		checkError(err)
	}
	if err := installGH.Run(); err != nil {
		checkError(err)
	}
	if err := installPS.Run(); err != nil {
		checkError(err)
	}
	if err := installCygwin.Run(); err != nil {
		checkError(err)
	}
	if err := installVS.Run(); err != nil {
		checkError(err)
	}
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

	if err := installHTTPD.Run(); err != nil {
		checkError(err)
	}
	if err := installTomcat.Run(); err != nil {
		checkError(err)
	}
	if err := installSQLite.Run(); err != nil {
		checkError(err)
	}
	if err := installPostgreSQL.Run(); err != nil {
		checkError(err)
	}
	if err := installMySQL.Run(); err != nil {
		checkError(err)
	}
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

	if err := installGCC.Run(); err != nil {
		checkError(err)
	}
	if err := installLLVM.Run(); err != nil {
		checkError(err)
	}
	if err := installNuGet.Run(); err != nil {
		checkError(err)
	}
	if err := installPerl.Run(); err != nil {
		checkError(err)
	}
	if err := installRuby.Run(); err != nil {
		checkError(err)
	}
	if err := installPython.Run(); err != nil {
		checkError(err)
	}
	if err := installLua.Run(); err != nil {
		checkError(err)
	}
	if err := installGo.Run(); err != nil {
		checkError(err)
	}
	if err := installRust.Run(); err != nil {
		checkError(err)
	}
	if err := installNode.Run(); err != nil {
		checkError(err)
	}
	if err := installPHP.Run(); err != nil {
		checkError(err)
	}
	if err := installJDK.Run(); err != nil {
		checkError(err)
	}
	if err := installGroovy.Run(); err != nil {
		checkError(err)
	}
	if err := installScala.Run(); err != nil {
		checkError(err)
	}
	if err := installClojure.Run(); err != nil {
		checkError(err)
	}
	if err := installErlang.Run(); err != nil {
		checkError(err)
	}
	if err := installElixir.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func winWLS() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing Windows Subsystem for Linux with Ubuntu..."
	ldBar.FinalMSG = " - Installed WSL2 with Ubuntu!\n"
	ldBar.Start()

	setWSL := exec.Command(pSh, "wsl", "--install")
	if err := setWSL.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func main() {
	if !checkAdmin() {
		runElevated()
	}
	if checkAdmin() {
		fmt.Println("\nDev4win v" + appVer + "\n")
		if checkNetStatus() == true {
			winBegin()
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
		endOpt:
			for {
				fmt.Printf("\nSelect command: ")
				_, err := fmt.Scanln(&cmdOpt)
				checkError(err)
				if cmdOpt == "1" {
					confG4s()
					restartWin()
				} else if cmdOpt == "2" {
					restartWin()
				} else if cmdOpt == "3" {
					confG4s()
				} else if cmdOpt == "0" || cmdOpt == "q" || cmdOpt == "e" || cmdOpt == "quit" || cmdOpt == "exit" {
				} else {
					fmt.Println("Wrong answer. Please choose between 1,2,3,0.")
					goto endOpt
				}
				break
			}
			fmt.Println("\n----------Finished!----------\n" +
				"Please RESTART your terminal and OS!\n" +
				lstDot + "Restart the terminal (CMD or PowerShell) for the changes to take effect.\n" +
				lstDot + "WSL has been setup. Restart OS for the changes to take effect.\n" +
				"\nPress 'Enter' to exit...")
			_, err := fmt.Scanln()
			checkError(err)
		} else {
			fmt.Println(lstDot + "Please check your internet connection and try again.\n")
		}

	}
}
