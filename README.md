# mcserverdl

A command-line tool and Go library to download and install various Minecraft server types with ease.

`mcserverdl` simplifies the process of setting up a Minecraft server by handling the download and installation steps for you. It can perform a direct download, or automatically patch a vanilla server, depending on the server type and version.

## Features

- **Multiple Server Types**: Supports Vanilla, Paper, Forge, Fabric, and NeoForge.
- **Automatic Version Detection**: Automatically fetches the latest loader/build version if not specified.
- **Smart Installation**:
  - Directly downloads ready-to-use JARs (Vanilla, Paper, Fabric).
  - Downloads the installer for modern Forge and NeoForge versions.
  - Automatically patches the vanilla server JAR for older Forge versions that use a patch file.
- **Easy to Use**: A simple and intuitive command-line interface.
- **Usable as a Go Library**: All functionalities are exported and can be used in your own Go projects.

## Installation

To install the command-line tool, use `go install`:

```shell
go install github.com/abulleDev/mcserverdl/cmd/msd@latest
```

## Usage

The main command is `msd`. It requires flags to specify the server type and game version.

```shell
msd -type <server_type> -game <game_version> [flags]
```

### Command-line Flags

| Flag      | Description                                                                                             | Required |
| :-------- | :------------------------------------------------------------------------------------------------------ | :------- |
| `-type`   | The type of server. Supported: `vanilla`, `paper`, `forge`, `fabric`, `neoforge`.                       | **Yes**  |
| `-game`   | The Minecraft game version (e.g., `1.21`).                                                              | **Yes**  |
| `-loader` | The version of the mod loader or the build number for Paper. Defaults to the latest version if omitted. | No       |
| `-path`   | The directory where the server will be installed. Defaults to the current directory (`.`).              | No       |

### Examples

```shell
# Download the latest Vanilla server for Minecraft 1.21 to the current directory.
msd -type vanilla -game 1.21

# Download Paper build 14 for Minecraft 1.21 into a folder named "my-paper-server".
msd -type paper -game 1.21 -loader 14 -path ./my-paper-server

# Download and automatically install the latest Forge server for Minecraft 1.21.6.
msd -type forge -game 1.21.6
```

## Library Usage

This project can also be used as a package in your own Go projects.

First, add the package to your project:

```shell
go get github.com/abulleDev/mcserverdl
```

Then, you can use the functions from the different packages to fetch versions, loaders, and download URLs.

### Example

This example demonstrates how to get the latest Paper build for a specific game version and then construct its download URL.

```go
package main

import (
    "fmt"
    "log"

    "github.com/abulleDev/mcserverdl/pkg/paper"
)

func main() {
    log.SetFlags(0)

    gameVersion := "1.21"
    fmt.Printf("--- Getting latest Paper server for %s ---\n", gameVersion)

    // 1. Get the list of all available builds for the game version.
    builds, err := paper.Builds(gameVersion, true)
    if err != nil {
        log.Fatalf("Failed to get Paper builds for %s: %v", gameVersion, err)
    }
    latestBuild := builds[0]
    fmt.Printf("Latest Paper build for %s: %d\n", gameVersion, latestBuild)

    // 2. Get the download URL for that specific build.
    downloadURL, err := paper.DownloadURL(gameVersion, latestBuild)
    if err != nil {
        log.Fatalf("Failed to get Paper download URL: %v", err)
    }
    fmt.Printf("Download URL: %s\n", downloadURL)
}
```

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
