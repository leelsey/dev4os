package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"golang.org/x/term"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var (
	appVer     = "0.5"
	lstDot     = " • "
	shrcPath   = homeDir() + ".zshrc"
	prfPath    = homeDir() + ".zprofile"
	arm64Path  = "/opt/homebrew/"
	amd64Path  = "/usr/local/"
	brewPrefix = checkBrewPrefix()
	cmdAdmin   = "sudo"
	cmdSh      = "/bin/bash"
	cmdPMS     = checkBrewPath()
	cmdGit     = "/usr/bin/git"
	cmdASDF    = checkASDFPath()
	optIns     = "install"
	optReIn    = "reinstall"
	//optUnIn    = "uninstall"
	//optRm      = "remove"
	optAlt    = "--cask"
	optRepo   = "tap"
	fontPath  = homeDir() + "Library/Fonts/"
	p10kPath  = homeDir() + ".config/p10k/"
	p10kCache = homeDir() + ".cache/p10k-" + userName()
	tryLoop   = 0
	clrReset  = "\033[0m"
	clrRed    = "\033[31m"
	clrGreen  = "\033[32m"
	clrYellow = "\033[33m"
	clrBlue   = "\033[34m"
	clrPurple = "\033[35m"
	clrCyan   = "\033[36m"
	clrGrey   = "\033[37m"
	runLdBar  = spinner.New(spinner.CharSets[11], 50*time.Millisecond)
	macLdBar  = spinner.New(spinner.CharSets[16], 50*time.Millisecond)
)

func messageError(handling, msg, code string) {
	errOccurred := clrRed + "\nError occurred " + clrReset + "at "
	errMsgFormat := "\n" + clrRed + "Error >> " + clrReset + msg + " (" + code + ")\n"
	if handling == "fatal" || handling == "stop" {
		fmt.Print(errors.New("\n" + lstDot + "Fatal error" + errOccurred))
		log.Fatalln(errMsgFormat)
	} else if handling == "print" || handling == "continue" {
		log.Println(errMsgFormat)
	} else if handling == "panic" || handling == "detail" {
		fmt.Print(errors.New("\n" + lstDot + "Panic error" + errOccurred))
		panic(errMsgFormat)
	} else {
		fmt.Print(errors.New("\n" + lstDot + "Unknown error" + errOccurred))
		log.Fatalln(errMsgFormat)
	}
}

func checkError(err error, msg string) {
	if err != nil {
		messageError("fatal", msg, err.Error())
	}
}

func checkCmdError(err error, msg, pkg string) {
	if err != nil {
		messageError("print", msg+" "+clrYellow+pkg+clrReset, err.Error())
	}
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

func checkExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func checkArchitecture() bool {
	switch runtime.GOARCH {
	case "arm64":
		return true
	}
	return false
}

func checkBrewPrefix() string {
	if checkArchitecture() == true {
		return arm64Path
	} else {
		return amd64Path
	}
}

func checkBrewPath() string {
	if checkArchitecture() == true {
		return arm64Path + "bin/brew"
	} else {
		return amd64Path + "bin/brew"
	}
}

func checkASDFPath() string {
	asdfPath := "opt/asdf/libexec/bin/asdf"
	if checkArchitecture() == true {
		return arm64Path + asdfPath
	} else {
		return amd64Path + asdfPath
	}
}

func checkPassword() (string, bool) {
	for tryLoop < 3 {
		fmt.Print("Password:")
		bytePw, _ := term.ReadPassword(0)

		runLdBar.Suffix = " Checking password... "
		runLdBar.Start()

		tryLoop++
		strPw := string(bytePw)
		inputPw := exec.Command("echo", strPw)
		checkPw := exec.Command(cmdAdmin, "-Sv")
		checkPw.Env = os.Environ()
		checkPw.Stdout = os.Stdout
		checkPw.Stdin, _ = inputPw.StdoutPipe()
		_ = checkPw.Start()
		_ = inputPw.Run()
		errSudo := checkPw.Wait()
		if errSudo != nil {
			runLdBar.FinalMSG = clrRed + "Password check failed" + clrReset + "\n"
			runLdBar.Stop()
			if tryLoop < 3 {
				fmt.Println(errors.New(lstDot + "Sorry, try again."))
			} else if tryLoop >= 3 {
				fmt.Println(errors.New(lstDot + "3 incorrect password attempts."))
			}
		} else {
			runLdBar.Stop()
			if tryLoop == 1 {
				clearLine(tryLoop)
			} else {
				clearLine(tryLoop * 2)
			}
			return strPw, true
		}
	}
	return "", false
}

func checkPermission(runOpt, brewStatus string) bool {
	expMsg := "Need " + clrYellow + "ROOT permission " + clrReset + "to install "

	if runOpt == "1" || runOpt == "2" {
		if brewStatus == "Install" {
			fmt.Println(expMsg + "homebrew")
			return true
		} else {
			return false
		}
	} else {
		if brewStatus == "Install" {
			fmt.Println(expMsg + "homebrew and few applications")
		} else if brewStatus == "Update" {
			fmt.Println(expMsg + "few applications")
		}
		return true
	}
}

func needPermission(strPw string) {
	inputPw := exec.Command("echo", strPw)
	checkPw := exec.Command(cmdAdmin, "-Sv")
	checkPw.Env = os.Environ()
	checkPw.Stdout = os.Stdout

	checkPw.Stdin, _ = inputPw.StdoutPipe()
	_ = checkPw.Start()
	_ = inputPw.Run()
	errSudo := checkPw.Wait()
	checkError(errSudo, "Failed to run root permission")

	runRoot := exec.Command(cmdAdmin, "whoami")
	runRoot.Env = os.Environ()
	whoAmI, _ := runRoot.Output()

	if string(whoAmI) != "root\n" {
		msg := "Incorrect user, please check permission of sudo.\n" +
			lstDot + "It need sudo command of \"" + clrRed + "root" + clrReset + "\" user's permission.\n" +
			lstDot + "Working username: " + string(whoAmI)
		messageError("fatal", msg, "User")
	}
}

func clearLine(line int) {
	for clear := 0; clear < line; clear++ {
		fmt.Printf("\033[1A\033[K")
	}
}

func netHTTP(urlPath string) string {
	resp, err := http.Get(urlPath)
	checkError(err, "Failed to connect "+urlPath)

	defer func() {
		errBodyClose := resp.Body.Close()
		checkError(errBodyClose, "Failed to download from "+urlPath)
	}()

	rawFile, err := io.ReadAll(resp.Body)
	checkError(err, "Failed to read file information from "+urlPath)
	return string(rawFile)
}

func netJSON(urlPath, key string) string {
	resp, err := http.Get(urlPath)
	checkError(err, "Failed to connect "+urlPath)

	defer func() {
		errBodyClose := resp.Body.Close()
		checkError(errBodyClose, "Failed to download from "+urlPath)
	}()

	jsonFile, err := io.ReadAll(resp.Body)
	checkError(err, "Failed to read file information from "+urlPath)

	var res map[string]interface{}
	errMarshal := json.Unmarshal(jsonFile, &res)
	checkError(errMarshal, "Failed to parse JSON file from "+urlPath)
	return res[key].(string)
}

func homeDir() string {
	homeDirPath, err := os.UserHomeDir()
	checkError(err, "Failed to get home directory")
	return homeDirPath + "/"
}

func workingDir() string {
	workingDirPath, err := os.Getwd()
	checkError(err, "Failed to get working directory")
	return workingDirPath + "/"
}

func userName() string {
	workingUser, err := user.Current()
	checkError(err, "Failed to get current user")
	return workingUser.Username
}

func systemUpdate() {
	runLdBar.Suffix = " Updating OS, please wait a moment ... "
	runLdBar.Start()

	osUpdate := exec.Command("softwareupdate", "--all", "--install", "--force")
	errOSUpdate := osUpdate.Run()
	checkError(errOSUpdate, "Failed to update Operating System")

	runLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "update OS!\n"
	runLdBar.Stop()
}

func systemReboot(adminCode string) {
	runLdBar.Suffix = " Restarting OS, please wait a moment ... "
	runLdBar.Start()

	needPermission(adminCode)
	reboot := exec.Command(cmdAdmin, "shutdown", "-r", "now")
	time.Sleep(time.Second * 3)
	if err := reboot.Run(); err != nil {
		runLdBar.FinalMSG = clrRed + "Error: " + clrReset
		runLdBar.Stop()
		fmt.Println(errors.New("failed to reboot Operating System"))
	}

	runLdBar.FinalMSG = "⣾ Restarting OS, please wait a moment ... "
	runLdBar.Stop()
}

func makeDirectory(dirPath string) {
	if checkExists(dirPath) != true {
		err := os.MkdirAll(dirPath, 0755)
		checkError(err, "Failed to make directory")
	}
}

func copyDirectory(srcPath, dstPath string) {
	if checkExists(dstPath) != true {
		cpDir := exec.Command("cp", "-rf", srcPath, dstPath)
		cpDir.Stderr = os.Stderr
		err := cpDir.Run()
		checkError(err, "Failed to copy directory from \""+srcPath+"\" to \""+dstPath+"\"")
	}
}

func makeFile(filePath, fileContents string, fileMode int) {
	targetFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(fileMode))
	checkError(err, "Failed to get file information to make new file from \""+filePath+"\"")

	defer func() {
		err := targetFile.Close()
		checkError(err, "Failed to finish make file to \""+filePath+"\"")
	}()

	_, err = targetFile.Write([]byte(fileContents))
	checkError(err, "Failed to fill in information to \""+filePath+"\"")
}

