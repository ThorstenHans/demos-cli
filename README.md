# Demos Over SSH

> TL;DR: Don't let bad conference network nuke your demos. Instead offload network intense demos to a dedicated machine and execute them over SSH. Resulting in less data being transferred over the conference network.

## Installation

Download the `demos` binary from the recent [release page](https://github.com/ThorstenHans/demos-over-ssh/releases) and ensure it is executable (`chmod +x demos`). You can either execute the binary from the folder you've downloaded it to, or move it into your `PATH`.

## Configuration

The `demos` application requires two different sets of configuration data: 

- Configuration for establishing an SSH connection to the "jump box"
- Your actual demos

### SSH Configuration

After you've downloaded the `demos` executable run the `demos configure` command. It'will prompt for your SSH configuration data. You must provide:

- IPv4 address of the jump box
- Desired SSH Port (default `22`)
- SSH user name
- Password for the SSH user

Configuration data is encrypted at REST and stored in your user profile (`$HOME/.demo/demo.config`).


### Configuring Demos

The `demos` app comes with a sample demo backed in... However, you want to provide your own demos (obviously). To give you a head start, you can use the `demos eject` command, which will create the `$HOME/.demo/demos.json` file for you. You can add as many demos to the JSON array as you want. Here some additional information for specifying your demos:

- A demo can consist of an unlimited number of steps. 
- A demo step can either be of `kind` code (`1`) or text (`0`)
  - Text steps are printed directly to `stdout`
  - Every code step is executed over SSH in a dedicated session and it's output is forwarded to the local `stdout`
- A demo must have a `name`, a `cliCommand`, an `alias` and a `description`
  - `cliCommand` and `alias` have unique constraints which are enforced when loading the demos into the binary at runtime

```json
[
  {
    "name": "Sample Load Test",
    "cliCommand": "load-test",
    "alias": "lt",
    "description": "Run a Sample Load Test",
    "steps": [
      { "command": "We'll now sent 100 requests Google", "kind": 0 },
      { "command": "hey -c 10 -n 100 https://www.google.com", "kind": 1 },
      { "command": "100 requests sent!", "kind": 0 }
    ]
  }
]
``` 

## Dynamic CLI Commands

For each demo provided, a new sub-command is added under `demos run` using the specified `cliCommand` as command name and setting the provided `alias` as command alias.

Taking the previously shown `demos.json` into context, you will end up with the following commands available as part of the `demos` CLI:

```bash
# default command name
demos run load-test

# command alias
demos run lt
```

## Printing your Demos

As demos are automatically hooked up into the CLI, you can simply execute `demos run` to get a list of all commands (demos) that could be executed. 