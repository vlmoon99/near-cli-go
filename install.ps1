$Repo = "vlmoon99/near-cli-go"
$LatestUrl = "https://api.github.com/repos/$Repo/releases/latest"
$InstallDir = "$env:ProgramFiles\NearCLI"
$BinaryName = "near-go.exe"

Write-Host "🔍 Detecting system architecture..."
$Arch = (Get-WmiObject Win32_Processor).Architecture

if ($Arch -eq 9) {  # x86_64 (AMD64)
    $Filename = "near-cli-windows-amd64.exe"
} else {
    Write-Host "❌ Unsupported architecture: $Arch"
    exit 1
}

Write-Host "⬇️ Downloading $Filename..."
$Response = Invoke-RestMethod -Uri $LatestUrl
$Url = ($Response.assets | Where-Object { $_.name -eq $Filename }).browser_download_url

if (-not $Url) {
    Write-Host "❌ Failed to find the latest release for $Filename"
    exit 1
}

$DownloadPath = "$env:TEMP\$Filename"
Invoke-WebRequest -Uri $Url -OutFile $DownloadPath

Write-Host "🔧 Installing Near CLI..."
if (!(Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Move-Item -Force -Path $DownloadPath -Destination "$InstallDir\$BinaryName"

# Add to system PATH
$CurrentPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
if ($CurrentPath -notlike "*$InstallDir*") {
    [System.Environment]::SetEnvironmentVariable("Path", "$CurrentPath;$InstallDir", [System.EnvironmentVariableTarget]::Machine)
    Write-Host "🔄 System PATH updated. Please restart your terminal."
}

Write-Host "✅ Installation complete! Run '$BinaryName' to start using it."
