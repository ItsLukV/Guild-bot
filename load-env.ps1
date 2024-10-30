# load-env.ps1
Get-Content .env | ForEach-Object {
    if ($_ -match "^(.*)=(.*)$") {
        $name = $matches[1]
        $value = $matches[2] -replace '^"|"$', ''  # Removes surrounding double quotes
        [System.Environment]::SetEnvironmentVariable($name, $value, "Process")
    }
}
