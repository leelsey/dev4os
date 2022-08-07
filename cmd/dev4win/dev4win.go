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
	pmsIns = "install"
	//cmdReIn = "reinstall"
	//cmdRm   = "uninstall"
	cmdYes = "-y"
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

	setExcludesFile := exec.Command(cmdGit, "config", "--global", "core.excludesfile", ignorePath)
	if err := setExcludesFile.Run(); err != nil {
		fmt.Println("error2")
		checkError(err)
	}

	fmt.Println(" " + lstDot + "Make \"gitignore_global\" file in " + ignoreDir)
}

func updateChoco() {
	updateChocolatey := exec.Command(cmdPMS, "upgrade", cmdYes, "all")
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

	chocoGit := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "git")
	chocoGit.Stderr = os.Stderr
	chocoGitLfs := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "git-lfs")
	chocoGitLfs.Stderr = os.Stderr

	if err := chocoGit.Run(); err != nil {
		checkError(err)
	}
	if err := chocoGitLfs.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func winDependency() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing dependencies for development work..."
	ldBar.FinalMSG = " - Installed dependencies!\n"
	ldBar.Start()

	chocoSSL := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "openssl")
	chocoSSL.Stderr = os.Stderr
	chocoGnuPG := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "gnupg")
	chocoGnuPG.Stderr = os.Stderr
	chococURL := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "curl")
	chococURL.Stderr = os.Stderr
	chocoWget := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "wget")
	chocoWget.Stderr = os.Stderr
	chocoGzip := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "gzip")
	chocoGzip.Stderr = os.Stderr
	chocoBzip2 := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "bzip2")
	chocoBzip2.Stderr = os.Stderr
	chocoCoreUtils := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "gnuwin32-coreutils.install")
	chocoCoreUtils.Stderr = os.Stderr
	chocoRe2C := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "re2c")
	chocoRe2C.Stderr = os.Stderr
	chocoGhostscript := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "ghostscript")
	chocoGhostscript.Stderr = os.Stderr

	if err := chocoSSL.Run(); err != nil {
		checkError(err)
	}
	if err := chocoGnuPG.Run(); err != nil {
		checkError(err)
	}
	if err := chococURL.Run(); err != nil {
		checkError(err)
	}
	if err := chocoWget.Run(); err != nil {
		checkError(err)
	}
	if err := chocoGzip.Run(); err != nil {
		checkError(err)
	}
	if err := chocoBzip2.Run(); err != nil {
		checkError(err)
	}
	if err := chocoCoreUtils.Run(); err != nil {
		checkError(err)
	}
	if err := chocoRe2C.Run(); err != nil {
		checkError(err)
	}
	if err := chocoGhostscript.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func winDevToolCLI() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing developer tools for CLI..."
	ldBar.FinalMSG = " - Installed developer utilities!\n"
	ldBar.Start()

	chocoGawk := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "gawk")
	chocoGawk.Stderr = os.Stderr
	chocoWatchman := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "watchman")
	chocoWatchman.Stderr = os.Stderr
	chocoQEMU := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "qemu")
	chocoQEMU.Stderr = os.Stderr
	chocoCcache := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "ccache")
	chocoCcache.Stderr = os.Stderr
	chocoMake := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "make")
	chocoMake.Stderr = os.Stderr
	chocoVim := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "vim")
	chocoVim.Stderr = os.Stderr
	chocoBat := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "bat")
	chocoBat.Stderr = os.Stderr
	chocoJQ := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "jq")
	chocoJQ.Stderr = os.Stderr
	chocoGH := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "gh")
	chocoGH.Stderr = os.Stderr
	chocoPS := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "powershell")
	chocoPS.Stderr = os.Stderr
	chocoCygwin := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "cygwin")
	chocoCygwin.Stderr = os.Stderr
	chocoVS := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "visualstudio2022community")
	chocoVS.Stderr = os.Stderr

	if err := chocoGawk.Run(); err != nil {
		checkError(err)
	}
	if err := chocoWatchman.Run(); err != nil {
		checkError(err)
	}
	if err := chocoQEMU.Run(); err != nil {
		checkError(err)
	}
	if err := chocoCcache.Run(); err != nil {
		checkError(err)
	}
	if err := chocoMake.Run(); err != nil {
		checkError(err)
	}
	if err := chocoVim.Run(); err != nil {
		checkError(err)
	}
	if err := chocoBat.Run(); err != nil {
		checkError(err)
	}
	if err := chocoJQ.Run(); err != nil {
		checkError(err)
	}
	if err := chocoGH.Run(); err != nil {
		checkError(err)
	}
	if err := chocoPS.Run(); err != nil {
		checkError(err)
	}
	if err := chocoCygwin.Run(); err != nil {
		checkError(err)
	}
	if err := chocoVS.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func winServer() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing developing tools for server..."
	ldBar.FinalMSG = " - Installed server and database!\n"
	ldBar.Start()

	chocoHTTPD := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "apache-httpd")
	chocoHTTPD.Stderr = os.Stderr
	chocoTomcat := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "tomcat")
	chocoTomcat.Stderr = os.Stderr
	chocoSQLite := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "sqlite")
	chocoSQLite.Stderr = os.Stderr
	chocoPostgreSQL := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "postgresql")
	chocoPostgreSQL.Stderr = os.Stderr
	chocoMySQL := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "mysql")
	chocoMySQL.Stderr = os.Stderr

	if err := chocoHTTPD.Run(); err != nil {
		checkError(err)
	}
	if err := chocoTomcat.Run(); err != nil {
		checkError(err)
	}
	if err := chocoSQLite.Run(); err != nil {
		checkError(err)
	}
	if err := chocoPostgreSQL.Run(); err != nil {
		checkError(err)
	}
	if err := chocoMySQL.Run(); err != nil {
		checkError(err)
	}
	ldBar.Stop()
}

