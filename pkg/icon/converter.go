package icon

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// Sizes for different platforms
var (
	// BrowserExtSizes contains standard icon sizes for browser extensions
	BrowserExtSizes = []int{16, 32, 48, 128}

	// MacAppSizes contains standard icon sizes for macOS applications
	MacAppSizes = []int{16, 32, 64, 128, 256, 512, 1024}

	// WindowsAppSizes contains standard icon sizes for Windows applications
	WindowsAppSizes = []int{16, 32, 48, 64, 128, 256}

	// FaviconPNFSizes contains common PNG sizes for favicons. These will also be used for the .ico file.
	FaviconPNFSizes = []int{16, 32, 48, 64}

	// FaviconOtherPNFSizes contains additional PNG sizes to generate for modern browsers.
	FaviconOtherPNFSizes = []int{192, 512}

	// AppleTouchIconSize is the standard size for apple-touch-icon.png
	AppleTouchIconSize = 180
)

// Config contains configuration options for icon conversion
type Config struct {
	// InputPath is the path to the source image
	InputPath string

	// OutputDir is the directory where output files will be saved
	OutputDir string

	// ForceFlag indicates whether to overwrite existing files
	ForceFlag bool
}

// Converter handles image resizing and format conversion for various platforms
type Converter struct {
	// Config holds the converter configuration
	Config Config
}

// NewConverter creates a new icon converter with the provided configuration
func NewConverter(config Config) *Converter {
	return &Converter{
		Config: config,
	}
}

// resizeAndSaveInternal resizes an image and saves it to the output directory with a specific filename and format
func (c *Converter) resizeAndSaveInternal(outputDir string, size int, format imaging.Format, filename string) (string, error) {
	// Load the source image
	src, err := imaging.Open(c.Config.InputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image %s: %w", c.Config.InputPath, err)
	}

	// Resize the image using Lanczos resampling
	resized := imaging.Resize(src, size, size, imaging.Lanczos)

	// Determine file extension based on format
	var extension string
	switch format {
	case imaging.PNG:
		extension = ".png"
	case imaging.JPEG:
		extension = ".jpg"
	case imaging.GIF:
		extension = ".gif"
	case imaging.BMP:
		extension = ".bmp"
	case imaging.TIFF:
		extension = ".tiff"
	default:
		extension = ".png"
	}

	outPath := filepath.Join(outputDir, filename+extension)

	// Check if file exists and force flag is not set
	if _, err := os.Stat(outPath); err == nil && !c.Config.ForceFlag {
		return "", fmt.Errorf("output file already exists: %s. Use -f or --force flag to overwrite", outPath)
	}

	// Save the resized image
	if err = imaging.Save(resized, outPath); err != nil {
		return "", fmt.Errorf("failed to save resized image to %s: %w", outPath, err)
	}

	return outPath, nil
}

// ResizeAndSave resizes an image and saves it to the output directory with the specified format
// The filename will be in the format icon_WxH.ext
func (c *Converter) ResizeAndSave(outputDir string, size int, format imaging.Format) error {
	// Create output filename based on size
	outName := fmt.Sprintf("icon_%dx%d", size, size)
	_, err := c.resizeAndSaveInternal(outputDir, size, format, outName)
	return err
}

// CopyFile copies a file from src to dst, respecting the force flag
func (c *Converter) CopyFile(src, dst string) error {
	// Check if destination file exists and force flag is not set
	if _, err := os.Stat(dst); err == nil && !c.Config.ForceFlag {
		return fmt.Errorf("output file already exists: %s. Use -f or --force flag to overwrite", dst)
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return destFile.Sync()
}

// ProcessForBrowserExtension generates icons for browser extensions
func (c *Converter) ProcessForBrowserExtension() error {
	extensionDir := filepath.Join(c.Config.OutputDir, "browser-extension")
	if err := os.MkdirAll(extensionDir, 0755); err != nil {
		return fmt.Errorf("failed to create browser extension directory: %w", err)
	}

	// Process each size
	for _, size := range BrowserExtSizes {
		if err := c.ResizeAndSave(extensionDir, size, imaging.PNG); err != nil {
			return fmt.Errorf("failed to process browser extension icon size %dx%d: %w", size, size, err)
		}
	}

	fmt.Printf("Browser extension icons generated in: %s\n", extensionDir)
	return nil
}

// ProcessForMacApp generates icons for macOS applications
func (c *Converter) ProcessForMacApp() error {
	macDir := filepath.Join(c.Config.OutputDir, "mac-app")
	contentsDir := filepath.Join(macDir, "Contents")
	resourcesDir := filepath.Join(contentsDir, "Resources")

	// Create standard macOS app bundle structure
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		return fmt.Errorf("failed to create Mac app resources directory: %w", err)
	}

	// Create temporary iconset directory for iconutil
	iconsetDir := filepath.Join(macDir, "AppIcon.iconset")
	if err := os.MkdirAll(iconsetDir, 0755); err != nil {
		return fmt.Errorf("failed to create Mac app iconset directory: %w", err)
	}

	// Define macOS icon naming convention according to Apple Human Interface Guidelines
	iconConventions := []struct {
		size     int
		filename string
		scale    int // 1 for standard, 2 for @2x
	}{
		{16, "icon_16x16.png", 1},
		{32, "icon_16x16@2x.png", 2},
		{32, "icon_32x32.png", 1},
		{64, "icon_32x32@2x.png", 2},
		{128, "icon_128x128.png", 1},
		{256, "icon_128x128@2x.png", 2},
		{256, "icon_256x256.png", 1},
		{512, "icon_256x256@2x.png", 2},
		{512, "icon_512x512.png", 1},
		{1024, "icon_512x512@2x.png", 2},
	}

	// Generate iconset files
	for _, convention := range iconConventions {
		// Use resizeAndSaveInternal to get the path, though we don't need it here
		_, err := c.resizeAndSaveInternal(iconsetDir, convention.size, imaging.PNG, convention.filename)
		if err != nil {
			return fmt.Errorf("failed to save macOS icon %s: %w", convention.filename, err)
		}
	}

	// Define paths for .icns files
	tmpIcnsPath := filepath.Join(macDir, "AppIcon.icns")
	finalIcnsPath := filepath.Join(resourcesDir, "AppIcon.icns")

	// Try to automatically convert .iconset to .icns using iconutil
	fmt.Println("Attempting to convert .iconset to .icns using iconutil...")

	cmd := exec.Command("iconutil", "-c", "icns", "-o", tmpIcnsPath, iconsetDir)
	if err := cmd.Run(); err != nil {
		fmt.Println("Automatic conversion failed. Check if iconutil is installed and in your PATH.")
		fmt.Printf("To create .icns file manually, run: cd \"%s\" && iconutil -c icns AppIcon.iconset\n", macDir)
		fmt.Printf("Then place the resulting .icns file at: %s\n", finalIcnsPath)
		return nil // We don't return error here to allow users to convert manually
	}

	// If iconutil succeeded, move the .icns file to the Resources directory
	if err := os.Rename(tmpIcnsPath, finalIcnsPath); err != nil {
		fmt.Printf("Created .icns file but failed to move it to Resources directory: %v\n", err)
		fmt.Printf("Please manually move %s to %s\n", tmpIcnsPath, finalIcnsPath)
		return nil // We don't return error here to allow users to move manually
	}

	fmt.Printf("Successfully created and placed .icns file at: %s\n", finalIcnsPath)
	fmt.Printf("macOS app icons generated in: %s\n", iconsetDir)
	fmt.Println("Complete macOS app directory structure created at:", macDir)

	return nil
}

