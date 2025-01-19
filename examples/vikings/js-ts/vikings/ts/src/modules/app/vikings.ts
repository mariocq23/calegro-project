// hello.ts

import * as path from 'path';
import * as fs from 'fs';
import * as os from 'os';
import * as child_process from 'child_process';

interface ProcessOptions extends child_process.SpawnOptions{
    command: string;
    args?: string[];
}

function runGame(dosboxPath: string, gamePath: string): void {

    try{
        const confFile = path.join(fs.mkdtempSync(path.join(os.tmpdir(), '')), '.conf');
        console.log(confFile);
        fs.writeFileSync(confFile, generateGameConfiguration(gamePath));

        console.log('Entering child_process');

        const options : ProcessOptions = {
            command: dosboxPath,
            args: new List<> y,
        }
        
        

        const dosBoxProcess = child_process.spawn(dosboxPath, ["-conf \"{confFile}\" -noconsole"]);

        
        dosBoxProcess.stdout.on('data', (data) => {
            console.log('stdout: ${data}');
        });

        dosBoxProcess.stderr.on('data', (data) => {
            console.error('stderr: $(data)');
        });

        dosBoxProcess.on('close', (code) => {
            console.log('child process exited with code ${code}');
        });

        fs.unlinkSync(confFile);

    }
    catch(error){
        console.log('error: ${error}');
    }
}

function generateGameConfiguration(gamePath: string): string {
    // Extract the drive and path for mounting
    const drive = "c";
    const gameDirectory = path.dirname(gamePath);

    return `
[autoexec]
mount ${drive} ${path.join(drive, gameDirectory)}
${drive}:
VIKINGS.EXE
`;
}

function main(args: string[]) {

    console.log("Testing");
    console.log(args[0])
    debugger;
    var dosboxPath: string = "C:/Program Files (x86)/DOSBox-0.74-3/DOSBox.exe";
    
    runGame(dosboxPath, args[0]);
}
  
  // Get arguments starting from the third element (index 2), as the first two are node and the script path.

  main(process.argv.slice(2));