func copyFile(srcPath, dstPath string) {
	srcFile, err := os.Open(srcPath)
	checkError(err, "Failed to get file information to copy from \""+srcPath+"\"")
	dstFile, err := os.Create(dstPath)
	checkError(err, "Failed to get file information to copy to \""+dstPath+"\"")

	defer func() {
		errSrcFileClose := srcFile.Close()
		checkError(errSrcFileClose, "Failed to finish copy file from \""+srcPath+"\"")
		errDstFileClose := dstFile.Close()
		checkError(errDstFileClose, "Failed to finish copy file to \""+dstPath+"\"")
	}()

	_, errCopy := io.Copy(dstFile, srcFile)
	checkError(errCopy, "Failed to copy file from \""+srcPath+"\" to \""+dstPath+"\"")
	errSync := dstFile.Sync()
	checkError(errSync, "Failed to sync file from \""+srcPath+"\" to \""+dstPath+"\"")
}

func removeFile(filePath string) {
	if checkExists(filePath) == true {
		err := os.Remove(filePath)
		checkError(err, "Failed to remove file \""+filePath+"\"")
	}
}

func linkFile(srcPath, dstPath, linkType, permission, adminCode string) {
	if linkType == "hard" {
		if permission == "root" || permission == "sudo" || permission == "admin" {
			needPermission(adminCode)
			lnFile := exec.Command(cmdAdmin, "ln", "-sfn", srcPath, dstPath)
			lnFile.Stderr = os.Stderr
			err := lnFile.Run()
			checkCmdError(err, "Add failed to hard link file", "\""+srcPath+"\"->\""+dstPath+"\"")
		} else {
			if checkExists(srcPath) == true {
				if checkExists(dstPath) == true {
					removeFile(dstPath)
				}
				errHardlink := os.Link(srcPath, dstPath)
				checkCmdError(errHardlink, "Add failed to hard link", "\""+srcPath+"\"->\""+dstPath+"\"")
			}
		}
	} else if linkType == "symbolic" {
		if permission == "root" || permission == "sudo" || permission == "admin" {
			needPermission(adminCode)
			lnFile := exec.Command(cmdAdmin, "ln", "-sfn", srcPath, dstPath)
			lnFile.Stderr = os.Stderr
			err := lnFile.Run()
			checkCmdError(err, "Add failed to symbolic link", "\""+srcPath+"\"->\""+dstPath+"\"")
		} else {
			if checkExists(srcPath) == true {
				if checkExists(dstPath) == true {
					removeFile(dstPath)
				}
				errSymlink := os.Symlink(srcPath, dstPath)
				checkCmdError(errSymlink, "Add failed to symbolic link\"", srcPath+"\"->\""+dstPath+"\"")
				errLinkOwn := os.Lchown(dstPath, os.Getuid(), os.Getgid())
				checkError(errLinkOwn, "Failed to change ownership of symlink \""+dstPath+"\"")
			}
		}
	} else {
		messageError("fatal", "Invalid link type", "Link file")
	}
}

func appendContents(filePath, fileContents string, fileMode int) {
	targetFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.FileMode(fileMode))
	checkError(err, "Failed to get file information to append contents from \""+filePath+"\"")

	defer func() {
		err := targetFile.Close()
		checkError(err, "Failed to finish append contents to \""+filePath+"\"")
	}()

	_, err = targetFile.Write([]byte(fileContents))
	checkError(err, "Failed to append contents to \""+filePath+"\"")
}

func downloadFile(filePath, urlPath string, fileMode int) {
	makeFile(filePath, netHTTP(urlPath), fileMode)
}

func runApplication(appName string) {
	runApp := exec.Command("open", "/Applications/"+appName+".app")
	err := runApp.Run()
	checkCmdError(err, "ASDF-VM failed to add", appName)
}

func changeAppIcon(appName, icnName, adminCode string) {
	srcIcn := workingDir() + ".dev4mac-app-icn.icns"
	downloadFile(srcIcn, "https://raw.githubusercontent.com/leelsey/ConfStore/main/icns/"+icnName, 0755)

	appSrc := strings.Replace(appName, " ", "\\ ", -1)
	appPath := "/Applications/" + appSrc + ".app"
	chicnPath := workingDir() + ".dev4mac-chicn.sh"
	cvtIcn := workingDir() + ".dev4mac-app-icn.rsrc"
	chIcnSrc := "sudo rm -rf \"" + appPath + "\"$'/Icon\\r'\n" +
		"sips -i " + srcIcn + " > /dev/null\n" +
		"DeRez -only icns " + srcIcn + " > " + cvtIcn + "\n" +
		"sudo Rez -append " + cvtIcn + " -o " + appPath + "$'/Icon\\r'\n" +
		"sudo SetFile -a C " + appPath + "\n" +
		"sudo SetFile -a V " + appPath + "$'/Icon\\r'"
	makeFile(chicnPath, chIcnSrc, 0644)

	needPermission(adminCode)
	chicn := exec.Command(cmdSh, chicnPath)
	chicn.Env = os.Environ()
	chicn.Stderr = os.Stderr
	err := chicn.Run()
	checkCmdError(err, "Failed change icon of", appName+".app")

	removeFile(srcIcn)
	removeFile(cvtIcn)
	removeFile(chicnPath)
}

func brewUpdate() {
	updateHomebrew := exec.Command(cmdPMS, "update", "--auto-update")
	err := updateHomebrew.Run()
	checkCmdError(err, "Brew failed to", "update repositories")
}

func brewUpgrade() {
	brewUpdate()
	upgradeHomebrew := exec.Command(cmdPMS, "upgrade", "--greedy")
	err := upgradeHomebrew.Run()
	checkCmdError(err, "Brew failed to", "upgrade packages")
}

func brewRepository(repo string) {
	brewRepo := strings.Split(repo, "/")
	repoPath := strings.Join(brewRepo[0:1], "") + "/homebrew-" + strings.Join(brewRepo[1:2], "")
	if checkExists(brewPrefix+"Homebrew/Library/Taps/"+repoPath) != true {
		brewRepo := exec.Command(cmdPMS, optRepo, repo)
		err := brewRepo.Run()
		checkCmdError(err, "Brew failed to add ", repo)
	}
}

func brewCleanup() {
	upgradeHomebrew := exec.Command(cmdPMS, "cleanup", "--prune=all", "-nsd")
	err := upgradeHomebrew.Run()
	checkCmdError(err, "Brew failed to", "cleanup old packages")
}

func brewRemoveCache() {
	upgradeHomebrew := exec.Command("rm", "-rf", "\"$(brew --cache)\"")
	err := upgradeHomebrew.Run()
	checkCmdError(err, "Brew failed to", "remove cache")
}

func brewInstall(pkg string) {
	if checkExists(brewPrefix+"Cellar/"+pkg) != true {
		brewUpdate()
		brewIns := exec.Command(cmdPMS, optIns, pkg)
		brewIns.Stderr = os.Stderr
		err := brewIns.Run()
		checkCmdError(err, "Brew failed to install", pkg)
	}
}

func brewInstallQuiet(pkg string) {
	if checkExists(brewPrefix+"Cellar/"+pkg) != true {
		brewUpdate()
		brewIns := exec.Command(cmdPMS, optIns, "--quiet", pkg)
		err := brewIns.Run()
		checkCmdError(err, "Brew failed to install", pkg)
	}
}