func winLanguage() {
	ldBar := spinner.New(spinner.CharSets[43], 500*time.Millisecond)
	ldBar.Suffix = " Installing computer programming language..."
	ldBar.FinalMSG = " - Installed basic languages!\n"
	ldBar.Start()

	chocoGCC := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "mingw")
	chocoGCC.Stderr = os.Stderr
	chocoLLVM := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "llvm")
	chocoLLVM.Stderr = os.Stderr
	chocoNuGet := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "nuget.commandline")
	chocoNuGet.Stderr = os.Stderr
	chocoPerl := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "strawberryperl")
	chocoPerl.Stderr = os.Stderr
	chocoRuby := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "ruby")
	chocoRuby.Stderr = os.Stderr
	chocoPython := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "python")
	chocoPython.Stderr = os.Stderr
	chocoLua := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "lua")
	chocoLua.Stderr = os.Stderr
	chocoGo := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "go")
	chocoGo.Stderr = os.Stderr
	chocoRust := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "rust")
	chocoRust.Stderr = os.Stderr
	chocoNode := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "nodejs")
	chocoNode.Stderr = os.Stderr
	chocoPHP := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "php")
	chocoPHP.Stderr = os.Stderr
	chocoJDK := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "openjdk")
	chocoJDK.Stderr = os.Stderr
	chocoGroovy := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "groovy")
	chocoGroovy.Stderr = os.Stderr
	chocoScala := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "scala")
	chocoScala.Stderr = os.Stderr
	chocoClojure := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "clojure")
	chocoClojure.Stderr = os.Stderr
	chocoErlang := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "erlang")
	chocoErlang.Stderr = os.Stderr
	chocoElixir := exec.Command(pSh, cmdPMS, pmsIns, cmdYes, "elixir")
	chocoElixir.Stderr = os.Stderr

	if err := chocoGCC.Run(); err != nil {
		checkError(err)
	}
	if err := chocoLLVM.Run(); err != nil {
		checkError(err)
	}
	if err := chocoNuGet.Run(); err != nil {
		checkError(err)
	}
	if err := chocoPerl.Run(); err != nil {
		checkError(err)
	}
	if err := chocoRuby.Run(); err != nil {
		checkError(err)
	}
	if err := chocoPython.Run(); err != nil {
		checkError(err)
	}
	if err := chocoLua.Run(); err != nil {
		checkError(err)
	}
	if err := chocoGo.Run(); err != nil {
		checkError(err)
	}
	if err := chocoRust.Run(); err != nil {
		checkError(err)
	}
	if err := chocoNode.Run(); err != nil {
		checkError(err)
	}
	if err := chocoPHP.Run(); err != nil {
		checkError(err)
	}
	if err := chocoJDK.Run(); err != nil {
		checkError(err)
	}
	if err := chocoGroovy.Run(); err != nil {
		checkError(err)
	}
	if err := chocoScala.Run(); err != nil {
		checkError(err)
	}
	if err := chocoClojure.Run(); err != nil {
		checkError(err)
	}
	if err := chocoErlang.Run(); err != nil {
		checkError(err)
	}
	if err := chocoElixir.Run(); err != nil {
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
