param(
    [string]$DatabaseUrl = "postgres://postgres:Tambunan140705@localhost:5432/jason_jewelry?sslmode=disable"
)

$ErrorActionPreference = "Stop"

& pwsh -NoProfile -ExecutionPolicy Bypass -File (Join-Path $PSScriptRoot "migrate.ps1") -Action up -DatabaseUrl $DatabaseUrl
