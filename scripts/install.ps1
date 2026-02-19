Param(
  [string]$Repo = "benji/cogito",
  [string]$Version = "latest",
  [string]$InstallDir = "$env:LOCALAPPDATA\\Programs\\cogito"
)

$ErrorActionPreference = "Stop"
$BinName = "cogito"

$arch = $env:PROCESSOR_ARCHITECTURE
switch ($arch) {
  "AMD64" { $Arch = "amd64" }
  "ARM64" { $Arch = "arm64" }
  default { throw "Unsupported architecture: $arch" }
}

if ($Version -eq "latest") {
  $release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
  $Version = $release.tag_name
}

if (-not $Version) { throw "Could not determine latest version." }

$zip = "${BinName}_${Version.TrimStart('v')}_windows_${Arch}.zip"
$url = "https://github.com/$Repo/releases/download/$Version/$zip"

$TempDir = New-Item -ItemType Directory -Path ([System.IO.Path]::GetTempPath()) -Name ([System.Guid]::NewGuid().ToString())
try {
  $zipPath = Join-Path $TempDir $zip
  Invoke-WebRequest -Uri $url -OutFile $zipPath
  Expand-Archive -Path $zipPath -DestinationPath $TempDir

  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
  Copy-Item -Path (Join-Path $TempDir "$BinName.exe") -Destination (Join-Path $InstallDir "$BinName.exe") -Force

  Write-Host "Installed $BinName to $InstallDir\\$BinName.exe"
  Write-Host "Add $InstallDir to your PATH to run it from anywhere."
} finally {
  Remove-Item -Recurse -Force $TempDir
}
