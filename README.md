# Icon Converter

A command-line tool to convert icons to various formats and sizes for different platforms like browser extensions, macOS applications, Windows applications, and websites (favicons).

## Features

- Convert icons for browser extensions in sizes: 16x16, 32x32, 48x48, 128x128
- Convert icons for macOS applications with proper iconset structure
- Generate .icns files for macOS applications (requires `iconutil`)
- Convert icons for Windows applications in various sizes
- Generate standard favicon PNG files (favicon-16x16.png, favicon-32x32.png, etc.) and apple-touch-icon.png
- Generate a multi-resolution `favicon.ico` file (requires ImageMagick `magick` command)
- Force overwrite of existing files with `-f` flag

## Installation

Make sure you have Go installed (version 1.18 or later), then run:

```bash
git clone https://github.com/yourusername/icon-converter.git
cd icon-converter
go build
```

Additionally, for generating `.ico` files, you need to install ImageMagick: [https://imagemagick.org/script/download.php](https://imagemagick.org/script/download.php)
Ensure the `magick` command is available in your system's PATH.

## Usage

```bash
# Basic usage
./icon-converter [input image] [output directory] [flags]

# Flags
-b, --browser-extension   Convert for browser extension requirements
-m, --mac-app             Convert for macOS application requirements
-w, --windows-app         Convert for Windows application requirements
   --favicon              Convert for website favicon requirements
   --fav                  Alias for --favicon
-f, --force               Force overwrite existing files
-h, --help                Help for the command

# Examples
./icon-converter icon.png ./output -b -m -w --favicon  # Convert for all platforms
./icon-converter icon.png ./output -b                # Convert for browser extensions only
./icon-converter icon.png ./output -m -f             # Convert for macOS and force overwrite
./icon-converter icon.png ./output --fav             # Convert for website favicons only (using alias)
```

## Code Architecture

The project follows a clean architecture approach with separation of concerns:

- `cmd/`: Contains the command-line interface logic
  - `commands.go`: Defines the command-line flags and calls the icon package

- `pkg/icon/`: Contains the core icon conversion logic
  - `converter.go`: Implements the icon conversion functionality

## Dependencies

- [github.com/disintegration/imaging](https://github.com/disintegration/imaging): Image processing library
- [github.com/spf13/cobra](https://github.com/spf13/cobra): Command-line interface library

## Requirements

- Go (version 1.18 or later)
- For macOS `.icns` file generation: `iconutil` command-line tool (included with macOS)
- For `.ico` file generation (using `--favicon` or `--fav`): ImageMagick (`magick` command in PATH)

## Automated Builds

This project uses GitHub Actions to automatically build binaries for multiple platforms:

- When pushing to the `main` or `master` branch, workflow builds binaries for:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)

- When creating a tag with format `v*` (e.g., `v1.0.0`), the workflow:
  1. Builds all platform binaries
  2. Creates a GitHub Release
  3. Attaches all binaries to the release

To create a new release:
```bash
git tag v1.0.0  # Replace with your version
git push origin v1.0.0
```

## License

MIT 