// ProcessForWindowsApp generates icons for Windows applications
func (c *Converter) ProcessForWindowsApp() error {
	winDir := filepath.Join(c.Config.OutputDir, "windows-app")
	if err := os.MkdirAll(winDir, 0755); err != nil {
		return fmt.Errorf("failed to create Windows app directory: %w", err)
	}

	// Process each size
	for _, size := range WindowsAppSizes {
		if err := c.ResizeAndSave(winDir, size, imaging.PNG); err != nil {
			return fmt.Errorf("failed to process Windows icon size %dx%d: %w", size, size, err)
		}
	}

	fmt.Printf("Windows app icons generated in: %s\n", winDir)
	fmt.Println("To create .ico file, use a third-party tool with these images")
	return nil
}

// ProcessForFavicon generates standard favicon files, including a multi-resolution .ico file if magick is available
func (c *Converter) ProcessForFavicon() error {
	faviconDir := filepath.Join(c.Config.OutputDir, "favicon")
	if err := os.MkdirAll(faviconDir, 0755); err != nil {
		return fmt.Errorf("failed to create favicon directory: %w", err)
	}

	var icoInputFiles []string

	// Process standard PNG sizes for .ico
	for _, size := range FaviconPNFSizes {
		filename := fmt.Sprintf("favicon-%dx%d", size, size)
		outPath, err := c.resizeAndSaveInternal(faviconDir, size, imaging.PNG, filename)
		if err != nil {
			return fmt.Errorf("failed to process favicon PNG size %dx%d: %w", size, size, err)
		}
		icoInputFiles = append(icoInputFiles, outPath)
	}

	// Process other PNG sizes
	for _, size := range FaviconOtherPNFSizes {
		filename := fmt.Sprintf("favicon-%dx%d", size, size)
		if _, err := c.resizeAndSaveInternal(faviconDir, size, imaging.PNG, filename); err != nil {
			return fmt.Errorf("failed to process favicon PNG size %dx%d: %w", size, size, err)
		}
	}

	// Process apple-touch-icon
	appleTouchFilename := "apple-touch-icon"
	if _, err := c.resizeAndSaveInternal(faviconDir, AppleTouchIconSize, imaging.PNG, appleTouchFilename); err != nil {
		return fmt.Errorf("failed to process apple-touch-icon: %w", err)
	}

	fmt.Printf("Favicon PNG files generated in: %s\n", faviconDir)

	// Attempt to generate .ico using magick command
	icoOutputPath := filepath.Join(faviconDir, "favicon.ico")
	fmt.Println("Attempting to generate favicon.ico using ImageMagick (magick command)...")

	// Check if file exists and force flag is not set
	if _, err := os.Stat(icoOutputPath); err == nil && !c.Config.ForceFlag {
		fmt.Printf("Skipping favicon.ico generation: %s already exists. Use -f to overwrite.\n", icoOutputPath)
		return nil
	}

	magickCmd := exec.Command("magick", append([]string{"convert"}, append(icoInputFiles, icoOutputPath)...)...)
	output, err := magickCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to generate favicon.ico using magick command: %v\n", err)
		fmt.Printf("ImageMagick output:\n%s\n", string(output))
		fmt.Println("Please ensure ImageMagick is installed and the 'magick' command is in your PATH.")
		fmt.Println("You may need to generate the favicon.ico manually from the generated PNGs.")
	} else {
		fmt.Printf("Successfully generated favicon.ico at: %s\n", icoOutputPath)
	}

	return nil
}
