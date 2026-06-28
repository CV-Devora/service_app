param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("up", "down")]
    [string]$Action,

    [Parameter(Mandatory = $true)]
    [string]$DatabaseUrl
)

$ErrorActionPreference = "Stop"

function Invoke-Psql {
    param(
        [string]$SqlFile
    )

    Write-Host "Running $SqlFile"
    & psql $DatabaseUrl -f $SqlFile
    if ($LASTEXITCODE -ne 0) {
        throw "psql failed for $SqlFile"
    }
}

$migrationDir = Join-Path $PSScriptRoot "..\migrations"
$files = Get-ChildItem -Path $migrationDir -Filter *.sql | Sort-Object Name

if ($Action -eq "up") {
    foreach ($file in $files) {
        Invoke-Psql -SqlFile $file.FullName
    }
    Write-Host "Migrations up completed."
    exit 0
}

if ($Action -eq "down") {
    foreach ($file in ($files | Sort-Object Name -Descending)) {
        $content = Get-Content $file.FullName -Raw
        if ($content -match 'CREATE TABLE IF NOT EXISTS\s+([a-zA-Z0-9_]+)') {
            $table = $Matches[1]
            $dropFile = Join-Path $env:TEMP "$table.drop.sql"
            Set-Content -Path $dropFile -Value "DROP TABLE IF EXISTS $table CASCADE;"
            Invoke-Psql -SqlFile $dropFile
            Remove-Item $dropFile -Force
        }
    }
    Write-Host "Migrations down completed."
    exit 0
}