func brewInstallCask(pkg, appName string) {
	if checkExists(brewPrefix+"Caskroom/"+pkg) != true {
		brewUpdate()
		if checkExists("/Applications/"+appName+".app") != true {
			brewIns := exec.Command(cmdPMS, optIns, optAlt, pkg)
			err := brewIns.Run()
			checkCmdError(err, "Brew failed to install cask", pkg)
		} else {
			brewIns := exec.Command(cmdPMS, optReIn, optAlt, pkg)
			err := brewIns.Run()
			checkCmdError(err, "Brew failed to reinstall cask", pkg)
		}
	}
}

func brewInstallCaskSudo(pkg, appName, appPath, adminCode string) {
	if checkExists(brewPrefix+"Caskroom/"+pkg) != true {
		brewUpdate()
		needPermission(adminCode)
		if checkExists(appPath) != true {
			brewIns := exec.Command(cmdPMS, optIns, optAlt, pkg)
			err := brewIns.Run()
			checkCmdError(err, "Brew failed to install cask", appName)
		} else {
			brewIns := exec.Command(cmdPMS, optReIn, optAlt, pkg)
			err := brewIns.Run()
			checkCmdError(err, "Brew failed to install cask", appName)
		}
	}
}

func asdfInstall(plugin, version string) {
	if checkExists(homeDir()+".asdf/plugins/"+plugin) != true {
		asdfPlugin := exec.Command(cmdASDF, "plugin", "add", plugin)
		err := asdfPlugin.Run()
		checkCmdError(err, "ASDF-VM failed to add", plugin)
	}

	asdfReshim()
	asdfIns := exec.Command(cmdASDF, optIns, plugin, version)
	asdfIns.Env = os.Environ()
	errIns := asdfIns.Run()
	checkCmdError(errIns, "ASDF-VM", plugin)

	asdfGlobal := exec.Command(cmdASDF, "global", plugin, version)
	asdfGlobal.Env = os.Environ()
	errConf := asdfGlobal.Run()
	checkCmdError(errConf, "ASDF-VM failed to install", plugin)
}

func asdfReshim() {
	reshim := exec.Command(cmdASDF, "reshim")
	err := reshim.Run()
	checkCmdError(err, "ASDF failed to", "reshim")
}

func addJavaHome(srcVer, dstVer, adminCode string) {
	if checkExists(brewPrefix+"Cellar/openjdk"+srcVer) == true {
		linkFile(brewPrefix+"opt/openjdk"+srcVer+" /libexec/openjdk.jdk", "/Library/Java/JavaVirtualMachines/openjdk"+dstVer+".jdk", "symbolic", "root", adminCode)
	}
}

func confA4s() {
	a4sPath := homeDir() + ".config/alias4sh"
	makeDirectory(a4sPath)
	makeFile(a4sPath+"/alias4.sh", "# ALIAS4SH", 0644)

	dlA4sPath := workingDir() + ".dev4mac-alias4sh.sh"
	downloadFile(dlA4sPath, "https://raw.githubusercontent.com/leelsey/Alias4sh/main/install.sh", 0644)

	installA4s := exec.Command("sh", dlA4sPath)
	if err := installA4s.Run(); err != nil {
		removeFile(dlA4sPath)
		checkError(err, "Failed to install Alias4sh")
	}

	removeFile(dlA4sPath)
}

func confG4s() {
	fmt.Println(clrCyan + "Git global configuration" + clrReset)

	fmt.Println(lstDot + "Add user information")
	consoleReader := bufio.NewScanner(os.Stdin)
	fmt.Print("  - User name: ")
	consoleReader.Scan()
	gitUserName := consoleReader.Text()
	fmt.Print("  - User email: ")
	consoleReader.Scan()
	gitUserEmail := consoleReader.Text()

	setGitUserName := exec.Command(cmdGit, "config", "--global", "user.name", gitUserName)
	errGitUserName := setGitUserName.Run()
	checkError(errGitUserName, "Failed to set git user name")
	setGitUserEmail := exec.Command(cmdGit, "config", "--global", "user.email", gitUserEmail)
	errGitUserEmail := setGitUserEmail.Run()
	checkError(errGitUserEmail, "Failed to set git user email")
	clearLine(3)
	fmt.Println(lstDot + "Saved user name(" + gitUserName + ") and email(" + gitUserEmail + ").")

	setGitBranch := exec.Command(cmdGit, "config", "--global", "init.defaultBranch", "main")
	errGitBranch := setGitBranch.Run()
	checkError(errGitBranch, "Failed to change branch default name (master -> main)")
	fmt.Println(lstDot + "Main git branch default name changed master -> main.")

	setGitColor := exec.Command(cmdGit, "config", "--global", "color.ui", "true")
	errGitColor := setGitColor.Run()
	checkError(errGitColor, "Failed to setup colourising")
	fmt.Println(lstDot + "Colourising enabled.")

	setGitEditor := exec.Command(cmdGit, "config", "--global", "core.editor", "vi")
	errGitEditor := setGitEditor.Run()
	checkError(errGitEditor, "Failed to setup editor vi (vim)")
	fmt.Println(lstDot + "Default editor set to vi (vim).")

	ignoreDirPath := homeDir() + ".config/git/"
	ignorePath := ignoreDirPath + "gitignore_global"
	makeDirectory(ignoreDirPath)
	downloadFile(ignorePath, "https://raw.githubusercontent.com/leelsey/Git4set/main/gitignore-sample", 0644)
	setExcludesFile := exec.Command(cmdGit, "config", "--global", "core.excludesfile", ignorePath)
	errExcludesFile := setExcludesFile.Run()
	checkError(errExcludesFile, "Failed to set git global ignore file")
	fmt.Println(lstDot + "Ignore list set in \"" + ignoreDirPath + "gitignore_global\".")
}

func installBrew(adminCode string) {
	insBrewPath := workingDir() + ".dev4mac-brew.sh"
	downloadFile(insBrewPath, "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh", 0755)

	needPermission(adminCode)
	installHomebrew := exec.Command(cmdSh, "-c", insBrewPath)
	installHomebrew.Env = append(os.Environ(), "NONINTERACTIVE=1")
	if err := installHomebrew.Run(); err != nil {
		removeFile(insBrewPath)
		checkError(err, "Failed to install Homebrew")
	}
	removeFile(insBrewPath)

	if checkExists(cmdPMS) == false {
		messageError("fatal", "Installed brew failed, please check your system", "Can't find Homebrew")
	}
}

func installXAMPP(adminCode string) {
	xamppVer := netJSON("https://formulae.brew.sh/api/cask/xampp-vm.json", "version")
	xamppName := "xampp-osx-" + xamppVer + "-vm"
	brewInstallCaskSudo("xampp-vm", xamppName, "/Applications/"+xamppName+".app", adminCode)
	changeAppIcon("xampp-osx-"+xamppVer+"-vm", "XAMPP.icns", adminCode)
}

func installHopper(adminCode string) {
	hopperRSS := strings.Split(netHTTP("https://www.hopperapp.com/rss/html_changelog.php"), " ")
	hopperVer := strings.Join(hopperRSS[1:2], "")

	dlHopperPath := workingDir() + ".Hopper.dmg"
	downloadFile(dlHopperPath, "https://d2ap6ypl1xbe4k.cloudfront.net/Hopper-"+hopperVer+"-demo.dmg", 0755)

	mountHopper := exec.Command("hdiutil", "attach", dlHopperPath)
	errMount := mountHopper.Run()
	checkError(errMount, "Failed to mount "+clrYellow+"Hopper.dmg"+clrReset)
	removeFile(dlHopperPath)

	appName := "Hopper Disassembler v4"
	copyDirectory("/Volumes/Hopper Disassembler/"+appName+".app", "/Applications/"+appName+".app")

	unmountDmg := exec.Command("hdiutil", "unmount", "/Volumes/Hopper Disassembler")
	errUnmount := unmountDmg.Run()
	checkError(errUnmount, "Failed to unmount "+clrYellow+"Hopper Disassembler"+clrReset)

	if checkArchitecture() == true {
		changeAppIcon(appName, "Hopper Disassembler ARM64.icns", adminCode)
	} else {
		changeAppIcon(appName, "Hopper Disassembler AMD64.icns", adminCode)
	}
}

