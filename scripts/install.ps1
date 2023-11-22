# Determine OS and Architecture
$os = "windows" # Assuming Windows, as this is a PowerShell script
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Define download URL
$tar = "celify-$os-$arch.zip"
$url = "https://github.com/rdalbuquerque/celify/releases/latest/download"

# Get the latest version
$version = Invoke-RestMethod -Uri "https://api.github.com/repos/rdalbuquerque/celify/releases/latest" | Select-Object -ExpandProperty tag_name
$filename = "celify_${version}_$os-$arch.exe"

# Download and install
Write-Host "Downloading version $version of celify-$os-$arch..."
Invoke-WebRequest -Uri "$url/$tar" -OutFile "$env:TEMP\$tar"
if (!(Test-Path "$env:TEMP\$tar")) {
    Write-Host "Failed to download celify-$os-$arch"
    throw "Failed to download celify-$os-$arch"
}

# Extracting the ZIP file
Expand-Archive -Path "$env:TEMP\$tar" -DestinationPath $env:TEMP -Force

# Move the file to a specific location
$destination = "$env:LOCALAPPDATA\celify\celify.exe"
Write-Host "Moving $env:TEMP\$filename to $destination"
$null = New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\celify"
if (Test-Path $destination) {
    $answer = Read-Host "celify already installed, would you like to overwrite it? (y/n)"
    if ($answer -eq "y") {
        Move-Item -Path "$env:TEMP\$filename" -Destination $destination -Force
    } else {
        Write-Host "Aborting..."
        Remove-Item "$env:TEMP\$tar"
        return
    }
}

# Add the celify.exe to the PATH if it's not
$paths = $env:Path -split ";"
if ($paths -notcontains "$env:LOCALAPPDATA\celify") {
    Write-Host "Adding $env:LOCALAPPDATA\celify to the PATH"
    $env:Path += ";$env:LOCALAPPDATA\celify"
    Write-Host "Please alter your PATH variable to include $env:LOCALAPPDATA\celify permanently"
}
celify --version
if ($LASTEXITCODE -ne 0) {
    throw "Failed to install celify"
} else {
    Write-Host "celify installed successfully"
}

# Clean up
Remove-Item "$env:TEMP\$tar"
