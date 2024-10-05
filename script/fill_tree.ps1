# Define the directory to search (replace with your directory path)
$rootDirectory = "../."

# Function to check if a directory is empty
function IsDirectoryEmpty($directory) {
    $childItems = Get-ChildItem $directory
    return ($childItems.Count -eq 0)
}

# Function to create README.md file in empty directories
function CreateReadmeInEmptyDirectories($rootDir) {
    Get-ChildItem -Directory -Recurse $rootDir | ForEach-Object {
        if (IsDirectoryEmpty($_.FullName)) {
            $readmePath = Join-Path $_.FullName "README.md"
            if (-not (Test-Path $readmePath)) {
                New-Item -ItemType File -Path $readmePath -Force
                Write-Output "Created README.md in $($_.FullName)"
            }
        }
    }
}

# Call the function to create README.md files in empty directories
CreateReadmeInEmptyDirectories $rootDirectory
