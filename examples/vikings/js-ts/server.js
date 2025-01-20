const express = require('express');
const { spawn } = require('child_process');
const path = require('path');

const app = express();
const port = 3000;

app.get('/run-dosbox', (req, res) => {
    const dosboxPath = 'C:\\Program Files (x86)\\DOSBox-0.74-3\\DOSBox.exe'; // Replace with your DOSBox path
    const gamePath = 'C:\\Games\\DosGames\\Lost Vikings'; // Replace with your game directory
    const gameExe = 'VIKINGS.EXE'; // Replace with the game's executable

    if (!dosboxPath || !gamePath || !gameExe) {
      return res.status(500).send('DOSBox or game paths not configured.');
    }

    const tempConfPath = path.join(__dirname, 'temp_dosbox.conf'); // Create a temporary config file
    const fs = require('fs');
    fs.writeFileSync(tempConfPath, `
[autoexec]
mount c "${gamePath}"
c:
${gameExe}
exit
`);

    const dosboxProcess = spawn(dosboxPath, ['-conf', tempConfPath]);

    dosboxProcess.on('error', (err) => {
        console.error('Failed to start DOSBox:', err);
        res.status(500).send('Failed to start DOSBox.');
        fs.unlinkSync(tempConfPath); //Clean up config file
    });

    dosboxProcess.on('close', (code) => {
        console.log(`DOSBox exited with code ${code}`);
        res.send('DOSBox started (in the background). Check your desktop.');
        fs.unlinkSync(tempConfPath); //Clean up config file
    });

    dosboxProcess.stdout.on('data', (data) => {
        console.log(`DOSBox stdout: ${data}`); // Log DOSBox output (optional)
    });

    dosboxProcess.stderr.on('data', (data) => {
        console.error(`DOSBox stderr: ${data}`); // Log DOSBox errors (optional)
    });
});

app.listen(port, () => {
    console.log(`Server listening at http://localhost:${port}`);
});