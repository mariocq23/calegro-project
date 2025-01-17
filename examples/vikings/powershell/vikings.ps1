param([string]$gamePath)

function Invoke-Game([string]$emulatorPath, [string]$gamePath){
    try {

        Write-Host "Emulator path: $emulatorPath";

        Write-Host "Game path: $gamePath";

        $gameDirectory = Split-Path -Path $gamePath -Parent -Resolve;

        Write-Host "Game directory: $gameDirectory";

        $confFile = New-ConfigurationFile $gamePath $gameDirectory;

        Write-Host "Conf file output: $confFile";

        Save-ConfigurationFile($gamePath, $confFile);

        Start-Process -FilePath $emulatorPath -ArgumentList "-conf \"{$confFile}\" -noconsole";
    }
    catch {
        Write-Host "Something went wrong!"
    }
}

$dosBoxPath = "C:\Program Files (x86)\DOSBox-0.74-3\DOSBox.exe";

Invoke-Game $dosBoxPath $gamePath;

function New-ConfigurationFile([string]$gamePath, [string]$gameDirectory){
    $drive = "C";

    Write-Host "Drive: $drive";

    $executableName = Split-Path -Path $gamePath -Leaf -Resolve;
    
    Write-Host "Executable: $executableName";


    return "
[autoexec]
mount  ${drive} ${gameDirectory}
${drive}:
${executableName}";
}

function Save-ConfigurationFile([string]$gameDirectory, [string]$confFile){
    Set-Content (Join-Path -Path ${gameDirectory} -ChildPath "dosbox.conf") $confFile;
}