#!/bin/zsh
BIN_FILENAME=icon

DEST_PATH=~/go/bin/$BIN_FILENAME

# build the binary
echo "Building the binary..."
go build -o $DEST_PATH

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful! Binary installed at $DEST_PATH"
    echo "You can now use '$BIN_FILENAME' command."
else
    echo "Build failed!"
    exit 1
fi
