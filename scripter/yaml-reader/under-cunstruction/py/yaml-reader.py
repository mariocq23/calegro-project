from typing import List, Optional, Dict, Any
import yaml
from dataclasses import dataclass, field

@dataclass
class Header:
    headerImport: List[str]
    id: str
    name: str
    isChildNode: bool
    inherits: str
    implements: List[str]

@dataclass
class Security:
    publicPassword: str
    privatePasswordLocation: str
    certificateLocation: int
    templateOrSource: str

@dataclass
class Configuration:
    user: str
    agent: str
    executionMode: int
    bypassSecurity: bool
    security:Security
    location:str
    contextName:str
    encoding:str

@dataclass
class Action:
    host: str
    port: int
    workers: int

@dataclass
class Contexts:
    host: str
    port: int
    workers: int

@dataclass
class Steps:
    host: str
    port: int
    workers: int

@dataclass
class YamlFile:
    header: Header
    configuration: Configuration
    action: Action
    contexts: Contexts
    steps: Steps

def load_config(filepath: str) -> YamlFile:
    """Loads and deserializes a YAML configuration file into a typed dataclass."""
    try:
        with open(filepath, 'r') as f:
            data = yaml.safe_load(f)
            if not isinstance(data, dict):
                raise ValueError("YAML data must be a dictionary.")

            # Deserialize the nested objects carefully
            header_data = data.get('header')
            if not isinstance(header_data, dict):
                raise ValueError("Header must be a dictionary.")
            header_yaml = Header(**header_data)

            configuration_data = data.get('configuration')
            if not isinstance(configuration_data, dict):
                raise ValueError("Configuration must be a dictionary.")
            configuration_yaml = Configuration(**configuration_data)

            action_data = data.get('action')
            if not isinstance(action_data, dict):
                raise ValueError("Action must be a dictionary.")
            action_yaml = Action(**action_data)

            contexts_data = data.get('contexts')
            if not isinstance(contexts_data, list):
                raise ValueError("Contexts must be a list.")
            contexts_yaml = [Contexts(**context_data) for context_data in contexts_data]

            steps_data = data.get('steps')
            if not isinstance(steps_data, list):
                raise ValueError("Steps must be a list.")
            steps_yaml = [Contexts(**step_data) for step_data in steps_data]

            app_config = YamlFile(
                header=header_yaml,
                configuration=configuration_yaml,
                action=action_yaml,
                contexts=contexts_yaml,
                steps=steps_yaml,
                debug=data.get('debug', False) # Handle optional field with get() and default
            )
            return app_config

    except FileNotFoundError:
        raise FileNotFoundError(f"Configuration file not found: {filepath}")
    except yaml.YAMLError as e:
        raise ValueError(f"Error parsing YAML: {e}")
    except (TypeError, KeyError, ValueError) as e: # Catch deserialization errors
        raise ValueError(f"Invalid configuration format: {e}")


# Example usage:
if __name__ == "__main__":
    try:
        config = load_config("C:\\git\\calegro-project\\examples\\vikings\\vikings.yaml")  # Replace with your YAML file
        print(config)
        print(config.database.host) # Access nested attributes
        for server in config.servers:
            print(f"Server: {server.host}:{server.port}")
    except (FileNotFoundError, ValueError) as e:
        print(f"Error: {e}")