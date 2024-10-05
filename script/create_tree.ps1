param (
    [string]$path = "."
)

function Get-Tree {
    param (
        [string]$dir,
        [string]$indent = ""
    )

    $items = Get-ChildItem -Path $dir

    foreach ($item in $items) {
        if ($item.PSIsContainer) {
            Write-Output "$indent|--- $($item.Name)"
            Get-Tree -dir $item.FullName -indent "$indent|   "
        } else {
            Write-Output "$indent|--- $($item.Name)"
        }
    }
}

# Start with the base directory
Write-Output $path
Get-Tree -dir $path