func macBegin(adminCode string) {
	if checkExists(cmdPMS) == true {
		macLdBar.Suffix = " Updating homebrew... "
		macLdBar.Start()
		macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "update homebrew!\n"
	} else {
		macLdBar.Suffix = " Installing homebrew... "
		macLdBar.Start()
		macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install and update homebrew!\n"

		installBrew(adminCode)
	}
	err := os.Chmod(brewPrefix+"share", 0755)
	checkError(err, "Failed to change permissions on "+brewPrefix+"share to 755")

	brewUpdate()
	brewRepository("homebrew/core")
	brewRepository("homebrew/cask")
	brewRepository("homebrew/cask-versions")
	brewUpgrade()

	macLdBar.Stop()
}

func macEnv() {
	macLdBar.Suffix = " Setting basic environment... "
	macLdBar.Start()

	if checkExists(prfPath) == true {
		copyFile(prfPath, homeDir()+".zprofile.bck")
	}
	if checkExists(shrcPath) == true {
		copyFile(shrcPath, homeDir()+".zshrc.bck")
	}

	profileContents := "#    ___________  _____   ____  ______ _____ _      ______ \n" +
		"#   |___  /  __ \\|  __ \\ / __ \\|  ____|_   _| |    |  ____|\n" +
		"#      / /| |__) | |__) | |  | | |__    | | | |    | |__   \n" +
		"#     / / |  ___/|  _  /| |  | |  __|   | | | |    |  __|  \n" +
		"#    / /__| |    | | \\ \\| |__| | |     _| |_| |____| |____ \n" +
		"#   /_____|_|    |_|  \\_\\\\____/|_|    |_____|______|______|\n#\n" +
		"#  " + userName() + "’s zsh profile\n\n" +
		"# HOMEBREW\n" +
		"eval \"$(" + cmdPMS + " shellenv)\"\n\n"
	makeFile(prfPath, profileContents, 0644)

	shrcContents := "#   ______ _____ _    _ _____   _____\n" +
		"#  |___  // ____| |  | |  __ \\ / ____|\n" +
		"#     / /| (___ | |__| | |__) | |\n" +
		"#    / /  \\___ \\|  __  |  _  /| |\n" +
		"#   / /__ ____) | |  | | | \\ \\| |____\n" +
		"#  /_____|_____/|_|  |_|_|  \\_\\\\_____|\n#\n" +
		"#  " + userName() + "’s zsh run commands\n\n"
	makeFile(shrcPath, shrcContents, 0644)

	makeDirectory(homeDir() + ".config")
	makeDirectory(homeDir() + ".cache")

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "setup zsh environment!\n"
	macLdBar.Stop()
}

func macDependency(runOpt string) {
	macLdBar.Suffix = " Installing dependencies... "
	macLdBar.Start()

	brewInstall("pkg-config")
	brewInstall("ca-certificates")
	brewInstall("ncurses")
	brewInstall("openssl@3")
	brewInstall("openssl@1.1")
	brewInstall("readline")
	brewInstall("autoconf")
	brewInstall("automake")
	brewInstall("mpdecimal")
	brewInstall("utf8proc")
	brewInstall("m4")
	brewInstall("gmp")
	brewInstall("mpfr")
	brewInstall("gettext")
	brewInstall("jpeg-turbo")
	brewInstall("libtool")
	brewInstall("libevent")
	brewInstall("libffi")
	brewInstall("libtiff")
	brewInstall("libvmaf")
	brewInstall("libpng")
	brewInstall("libyaml")
	brewInstall("giflib")
	brewInstall("xz")
	brewInstall("gdbm")
	brewInstall("sqlite")
	brewInstall("lz4")
	brewInstall("zstd")
	brewInstall("hiredis")
	brewInstall("berkeley-db")
	brewInstall("asciidoctor")
	brewInstall("freetype")
	brewInstall("fontconfig")
	brewInstall("pcre")
	brewInstall("pcre2")

	shrcAppend := "# NCURSES\n" +
		"export PATH=\"" + brewPrefix + "opt/ncurses/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/ncurses/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/ncurses/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/ncurses/lib/pkgconfig\"\n\n" +
		"# OPENSSL-3\n" +
		"export PATH=\"" + brewPrefix + "opt/openssl@3/bin:$PATH\"\n" +
		"export LDFLAGS=\"-L" + brewPrefix + "opt/openssl@3/lib\"\n" +
		"export CPPFLAGS=\"-I" + brewPrefix + "opt/openssl@3/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/openssl@3/lib/pkgconfig\"\n\n" +
		"# OPENSSL-1.1\n" +
		"export PATH=\"" + brewPrefix + "opt/openssl@1.1/bin:$PATH\"\n" +
		"export LDFLAGS=\"-L" + brewPrefix + "opt/openssl@1.1/lib\"\n" +
		"export CPPFLAGS=\"-I" + brewPrefix + "opt/openssl@1.1/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/openssl@1.1/lib/pkgconfig\"\n\n"
	appendContents(shrcPath, shrcAppend, 0644)

	if runOpt != "2" && runOpt != "3" {
		brewInstall("ccache")
		brewInstall("gawk")
		brewInstall("tcl-tk")
		brewInstall("bash")
		brewInstall("zsh")
		brewInstall("perl")
		brewInstall("ruby")
		brewInstall("python@3.10")
		brewInstall("openjdk")
		brewInstall("ghc")
	}

	if runOpt == "6" || runOpt == "7" {
		brewInstall("krb5")
		brewInstall("libsodium")
		brewInstall("nettle")
		brewInstall("coreutils")
		brewInstall("gnu-getopt")
		brewInstall("ldns")
		brewInstall("isl")
		brewInstall("npth")
		brewInstall("gzip")
		brewInstall("bzip2")
		brewInstall("fop")
		brewInstall("little-cms2")
		brewInstall("imath")
		brewInstall("openldap")
		brewInstall("openexr")
		brewInstall("openjpeg")
		brewInstall("jpeg-xl")
		brewInstall("webp")
		brewInstall("rtmpdump")
		brewInstall("aom")
		brewInstall("screenresolution")
		brewInstall("brotli")
		brewInstall("bison")
		brewInstall("swig")
		brewInstall("re2c")
		brewInstall("icu4c")
		brewInstall("bdw-gc")
		brewInstall("guile")
		brewInstall("wxwidgets")
		brewInstall("sphinx-doc")
		brewInstall("docbook")
		brewInstall("docbook2x")
		brewInstall("docbook-xsl")
		brewInstall("xmlto")
		brewInstall("html-xml-utils")
		brewInstall("shared-mime-info")
		brewInstall("x265")
		brewInstall("oniguruma")
		brewInstall("libgpg-error")
		brewInstall("libgcrypt")
		brewInstall("libunistring")
		brewInstall("libatomic_ops")
		brewInstall("libiconv")
		brewInstall("libmpc")
		brewInstall("libidn")
		brewInstall("libidn2")
		brewInstall("libssh2")
		brewInstall("libnghttp2")
		brewInstall("libxml2")
		brewInstall("libtasn1")
		brewInstall("libxslt")
		brewInstall("libavif")
		brewInstall("libzip")
		brewInstall("libde265")
		brewInstall("libheif")
		brewInstall("libksba")
		brewInstall("libusb")
		brewInstall("liblqr")
		brewInstall("libomp")
		brewInstall("libassuan")
		brewInstall("p11-kit")
		brewInstall("gnutls")
		brewInstall("gd")
		brewInstall("ghostscript")
		brewInstall("imagemagick")
		brewInstall("pinentry")
		brewInstall("gnupg")
		brewInstall("curl")
		brewInstall("wget")
		brewInstall("glib")
		brewInstall("zlib")

		shrcAppend := "# KRB5\n" +
			"export PATH=\"" + brewPrefix + "opt/krb5/bin:$PATH\"\n" +
			"export PATH=\"" + brewPrefix + "opt/krb5/sbin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/krb5/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/krb5/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/krb5/lib/pkgconfig\"\n\n" +
			"# COREUTILS\n" +
			"#export PATH=\"" + brewPrefix + "opt/coreutils/libexec/gnubin:$PATH\"\n\n" +
			"# GNU GETOPT\n" +
			"export PATH=\"" + brewPrefix + "opt/gnu-getopt/bin:$PATH\"\n\n" +
			"# TCL-TK\n" +
			"export PATH=\"" + brewPrefix + "opt/tcl-tk/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/tcl-tk/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/tcl-tk/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/tcl-tk/lib/pkgconfig\"\n\n" +
			"# BZIP2\n" +
			"export PATH=\"" + brewPrefix + "opt/bzip2/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/bzip2/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/bzip2/include\"\n\n" +
			"# BISON\n" +
			"export PATH=\"" + brewPrefix + "opt/bison/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/bison/lib\"\n\n" +
			"# ICU4C\n" +
			"export PATH=\"" + brewPrefix + "opt/icu4c/bin:$PATH\"\n" +
			"export PATH=\"" + brewPrefix + "opt/icu4c/sbin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/icu4c/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/icu4c/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/icu4c/lib/pkgconfig\"\n\n" +
			"# DOCBOOK\n" +
			"export XML_CATALOG_FILES=\"" + brewPrefix + "etc/xml/catalog\"\n\n" +
			"# LIBICONV\n" +
			"export PATH=\"" + brewPrefix + "opt/libiconv/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/libiconv/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/libiconv/include\"\n\n" +
			"# LIBXML2\n" +
			"export PATH=\"" + brewPrefix + "opt/libxml2/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/libxml2/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/libxml2/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/libxml2/lib/pkgconfig\"\n\n" +
			"# LIBXSLT\n" +
			"export PATH=\"" + brewPrefix + "opt/libxslt/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/libxslt/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/libxslt/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/libxslt/lib/pkgconfig\"\n\n" +
			"# CURL\n" +
			"export PATH=\"" + brewPrefix + "opt/curl/bin:$PATH\"\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/curl/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/curl/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/curl/lib/pkgconfig\"\n\n" +
			"# ZLIB\n" +
			"export LDFLAGS=\"" + brewPrefix + "opt/zlib/lib\"\n" +
			"export CPPFLAGS=\"" + brewPrefix + "opt/zlib/include\"\n" +
			"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/zlib/lib/pkgconfig\"\n\n"
		appendContents(shrcPath, shrcAppend, 0644)
	}

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install dependencies!\n"
	macLdBar.Stop()
}

