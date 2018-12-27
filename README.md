# uprun - get your tasks up and running
This little helper manages the lifecycle of your tasks.  It is intended to be used inside docker containers.

## Features
* start multiple services
* redirect stdout/stderr
* shutdown everything if one task fails
* export docker secrets to environment variables
* capture SIGINT / SIGTERM to properly shutdown all services

## Configuration
    secret_dir = "/run/secrets"
    
    service "super-awesome-node" {
      command = "node test2.js"
      capture_stdout = true
      capture_stderr = true
    }
    
    service "boring-node" {
      command = "node test.js"
      secret_prefix = "test_"
      capture_stdout = true
      capture_stderr = true
      ignore_failure = true
    }
    
## Commandline Options
    $: ./uprun --help
    Usage of ./uprun:
      -colors
            force logging with colors
      -conf string
            path to config file (default "uprun.hcl")
