
if (!$args[0]) {
    Write-Output "No input file given"
    exit
}
$fp=(Get-Item -Path ".\" -Verbose).FullName
$infile = $($args[0])
$tempfile = [io.path]::GetFileNameWithoutExtension($infile) + ".temp"
$outfile = $fp + "\" + [io.path]::GetFileNameWithoutExtension($infile) + ".exe"
dasm $infile -f2 -o"$tempfile"
dasmf $tempfile $outfile
