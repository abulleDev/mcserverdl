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
go install github.com/abulleDev/mcserverdl/v2/cmd/mcserverdl@latest
```

## Usage

The main command is `mcserverdl`. It requires flags to specify the server type and game version.

```shell
mcserverdl -type <server_type> -game <game_version> [flags]
```

### Command-line Flags

| Flag       | Description                                                                                             | Required |
| :--------- | :------------------------------------------------------------------------------------------------------ | :------- |
| `-type`    | The type of server. Supported: `vanilla`, `paper`, `forge`, `fabric`, `neoforge`.                       | **Yes**  |
| `-game`    | The Minecraft game version (e.g., `1.21`).                                                              | **Yes**  |
| `-server`  | The version of the mod loader or the build number for Paper. Defaults to the latest version if omitted. | No       |
| `-path`    | The directory where the server will be installed. Defaults to the current directory (`.`).              | No       |
| `-version` | Prints the current version of the tool.                                                                 | No       |

### Examples

```shell
# Download the latest Vanilla server for Minecraft 1.21 to the current directory.
mcserverdl -type vanilla -game 1.21

# Download Paper build 14 for Minecraft 1.21 into a folder named "my-paper-server".
mcserverdl -type paper -game 1.21 -server 14 -path ./my-paper-server

# Download and automatically install the latest NeoForge server for Minecraft 1.21.6.
mcserverdl -type neoforge -game 1.21.6
```

## Library Usage

This project can also be used as a package in your own Go projects.

First, add the package to your project:

```shell
go get github.com/abulleDev/mcserverdl/v2
```

Then, you can use the functions from the different packages to fetch versions, loaders, and download URLs.

### Example

This example demonstrates how to use the factory to get a provider, fetch the available versions, and get the download URL.

```go
package main

import (
	"fmt"
	"log"

	"github.com/abulleDev/mcserverdl/v2/pkg/factory"
)

func main() {
	log.SetFlags(0)

	serverType := "paper" // Can be vanilla, paper, fabric, forge, neoforge
	gameVersion := "1.21"

	// 1. Initialize the provider using the factory.
	p, err := factory.New(serverType)
	if err != nil {
		log.Fatalf("Error creating provider: %v", err)
	}

	fmt.Printf("--- Getting latest %s server for %s ---\n", serverType, gameVersion)

	// 2. Get the list of all available server versions (builds/loaders) for the game version.
	// Note: Vanilla server typically determines the version from the game version itself.
	versions, err := p.ServerVersions(gameVersion)
	if err != nil {
		log.Fatalf("Failed to get server versions for %s: %v", gameVersion, err)
	}

	if len(versions) == 0 {
		log.Fatal("No server versions found")
	}

	// Assuming the first version is the latest (provider dependent)
	latestVersion := versions[0]
	fmt.Printf("Latest version matches: %s\n", latestVersion)

	// 3. Get the download URL for that specific version.
	downloadURL, err := p.DownloadURL(gameVersion, latestVersion)
	if err != nil {
		log.Fatalf("Failed to get download URL: %v", err)
	}
	fmt.Printf("Download URL: %s\n", downloadURL)
}
```

### Direct Package Usage

If you only need a specific server type, you can import the provider package directly to reduce dependencies or for simpler usage.

```go
package main

import (
	"fmt"
	"log"

	"github.com/abulleDev/mcserverdl/v2/pkg/provider/forge"
)

func main() {
	// Directly initialize the Forge provider
	provider := forge.New()

	// Download specific version directly
	// Arguments: Game Version, Loader Version, Install Path, Progress Callback (nil here)
	err := provider.Download("1.5.1", "7.7.2.682", "./server", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Forge server installed successfully!")
}
```

### Custom Logging

You can inject a custom logger (or the standard one) to see internal logs from the provider, such as fetching status or debug info.

```go
package main

import (
	"log"
	"os"

	"github.com/abulleDev/mcserverdl/v2/pkg/factory"
)

func main() {
	// Create a standard logger
	logger := log.New(os.Stdout, "[MC-DL] ", log.Ltime)

	p, _ := factory.New("fabric")

	// Inject the logger into the provider
	p.SetLogger(logger)

	// Now operations will log their progress
	// Output: [MC-DL] 01:06:37 Fetching Fabric server versions (loaders) for 1.20.1...
	p.ServerVersions("1.20.1")
}
```

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
