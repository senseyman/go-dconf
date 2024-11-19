# Dynamic Configuration Manager (`go-dconf`)

`go-dconf` is a Go library designed for managing dynamic configurations in applications. It allows you to update and retrieve configuration values without restarting your app. The library supports periodic polling and real-time updates using Redis (as an example database) and provides a flexible, generic, and thread-safe configuration management mechanism.
You can use your own storage if you just implement the basic `Repository` interface.

---

## Features

- **Dynamic Configuration Updates**:
    - Fetch and apply updated configuration values at runtime.
- **Periodic Polling**:
    - Automatically refresh configurations at regular intervals.
- **Thread-Safe Access**:
    - Use `sync.Mutex` to ensure safe concurrent access to configurations.
- **Generic Support**:
    - Use the library with any custom configuration structure.
- **Extensibility**:
    - Abstract repository layer allows support for other databases (e.g., PostgreSQL, MongoDB).

---

## Installation

```bash
go get github.com/senseyman/go-dconf
```

# Getting Started
## Define Your Configuration
Create a custom struct for your application's configuration:
```go
type AppConfig struct {
    Name string
    Port int
}
```
## Initialize the Manager
Set up the configuration manager with an initial configuration:
```go
import (
    "context"
    "log"
    "time"

    "github.com/senseyman/go-dconf/manager"
    "github.com/senseyman/go-dconf/repository/redis"
)

func main() {
    ctx := context.Background()

    // Define your initial configuration
    initConfig := AppConfig{
        Name: "MyApp",
        Port: 8080,
    }

    // Set up the Redis client
    repo, err := redis.New(ctx, redis.Config{
            Address:  "localhost:6379",
            Password: "",
            DB:       0,
            }, "myApp")
    if err != nil {
        log.Fatal(err)
    }
	
    // Initialize the ConfigManager
    configManager := manager.New(repo, initConfig, 10*time.Second)

    // Start the ConfigManager
    go func() {
        if err := configManager.Run(ctx, nil); err != nil {
            log.Fatalf("Failed to start config manager: %v", err)
        }
    }()

    // Access the current configuration
    currentConfig := configManager.GetConfig()
    log.Printf("Current Configuration: %+v", currentConfig)
}

```
# Key Components
## ConfigManager
Manages dynamic configurations. Key methods include:

GetConfig: Safely retrieve the current configuration.
LoadConfig: Manually reload the configuration.
Run: Periodically refresh the configuration or subscribe to real-time updates.
## Repository
An abstract layer for data storage. A Redis implementation is provided as an example.

### Example: Redis Repository
```go
repo := redis.New(ctx, redis.Config{
                Address:  "localhost:6379",
                Password: "",
                DB:       0,
            }, "myApp")
```

# Contributing
Contributions are welcome! Please submit pull requests or report issues on the GitHub repository.