func macTerminal(runOpt string) {
	macLdBar.Suffix = " Installing zsh with useful tools... "
	macLdBar.Start()

	confA4s()
	brewInstall("zsh-completions")
	brewInstall("zsh-syntax-highlighting")
	brewInstall("zsh-autosuggestions")
	brewInstall("z")
	brewInstall("tree")

	makeFile(homeDir()+".z", "", 0644)
	makeDirectory(p10kPath)
	makeDirectory(p10kCache)

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInstall("fzf")
		brewInstall("tmux")
		brewInstall("tmuxinator")
		brewInstall("neofetch")

		dliTerm2Conf := homeDir() + "Library/Preferences/com.googlecode.iterm2.plist"
		downloadFile(dliTerm2Conf, "https://raw.githubusercontent.com/leelsey/ConfStore/main/iterm2/iTerm2.plist", 0644)
	}

	brewRepository("romkatv/powerlevel10k")
	brewInstall("romkatv/powerlevel10k/powerlevel10k")
	downloadFile(p10kPath+"p10k-term.zsh", "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-minimalism.zsh", 0644)

	if runOpt == "2" || runOpt == "3" || runOpt == "4" {
		profileAppend := "# POWERLEVEL10K\n" +
			"source " + brewPrefix + "opt/powerlevel10k/powerlevel10k.zsh-theme\n" +
			"if [[ -r \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\" ]]; then\n" +
			"  source \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\"\n" +
			"fi\n" +
			"[[ ! -f " + p10kPath + "p10k-terminal.zsh ]] || source " + p10kPath + "p10k-terminal.zsh\n\n"
		appendContents(prfPath, profileAppend, 0644)
	} else if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		downloadFile(p10kPath+"p10k-iterm2.zsh", "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-atelier.zsh", 0644)
		downloadFile(p10kPath+"p10k-tmux.zsh", "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-seeking.zsh", 0644)
		downloadFile(p10kPath+"p10k-ops.zsh", "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-operations.zsh", 0644)
		downloadFile(p10kPath+"p10k-etc.zsh", "https://raw.githubusercontent.com/leelsey/ConfStore/main/p10k/p10k-engineering.zsh", 0644)
		downloadFile(fontPath+"MesloLGS NF Bold Italic.ttf", "https://raw.githubusercontent.com/romkatv/dotfiles-public/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Bold%20Italic.ttf", 0644)
		downloadFile(fontPath+"MesloLGS NF Bold.ttf", "https://raw.githubusercontent.com/romkatv/dotfiles-public/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Bold.ttf", 0644)
		downloadFile(fontPath+"MesloLGS NF Italic.ttf", "https://raw.githubusercontent.com/romkatv/dotfiles-public/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Italic.ttf", 0644)
		downloadFile(fontPath+"MesloLGS NF Regular.ttf", "https://raw.githubusercontent.com/romkatv/dotfiles-public/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Regular.ttf", 0644)

		profileAppend := "# ZSH\n" +
			"export SHELL=zsh\n\n" +
			"# POWERLEVEL10K\n" +
			"source " + brewPrefix + "opt/powerlevel10k/powerlevel10k.zsh-theme\n" +
			"if [[ -r \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\" ]]; then\n" +
			"  source \"${XDG_CACHE_HOME:-" + p10kCache + "}/p10k-instant-prompt-${(%):-%n}.zsh\"\n" +
			"fi\n" +
			"if [[ -d /Applications/iTerm.app ]]; then\n" +
			"  if [[ $TERM_PROGRAM = \"Apple_Terminal\" ]]; then\n" +
			"    [[ ! -f " + p10kPath + "p10k-term.zsh ]] || source " + p10kPath + "p10k-term.zsh\n" +
			"  elif [[ $TERM_PROGRAM = \"iTerm.app\" ]]; then\n" +
			"    [[ ! -f " + p10kPath + "p10k-iterm2.zsh ]] || source " + p10kPath + "p10k-iterm2.zsh\n" +
			"    echo ''; neofetch --bold off\n" +
			"  elif [[ $TERM_PROGRAM = \"tmux\" ]]; then\n" +
			"    [[ ! -f " + p10kPath + "p10k-tmux.zsh ]] || source " + p10kPath + "p10k-tmux.zsh\n" +
			"    echo ''; neofetch --bold off\n" +
			"  else\n" +
			"    [[ ! -f " + p10kPath + "p10k-etc.zsh ]] || source " + p10kPath + "p10k-etc.zsh\n" +
			"  fi\n" +
			"else\n" +
			"  [[ ! -f " + p10kPath + "p10k-term.zsh ]] || source " + p10kPath + "p10k-term.zsh\n" +
			"fi\n\n"
		appendContents(prfPath, profileAppend, 0644)
	}

	profileAppend := "# ZSH-COMPLETIONS\n" +
		"if type brew &>/dev/null; then\n" +
		"  FPATH=" + brewPrefix + "share/zsh-completions:$FPATH\n" +
		"  autoload -Uz compinit\n" +
		"  compinit\n" +
		"fi\n\n" +
		"# ZSH SYNTAX HIGHLIGHTING\n" +
		"source " + brewPrefix + "share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh\n\n" +
		"# ZSH AUTOSUGGESTIONS\n" +
		"source " + brewPrefix + "share/zsh-autosuggestions/zsh-autosuggestions.zsh\n\n" +
		"# Z\n" +
		"source " + brewPrefix + "etc/profile.d/z.sh\n\n" +
		"# ALIAS4SH\n" +
		"source " + homeDir() + "/.config/alias4sh/alias4.sh\n\n" +
		"# Edit\n" +
		"export EDITOR=/usr/bin/vi\n" +
		"edit () { $EDITOR \"$@\" }\n" +
		"#vi () { $EDITOR \"$@\" }\n\n"
	appendContents(prfPath, profileAppend, 0644)

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install and configure for terminal!\n"
	macLdBar.Stop()
}

