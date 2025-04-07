use std::env;
use std::fs;
use std::path::PathBuf;
use serde::Deserialize;
use serde_yaml;
use std::process::Command;

#[derive(Debug, Deserialize)]
struct Config {
    dosbox_executable: String,
    game_executable: String,
    // actions: Option<Vec<Action>>,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Get command-line arguments
    let args: Vec<String> = env::args().collect();

    println!("Command-line arguments: {:?}", args);

    if args.len() != 2 {
        eprintln!("Usage: {} <path_to_config.yaml>", args[0]);
        return Ok(()); // Or Err("Missing or incorrect number of arguments".into());
    }

    let config_file_path = &args[1];

    println!("Using configuration file: {}", config_file_path);

    let config_file = fs::read_to_string(config_file_path)?;

    let config: Config = serde_yaml::from_str(&config_file)?;

    println!("Configuration loaded: {:?}", config);

    // Launch DOSBox with Lost Vikings
    let mut command = Command::new(&config.dosbox_executable);
    command.arg(&config.game_executable);

    println!("Launching DOSBox with command: {:?}", command);

    let status = command.spawn()?;

    println!("DOSBox process started with PID: {:?}", status.id());

    Ok(())
}