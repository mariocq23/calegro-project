import subprocess
import os
import shutil

def run_lost_vikings_dosbox(dosbox_path, game_path, config_path=None):

    if not os.path.exists(dosbox_path):
        raise FileNotFoundError(f"DOSBox executable not found at: {dosbox_path}")
    if not os.path.exists(game_path):
        raise FileNotFoundError(f"Game directory not found at: {game_path}")
    
    # Create a temporary configuration file if none is provided
    if config_path is None:
        temp_config = os.path.join(os.getcwd(), "temp_dosbox.conf")
        with open(temp_config, "w") as f:
            f.write("[autoexec]\n")
            f.write(f"mount c \"{game_path}\"\n") # Mount the game directory as C:
            f.write("c:\n")
            #Find the executable. There can be more than one.
            possible_exes = [f for f in os.listdir(game_path) if f.lower().endswith(".exe")]
            if not possible_exes:
                raise FileNotFoundError("No executable found in game directory")
            f.write(possible_exes[0] + "\n") # Run the game's .EXE file
            f.write("exit\n") # Close DOSBox after the game exits (optional)
        config_path = temp_config
    elif not os.path.exists(config_path):
        raise FileNotFoundError(f"Config file not found at: {config_path}")


    try:
        subprocess.run([dosbox_path, "-conf", config_path], check=True)
    except subprocess.CalledProcessError as e:
        print(f"DOSBox exited with error: {e}")
    finally:
        if config_path == temp_config:
            os.remove(temp_config) #Clean up the temporary config file.

# Example usage: Replace with your actual paths
dosbox_exe = "C:\\Program Files (x86)\\DOSBox-0.74-3\\DOSBox.exe" # Or your DOSBox path
lost_vikings_dir = "C:\\Games\\DosGames\\lost-vikings" # Or your Lost Vikings game directory

try:
    run_lost_vikings_dosbox(dosbox_exe, lost_vikings_dir)
    print("Lost Vikings started successfully (hopefully!).")
except FileNotFoundError as e:
    print(f"Error: {e}")
except Exception as e:
    print(f"An unexpected error occurred: {e}")