func macLanguage(runOpt, adminCode string) {
	macLdBar.Suffix = " Installing computer programming language... "
	macLdBar.Start()

	shrcAppend := "# CCACHE\n" +
		"export PATH=\"" + brewPrefix + "opt/ccache/libexec:$PATH\"\n\n" +
		"# RUBY\n" +
		"export PATH=\"" + brewPrefix + "opt/ruby/bin:$PATH\"\n"
	appendContents(shrcPath, shrcAppend, 0644)

	if runOpt == "4" || runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInstall("php")
		if checkArchitecture() == false {
			brewInstall("openjdk@8")
		}
		brewInstall("openjdk@11")
		brewInstall("openjdk@17")
		addJavaHome("", "", adminCode)
		addJavaHome("@17", "-17", adminCode)
		addJavaHome("@11", "-11", adminCode)
		if checkArchitecture() == false {
			addJavaHome("@8", "-8", adminCode)
		}
	}

	if runOpt == "3" || runOpt == "4" || runOpt == "5" {
		brewInstall("nvm")
		brewInstall("pyenv")
		brewInstall("pyenv-virtualenv")

		shrcAppend := "# NVM\n" +
			"export NVM_DIR=\"$HOME/.nvm\"\n" +
			"[ -s \"" + brewPrefix + "opt/nvm/nvm.sh\" ] && source \"" + brewPrefix + "opt/nvm/nvm.sh\"\n" +
			"[ -s \"" + brewPrefix + "opt/nvm/etc/bash_completion.d/nvm\" ] && source \"" +
			brewPrefix + "opt/nvm/etc/bash_completion.d/nvm\"\n\n" +
			"# PYENV" +
			"export PYENV_ROOT=\"$HOME/.pyenv\"\n" +
			"export PATH=\"$PYENV_ROOT/bin:$PATH\"\n" +
			"eval \"$(pyenv init --path)\"\n" +
			"eval \"$(pyenv init -)\"\n\n"
		appendContents(shrcPath, shrcAppend, 0644)

		//nvmIns := exec.Command("nvm", optIns, "--lts")
		//nvmIns.Stderr = os.Stderr
		//err := nvmIns.Run()
		//checkCmdError(err, "NVM failed to install", "LTS")
	} else if runOpt == "6" || runOpt == "7" {
		brewInstall("llvm")
		brewInstall("gcc") // fortran
		brewInstall("go")
		brewInstall("rust")
		brewInstall("lua")
		brewInstall("node")
		brewRepository("dart-lang/dart")
		brewInstall("dart")
		brewInstall("groovy")
		brewInstall("kotlin")
		brewInstall("scala")
		brewInstall("clojure")
		brewInstall("erlang")
		brewInstall("elixir")
		brewInstall("typescript")
		brewInstall("cabal-install")
		brewInstall("haskell-language-server")
		brewInstall("stylish-haskell")
	}

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install languages!\n"
	macLdBar.Stop()
}

func macServer() {
	macLdBar.Suffix = " Installing developing tools for server... "
	macLdBar.Start()

	brewInstall("httpd")
	brewInstall("tomcat")
	brewInstall("nginx")

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install servers!\n"
	macLdBar.Stop()
}

func macDatabase() {
	macLdBar.Suffix = " Installing developing tools for database... "
	macLdBar.Start()

	shrcAppend := "# SQLITE3\n" +
		"export PATH=\"" + brewPrefix + "opt/sqlite/bin:$PATH\"\n" +
		"export LDFLAGS=\"" + brewPrefix + "opt/sqlite/lib\"\n" +
		"export CPPFLAGS=\"" + brewPrefix + "opt/sqlite/include\"\n" +
		"export PKG_CONFIG_PATH=\"" + brewPrefix + "opt/sqlite/lib/pkgconfig\"\n\n"
	appendContents(shrcPath, shrcAppend, 0644)

	brewInstall("sqlite-analyzer")
	brewInstall("postgresql")
	brewInstall("mysql")
	brewInstall("redis")
	brewRepository("mongodb/brew")
	brewInstall("mongodb-community")

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install databases!\n"
	macLdBar.Stop()
}

func macDevVM() {
	macLdBar.Suffix = " Installing developer tools version management tool with plugin... "
	macLdBar.Start()

	brewInstall("asdf")

	shrcAppend := "# ASDF VM\n" +
		"source " + brewPrefix + "opt/asdf/libexec/asdf.sh\n" +
		"export RUBY_CONFIGURE_OPTS=\"--with-openssl-dir=$(brew --prefix openssl@1.1)\"\n"
	appendContents(shrcPath, shrcAppend, 0644)

	asdfrcContents := "#              _____ _____  ______  __      ____  __ \n" +
		"#       /\\    / ____|  __ \\|  ____| \\ \\    / /  \\/  |\n" +
		"#      /  \\  | (___ | |  | | |__ ____\\ \\  / /| \\  / |\n" +
		"#     / /\\ \\  \\___ \\| |  | |  __|_____\\ \\/ / | |\\/| |\n" +
		"#    / ____ \\ ____) | |__| | |         \\  /  | |  | |\n" +
		"#   /_/    \\_\\_____/|_____/|_|          \\/   |_|  |_|\n#\n" +
		"#  " + userName() + "’s ASDF-VM run commands\n\n" +
		"legacy_version_file = yes\n" +
		"use_release_candidates = no\n" +
		"always_keep_download = no\n" +
		"plugin_repository_last_check_duration = 0\n" +
		"disable_plugin_short_name_repository = no\n" +
		"java_macos_integration_enable = yes\n"
	makeFile(homeDir()+".asdfrc", asdfrcContents, 0644)

	asdfInstall("perl", "latest")
	asdfInstall("ruby", "latest")
	asdfInstall("python", "latest")
	asdfInstall("java", "openjdk-11.0.2") // JDK LTS 11
	asdfInstall("java", "openjdk-17.0.2") // JDK LTS 17
	asdfInstall("rust", "latest")
	asdfInstall("golang", "latest")
	asdfInstall("lua", "latest")
	asdfInstall("nodejs", "latest")
	asdfInstall("dart", "latest")
	asdfInstall("php", "latest")
	asdfInstall("groovy", "latest")
	asdfInstall("kotlin", "latest")
	asdfInstall("scala", "latest")
	asdfInstall("clojure", "latest")
	asdfInstall("erlang", "latest")
	asdfInstall("elixir", "latest")
	asdfInstall("gleam", "latest")
	asdfInstall("haskell", "latest")
	asdfReshim()

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install ASDF-VM with languages!\n"
	macLdBar.Stop()
}

func macCLIApp(runOpt string) {
	macLdBar.Suffix = " Installing CLI applications... "
	macLdBar.Start()

	brewInstall("unzip")
	brewInstall("diffutils")
	brewInstall("transmission-cli")

	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInstall("openssh")
		brewInstall("mosh")
		brewInstall("inetutils")
		brewInstall("git")
		brewInstall("git-lfs")
		brewInstall("gh")
		brewInstall("hub")
		brewInstall("tig")
		brewInstall("exa")
		brewInstall("bat")
		brewInstall("diffr")
		brewInstall("tldr")
		brewInstall("watchman")
		brewInstall("direnv")

		shrcAppend := "# DIRENV\n" +
			"eval \"$(direnv hook zsh)\"\n\n"
		appendContents(shrcPath, shrcAppend, 0644)
	}

	if runOpt == "6" || runOpt == "7" {
		brewInstall("make")
		brewInstallQuiet("cmake")
		brewInstall("ninja")
		brewInstall("maven")
		brewInstall("gradle")
		brewInstall("rustup-init")
		brewInstall("htop")
		brewInstall("qemu")
		brewInstall("vim")
		brewInstall("neovim")
		brewInstall("curlie")
		brewInstall("jq")
		brewInstall("yq")
		brewInstall("dasel")
		brewInstall("asciinema")
	}

	if runOpt == "7" {
		brewInstall("tor")
		brewInstall("torsocks")
		brewInstall("nmap")
		brewInstall("radare2")
		brewInstall("sleuthkit")
		brewInstall("autopsy")
		brewInstall("virustotal-cli")
	}

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install CLI applications!\n"
	macLdBar.Stop()
}

