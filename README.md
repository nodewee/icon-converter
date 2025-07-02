# Icon Converter

A command-line tool to convert **any image** into various icon formats and sizes for different platforms, such as browser extensions, macOS applications, Windows applications, and websites (favicons).

## Features

- Convert icons for browser extensions in sizes: 16x16, 32x32, 48x48, 128x128
- Convert icons for macOS applications with proper iconset structure
- Generate .icns files for macOS applications (requires `iconutil`)
- Convert icons for Windows applications in various sizes
- Generate standard favicon PNG files (favicon-16x16.png, favicon-32x32.png, etc.) and apple-touch-icon.png
- Generate a multi-resolution `favicon.ico` file (requires ImageMagick `magick` command)
- Overwrite existing files with `--overwrite` flag

## Installation

Make sure you have Go installed (version 1.18 or later), then run:

```bash
git clone https://github.com/nodewee/icon-converter.git
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
-f, --favicon             Convert for website favicon requirements
    --overwrite           Overwrite existing files
-h, --help                Help for the command

# Examples
./icon-converter icon.png ./output -b -m -w -f                # Convert for all platforms
./icon-converter icon.png ./output -b                         # Convert for browser extensions only
./icon-converter icon.png ./output -m --overwrite             # Convert for macOS and overwrite existing files
./icon-converter icon.png ./output -f                         # Convert for website favicons only
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
- For `.ico` file generation (using `--favicon`): ImageMagick (`magick` command in PATH)

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

## Example Results

Here's what you can expect after running the converter on an input image:

### Browser Extension Icons
```
output/
├── icon-16.png  (16x16)
├── icon-32.png  (32x32)
├── icon-48.png  (48x48)
└── icon-128.png (128x128)
```

### macOS Application Icons
```
output/
└── AppIcon.iconset/
    ├── icon_16x16.png
    ├── icon_16x16@2x.png
    ├── icon_32x32.png
    ├── icon_32x32@2x.png
    ├── icon_128x128.png
    ├── icon_128x128@2x.png
    ├── icon_256x256.png
    ├── icon_256x256@2x.png
    ├── icon_512x512.png
    └── icon_512x512@2x.png
└── AppIcon.icns (if iconutil is available)
```

### Windows Application Icons
```
output/
├── icon-16.ico
├── icon-32.ico
├── icon-48.ico
└── icon-256.ico
```

### Website Favicons
```
output/
├── favicon-16x16.png
├── favicon-32x32.png
├── favicon-48x48.png
├── favicon.ico (multi-resolution if ImageMagick is available)
└── apple-touch-icon.png (180x180)
```

## License

MIT 