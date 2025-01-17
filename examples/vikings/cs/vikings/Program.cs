// See https://aka.ms/new-console-template for more information
using System.Diagnostics;

var dosboxPath = @"C:\Program Files (x86)\DOSBox-0.74-3\DOSBox.exe"; // Replace with your DOSBox path
var gamePath = args[0]; // Replace with the path to vikings.exe
RunGame(dosboxPath, gamePath);

Console.ReadKey();


static void RunGame(string dosboxPath, string gamePath)
{
    try
    {
        // 1. Create a temporary configuration file (.conf)
        var tempConfFile = Path.GetTempFileName() + ".conf";
        File.WriteAllText(tempConfFile, GenerateGameConfiguration(gamePath));

        // 2. Start DOSBox process
        var psi = new ProcessStartInfo
        {
            FileName = dosboxPath,
            Arguments = $"-conf \"{tempConfFile}\" -noconsole",
            UseShellExecute = false,
            CreateNoWindow = true,
            RedirectStandardOutput = true, // Capture output if needed (for debugging)
            RedirectStandardError = true // Capture errors if needed (for debugging)
        };

        using (var process = new Process { StartInfo = psi })
        {
            process.Start();

            var output = process.StandardOutput.ReadToEnd();
            var error = process.StandardError.ReadToEnd();

            process.WaitForExit();

            if (process.ExitCode != 0)
            {
                Console.WriteLine($"DOSBox exited with code {process.ExitCode}");
                if (!string.IsNullOrEmpty(error))
                {
                    Console.WriteLine($"DOSBox Error: {error}");
                }
            }
            else if (!string.IsNullOrEmpty(output))
            {
                Console.WriteLine($"DOSBox Output: {output}");
            }

            File.Delete(tempConfFile);
        }
    }
    catch (Exception ex)
    {
        Console.WriteLine($"Error running Lost Vikings: {ex.Message}");
    }
}

static string GenerateGameConfiguration(string gamePath)
{
    // Extract the drive and path for mounting
    var drive = Path.GetPathRoot(gamePath).Substring(0, 1).ToLower(); // e.g., "c"
    var gameDirectory = Path.GetDirectoryName(gamePath);
    var executableName = Path.GetFileName(gamePath);


    return $@"
[autoexec]
mount  {drive} {Path.Combine(drive, gameDirectory)}
{drive}:
START.BAT
";
}