func macGUIApp(runOpt, adminCode string) {
	macLdBar.Suffix = " Installing GUI applications... "
	macLdBar.Start()

	if runOpt != "7" {
		brewInstallCask("appcleaner", "AppCleaner")
		changeAppIcon("AppCleaner", "AppCleaner.icns", adminCode)
	} else if runOpt == "7" {
		brewInstallCask("sensei", "Sensei")
	}

	brewInstallCask("keka", "Keka")
	brewInstallCask("iina", "IINA")
	brewInstallCask("transmission", "Transmission")
	changeAppIcon("Transmission", "Transmission.icns", adminCode)
	brewInstallCask("rectangle", "Rectangle")
	brewInstallCask("google-chrome", "Google Chrome")
	brewInstallCask("firefox", "Firefox")
	changeAppIcon("Firefox", "Firefox.icns", adminCode)
	brewInstallCask("tor-browser", "Tor Browser")
	changeAppIcon("Tor Browser", "Tor Browser.icns", adminCode)
	brewInstallCask("spotify", "Spotify")
	changeAppIcon("Spotify", "Spotify.icns", adminCode)
	brewInstallCask("signal", "Signal")
	brewInstallCask("discord", "Discord")
	brewInstallCask("slack", "Slack")
	if runOpt == "5" || runOpt == "6" || runOpt == "7" {
		brewInstallCask("jetbrains-space", "JetBrains Space")
		changeAppIcon("JetBrains Space", "JetBrains Space.icns", adminCode)
	}

	if runOpt == "3" || runOpt == "6" || runOpt == "7" {
		brewInstallCask("dropbox", "Dropbox")
		brewInstallCask("dropbox-capture", "Dropbox Capture")
		brewInstallCask("sketch", "Sketch")
		brewInstallCask("zeplin", "Zeplin")
		brewInstallCask("blender", "Blender")
		changeAppIcon("Blender", "Blender.icns", adminCode)
		brewInstallCask("obs", "OBS")
		brewInstallCaskSudo("loopback", "Loopback", "/Applications/Loopback.app", adminCode)
		runApplication("Loopback")
	}

	if runOpt == "3" || runOpt == "4" || runOpt == "5" {
		brewInstallCaskSudo("blackhole-64ch", "BlackHole (64ch)", "/Library/Audio/Plug-Ins/HAL/BlackHoleXch.driver", adminCode)
	}

	if runOpt == "3" || runOpt == "4" {
		brewInstallCask("eclipse-ide", "Eclipse")
		changeAppIcon("Eclipse", "Eclipse.icns", adminCode)
		brewInstallCask("intellij-idea-ce", "IntelliJ IDEA CE")
		changeAppIcon("IntelliJ IDEA CE", "IntelliJ IDEA CE.icns", adminCode)
		brewInstallCask("android-studio", "Android Studio")
		changeAppIcon("Android Studio", "Android Studio.icns", adminCode)
		brewInstallCask("visual-studio-code", "Visual Studio Code")
		brewInstallCask("atom", "Atom")
		brewInstallCask("fork", "Fork")
		brewInstallCask("postman", "Postman")
		brewInstallCask("drawio", "draw.io")
		brewInstallCask("httpie", "HTTPie")
		installXAMPP(adminCode)
	} else if runOpt == "5" {
		brewInstallCask("iterm2", "iTerm")
		brewInstallCask("intellij-idea", "IntelliJ IDEA")
		changeAppIcon("IntelliJ IDEA", "IntelliJ IDEA.icns", adminCode)
		brewInstallCask("visual-studio-code", "Visual Studio Code")
		brewInstallCask("atom", "Atom")
		brewInstallCask("neovide", "Neovide")
		changeAppIcon("Neovide", "Neovide.icns", adminCode)
		brewInstallCask("github", "Github")
		brewInstallCask("fork", "Fork")
		brewInstallCask("docker", "Docker")
		brewInstallCask("tableplus", "TablePlus")
		brewInstallCask("postman", "Postman")
		brewInstallCask("httpie", "HTTPie")
		brewInstallCask("boop", "Boop")
		brewInstallCask("drawio", "draw.io")
		brewInstallCask("firefox-developer-edition", "Firefox Developer Edition")
		changeAppIcon("Firefox Developer Edition", "Firefox Developer Edition.icns", adminCode)
	} else if runOpt == "6" || runOpt == "7" {
		brewInstallCask("iterm2", "iTerm")
		brewInstallCask("intellij-idea", "IntelliJ IDEA")
		changeAppIcon("IntelliJ IDEA", "IntelliJ IDEA.icns", adminCode)
		brewInstallCask("visual-studio-code", "Visual Studio Code")
		brewInstallCask("atom", "Atom")
		brewInstallCask("neovide", "Neovide")
		changeAppIcon("Neovide", "Neovide.icns", adminCode)
		brewInstallCask("github", "Github")
		brewInstallCask("fork", "Fork")
		brewInstallCask("docker", "Docker")
		brewInstallCaskSudo("vmware-fusion", "VMware Fusion", "/Applications/VMware Fusion.app", adminCode)
		changeAppIcon("VMware Fusion", "VMware Fusion.icns", adminCode)
		brewInstallCask("tableplus", "TablePlus")
		brewInstallCask("proxyman", "Proxyman")
		brewInstallCask("postman", "Postman")
		brewInstallCask("paw", "Paw")
		brewInstallCask("httpie", "HTTPie")
		brewInstallCask("boop", "Boop")
		brewInstallCask("drawio", "draw.io")
		brewInstallCask("staruml", "StarUML")
		changeAppIcon("StarUML", "StarUML.icns", adminCode)
		brewInstallCask("vnc-viewer", "VNC Viewer")
		changeAppIcon("VNC Viewer", "VNC Viewer.icns", adminCode)
		brewInstallCask("forklift", "ForkLift")
	}

	shrcAppend := "# ANDROID STUDIO\n" +
		"export ANDROID_HOME=$HOME/Library/Android/sdk\n" +
		"export PATH=$PATH:$ANDROID_HOME/emulator\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools\n" +
		"export PATH=$PATH:$ANDROID_HOME/tools/bin\n" +
		"export PATH=$PATH:$ANDROID_HOME/platform-tools\n\n"
	appendContents(shrcPath, shrcAppend, 0644)

	if runOpt == "7" {
		brewInstallCaskSudo("codeql", "CodeQL", brewPrefix+"Caskroom/Codeql", adminCode)
		brewInstallCask("little-snitch", "Little Snitch")
		changeAppIcon("Little Snitch", "Little Snitch.icns", adminCode)
		brewInstallCask("micro-snitch", "Micro Snitch")
		changeAppIcon("Micro Snitch", "Micro Snitch.icns", adminCode)
		brewInstallCask("fsmonitor", "FSMonitor")
		changeAppIcon("FSMonitor", "FSMonitor.icns", adminCode)
		brewInstallCask("burp-suite-professional", "Burp Suite Professional")
		brewInstallCask("burp-suite", "Burp Suite Community Edition")
		brewInstallCask("owasp-zap", "OWASP ZAP")
		changeAppIcon("OWASP ZAP", "OWASP ZAP.icns", adminCode)
		brewInstallCaskSudo("wireshark", "Wireshark", "/Applications/Wireshark.app", adminCode)
		changeAppIcon("Wireshark", "Wireshark.icns", adminCode)
		brewInstallCaskSudo("zenmap", "Zenmap", "/Applications/Zenmap.app", adminCode)
		changeAppIcon("Zenmap", "Zenmap.icns", adminCode)
		installHopper(adminCode)
		brewInstallCask("cutter", "Cutter")
		// Install Ghidra // TODO: will add
		brewInstallCask("imazing", "iMazing")
		changeAppIcon("iMazing", "iMazing.icns", adminCode)
		brewInstallCask("apparency", "Apparency")
	}

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "install GUI applications!\n"
	macLdBar.Stop()
}

func macEnd() {
	macLdBar.Suffix = " Finishing... "
	macLdBar.Start()

	shrcAppend := "\n######## ADD CUSTOM VALUES UNDER HERE ########\n\n\n"
	appendContents(shrcPath, shrcAppend, 0644)

	brewUpgrade()
	brewCleanup()
	brewRemoveCache()

	macLdBar.FinalMSG = lstDot + clrGreen + "Succeed " + clrReset + "clean up homebrew's cache!\n"
	macLdBar.Stop()
}

