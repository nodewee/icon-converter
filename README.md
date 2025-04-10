# Icon Converter

A command-line tool to convert icons to various formats and sizes for different platforms like browser extensions, macOS applications, and Windows applications.

## Features

- Convert icons for browser extensions in sizes: 16x16, 32x32, 48x48, 128x128
- Convert icons for macOS applications with proper iconset structure
- Generate .icns files for macOS applications (requires iconutil)
- Convert icons for Windows applications in various sizes
- Force overwrite of existing files with `-f` flag

## Installation

Make sure you have Go installed (version 1.18 or later), then run:

```bash
git clone https://github.com/yourusername/icon-converter.git
cd icon-converter
go build
```

## Usage

```bash
# Basic usage
./icon-converter [input image] [output directory] [flags]

# Flags
-b, --browser-extension   Convert for browser extension requirements
-m, --mac-app             Convert for macOS application requirements
-w, --windows-app         Convert for Windows application requirements
-f, --force               Force overwrite existing files
-h, --help                Help for the command

# Examples
./icon-converter icon.png ./output -b -m -w  # Convert for all platforms
./icon-converter icon.png ./output -b        # Convert for browser extensions only
./icon-converter icon.png ./output -m -f     # Convert for macOS and force overwrite
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

- For macOS .icns file generation, the `iconutil` command-line tool is required (included with macOS)

## License

MIT 