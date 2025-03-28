# Icon Converter

A command-line tool to convert and resize icons for different platforms.

## Features

- Convert icons for browser extensions
- Convert icons for macOS applications 
- Convert icons for Windows applications
- Automatically creates output directories if they don't exist
- Checks for existing files and offers force overwrite option

## Installation

```
go install github.com/nodewee/icon-converter@latest
```

Or you can build from source:

```
git clone https://github.com/nodewee/icon-converter.git
cd icon-converter
go build -o icon
```

## Usage

```
icon [input image] [output directory] [flags]
```

### Required Arguments

- `input image`: Path to the source image file
- `output directory`: Path to the directory where the converted icons will be saved

### Flags

- `--browser-extension`, `-b`: Generate icons for browser extensions (16x16, 32x32, 48x48, 128x128)
- `--mac-app`, `-m`: Generate icons for macOS applications (16x16 to 1024x1024, .iconset format)
- `--windows-app`, `-w`: Generate icons for Windows applications (16x16 to 256x256)
- `--force`, `-f`: Force overwrite existing files (without this flag, the tool will exit with error if output files already exist)

### Behavior

- At least one output type flag (`-b`, `-m`, or `-w`) must be specified, otherwise the tool will display a warning and exit without any action.
- If output files already exist, the tool will exit with an error unless the `-f` flag is used.

## Examples

Convert an image for browser extension:

```
icon logo.png output-icons -b
```

Convert an image for both macOS and Windows applications:

```
icon myicon.png app-icons -m -w
```

Force overwrite existing files:

```
icon logo.png existing-icons -b -f
```

## Requirements

- Go 1.16 or higher

## License

MIT 