func macMain(runOpt, runType, brewSts, adminCode string) {
	runEgMsg := lstDot + "Run " + clrPurple + runType + clrReset + " installation\n" + lstDot + brewSts + " homebrew with configure shell"
	if runOpt == "1" {
		fmt.Println(runEgMsg + ".")
	} else if runOpt == "2" {
		fmt.Println(runEgMsg + ", then install Dependencies, Languages and Terminal/CLI applications with set basic preferences.")
	} else if runOpt == "3" {
		fmt.Println(runEgMsg + ", then install Dependencies, Languages and Terminal/CLI/GUI applications with set basic preferences.")
	} else if runOpt == "4" {
		fmt.Println(runEgMsg + ", then install Dependencies, Languages and Terminal/CLI/GUI applications with set basic preferences.")
	} else if runOpt == "5" {
		fmt.Println(runEgMsg + ", then install Dependencies, Languages, Server, Database and Terminal/CLI/GUI applications with set basic preferences.")
	} else if runOpt == "6" {
		fmt.Println(runEgMsg + ", then install Dependencies, Languages, Server, Database, management DevTools and Terminal/CLI/GUI applications with set basic preferences.")
	} else if runOpt == "7" {
		fmt.Println(runEgMsg + ", then install Dependencies, Languages, Server, Database, management DevTools and Terminal/CLI/GUI applications with set basic preferences.")
	}

	alMsg := lstDot + "Use root permission to install "
	if runOpt == "1" || runOpt == "2" {
		if brewSts == "Install" {
			fmt.Println(alMsg + "homebrew")
		}
	} else {
		if brewSts == "Install" {
			alMsg = alMsg + "homebrew " + clrReset + "and " + clrPurple + "few applications" + clrReset + ": "
		} else if brewSts == "Update" {
			alMsg = alMsg + "few applications" + clrReset + ": "
		}
		if runOpt == "3" {
			fmt.Println(alMsg + "Loopback and BlackHole")
		} else if runOpt == "4" {
			fmt.Println(alMsg + "Java and BlackHole")
		} else if runOpt == "5" {
			fmt.Println(alMsg + "Java, BlackHole and VMware Fusion")
		} else if runOpt == "6" {
			fmt.Println(alMsg + "Java, Loopback and VMware Fusion")
		} else if runOpt == "7" {
			fmt.Println(alMsg + "Java, Loopback, VMware Fusion, Wireshark and Zenmap")
		}
	}

	if runOpt == "1" {
		macBegin(adminCode)
		macEnv()
	} else if runOpt == "2" {
		macBegin(adminCode)
		macEnv()
		macDependency(runOpt)
		macTerminal(runOpt)
		macLanguage(runOpt, adminCode)
		macCLIApp(runOpt)
	} else if runOpt == "3" {
		macBegin(adminCode)
		macEnv()
		macDependency(runOpt)
		macTerminal(runOpt)
		macLanguage(runOpt, adminCode)
		macCLIApp(runOpt)
		macGUIApp(runOpt, adminCode)
	} else if runOpt == "4" {
		macBegin(adminCode)
		macEnv()
		macDependency(runOpt)
		macTerminal(runOpt)
		macLanguage(runOpt, adminCode)
		macServer()
		macDatabase()
		macCLIApp(runOpt)
		macGUIApp(runOpt, adminCode)
	} else if runOpt == "5" {
		macBegin(adminCode)
		macEnv()
		macDependency(runOpt)
		macTerminal(runOpt)
		macLanguage(runOpt, adminCode)
		macServer()
		macDatabase()
		macCLIApp(runOpt)
		macGUIApp(runOpt, adminCode)
	} else if runOpt == "6" {
		macBegin(adminCode)
		macEnv()
		macDependency(runOpt)
		macTerminal(runOpt)
		macLanguage(runOpt, adminCode)
		macServer()
		macDatabase()
		macDevVM()
		macCLIApp(runOpt)
		macGUIApp(runOpt, adminCode)
	} else if runOpt == "7" {
		macBegin(adminCode)
		macEnv()
		macDependency(runOpt)
		macTerminal(runOpt)
		macLanguage(runOpt, adminCode)
		macServer()
		macDatabase()
		macDevVM()
		macCLIApp(runOpt)
		macGUIApp(runOpt, adminCode)
	}
	macEnd()
}

func macExtend(runOpt, adminCode string) {
	if runOpt != "1" {
		var (
			g4sOpt      string
			osUpdateOpt string
			osRebootOpt string
		)
		askOpt := "If you wish to continue type (Y) then press return: "

		fmt.Print(clrCyan + "\nConfigure git global easily\n" + clrReset + "To continue we setup git global configuration.\n" + askOpt)
		_, errG4sOpt := fmt.Scanln(&g4sOpt)
		if errG4sOpt != nil {
			g4sOpt = "Enter"
		}
		if g4sOpt == "y" || g4sOpt == "Y" || g4sOpt == "yes" || g4sOpt == "Yes" || g4sOpt == "YES" {
			clearLine(3)
			confG4s()
		} else {
			clearLine(4)
		}

		fmt.Println("\nFinished all things!\n") // Finish messages for update or restart OS

		fmt.Print(clrCyan + "macOS software update\n" + clrReset + "To continue we update macOS software update.\n" + askOpt)
		_, errUpdateOpt := fmt.Scanln(&osUpdateOpt)
		if errUpdateOpt != nil {
			osUpdateOpt = "Enter"
		}
		if osUpdateOpt == "y" || osUpdateOpt == "Y" || osUpdateOpt == "yes" || osUpdateOpt == "Yes" || osUpdateOpt == "YES" {
			clearLine(4)
			systemUpdate()
			systemReboot(adminCode)
		} else {
			clearLine(3)
			if runOpt == "3" || runOpt == "6" || runOpt == "7" {
				fmt.Print(clrCyan + "Restart macOS to apply the changes\n" + clrReset + "To continue we restart macOS.\n" + askOpt)
				_, errRebootOpt := fmt.Scanln(&osRebootOpt)
				if errRebootOpt != nil {
					osRebootOpt = "Enter"
				}
				if osRebootOpt == "y" || osRebootOpt == "Y" || osRebootOpt == "yes" || osRebootOpt == "Yes" || osRebootOpt == "YES" {
					clearLine(3)
					systemReboot(adminCode)
				} else {
					clearLine(6)
				}
			} else {
				clearLine(4)
			}
		}
	}
}

func main() {
	fmt.Println(clrBlue + "\nDev4mac\n" + clrGrey + "Dev4os version " + appVer + clrReset + "\n")

	runLdBar.Suffix = " Checking network status... "
	runLdBar.Start()

	var (
		brewSts string
		runOpt  string
		runType string
		endMsg  string
	)

	if checkExists(cmdPMS) == true {
		brewSts = "Update"
	} else {
		brewSts = "Install"
	}

	if checkNetStatus() != true {
		runLdBar.FinalMSG = clrRed + "Network connect failed" + clrReset + "\n"
		runLdBar.Stop()
		fmt.Println(errors.New(lstDot + "Please check your internet connection.\n"))
		goto exitPoint
	}

	runLdBar.Stop()

	fmt.Println(clrCyan + "The Development tools of Essential and Various for macOS\n" + clrReset +
		lstDot + "Choose an installation option.\n" + lstDot + "If you need help, visit https://github.com/leelsey/Dev4os.\n" +
		"\t1. Minimal\n\t2. Basic\n\t3. Creator\n\t4. Beginner\n\t5. Developer\n\t6. Professional\n\t7. Specialist\n\t0. Exit\n")

runOptPoint:
	for {
		fmt.Print("Select command: ")
		_, err := fmt.Scanln(&runOpt)
		if err != nil {
			runOpt = "Null"
		}
		if runOpt == "1" {
			runType = "Minimal"
		} else if runOpt == "2" {
			runType = "Basic"
		} else if runOpt == "3" {
			runType = "Creator"
		} else if runOpt == "4" {
			runType = "Beginner"
		} else if runOpt == "5" {
			runType = "Developer"
		} else if runOpt == "6" {
			runType = "Developer"
		} else if runOpt == "7" {
			runType = "Specialist"
		} else if runOpt == "0" || runOpt == "q" || runOpt == "e" || runOpt == "quit" || runOpt == "exit" {
			fmt.Println(lstDot + "Exited Dev4mac.")
			goto exitPoint
		} else {
			fmt.Println(fmt.Errorf(lstDot + clrYellow + runOpt + clrReset +
				" is invalid option. Please choose number " + clrRed + "0-7" + clrReset + "."))
			tryLoop++
			goto runOptPoint
		}
		break
	}
	clearLine(12 + tryLoop*2)

	if checkPermission(runOpt, brewSts) == true {
		if adminCode, adminStatus := checkPassword(); adminStatus == true {
			needPermission(adminCode)
			macMain(runOpt, runType, brewSts, adminCode)
			macExtend(runOpt, adminCode)
		} else {
			goto exitPoint
		}
	} else {
		macMain(runOpt, runType, brewSts, "")
		macExtend(runOpt, "")
	}

	endMsg = "\n----------Finished!----------\nPlease" + clrRed + " RESTART " + clrReset + "your terminal!\n" +
		lstDot + "Enter this on terminal: source ~/.zprofile && source ~/.zshrc\n" + lstDot + "Or restart the Terminal.app by yourself.\n"
	if runOpt == "3" || runOpt == "6" || runOpt == "7" {
		fmt.Println(endMsg + lstDot + "Also you need " + clrRed + "RESTART macOS " + clrReset + " to apply " + "the changes.\n")
	} else {
		fmt.Println(endMsg)
	}

exitPoint:
	return
}
