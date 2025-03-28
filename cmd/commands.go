package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
)

var (
	browserExtFlag bool
	macAppFlag     bool
	windowsAppFlag bool
	forceFlag      bool

	// Browser extension icon sizes
	browserExtSizes = []int{16, 32, 48, 128}

	// macOS app icon sizes (all sizes needed for iconset)
	macAppSizes = []int{16, 32, 64, 128, 256, 512, 1024}

	// Windows app icon sizes
	windowsAppSizes = []int{16, 32, 48, 64, 128, 256}
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "icon [input image] [output directory]",
	Short: "Convert icons to various formats and sizes",
	Long: `Icon Converter is a command line tool to convert icons to various formats and sizes
for different platforms like browser extensions, macOS applications, and Windows applications.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputDir := args[1]

		// Check if input file exists
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", inputPath)
		}

		// Check if any output type is specified
		if !browserExtFlag && !macAppFlag && !windowsAppFlag {
			fmt.Println("警告: 未指定任何输出类型，未执行任何操作")
			fmt.Println("提示: 使用 -b, -m, -w 指定输出类型，或使用 --help 查看帮助")
			return nil
		}

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}

		// Process based on flags
		if browserExtFlag {
			if err := processForBrowserExtension(inputPath, outputDir); err != nil {
				return err
			}
		}

		if macAppFlag {
			if err := processForMacApp(inputPath, outputDir); err != nil {
				return err
			}
		}

		if windowsAppFlag {
			if err := processForWindowsApp(inputPath, outputDir); err != nil {
				return err
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Define flags
	rootCmd.Flags().BoolVarP(&browserExtFlag, "browser-extension", "b", false, "Convert for browser extension requirements")
	rootCmd.Flags().BoolVarP(&macAppFlag, "mac-app", "m", false, "Convert for macOS application requirements")
	rootCmd.Flags().BoolVarP(&windowsAppFlag, "windows-app", "w", false, "Convert for Windows application requirements")
	rootCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force overwrite existing files")
}

// resizeAndSave resizes an image and saves it to the output directory
func resizeAndSave(inputPath, outputDir string, size int, format imaging.Format) error {
	// Load the image
	src, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}

	// Resize the image
	resized := imaging.Resize(src, size, size, imaging.Lanczos)

	// Create output filename
	outName := fmt.Sprintf("icon_%dx%d", size, size)
	var outPath string
	switch format {
	case imaging.PNG:
		outPath = filepath.Join(outputDir, outName+".png")
	case imaging.JPEG:
		outPath = filepath.Join(outputDir, outName+".jpg")
	case imaging.GIF:
		outPath = filepath.Join(outputDir, outName+".gif")
	case imaging.BMP:
		outPath = filepath.Join(outputDir, outName+".bmp")
	case imaging.TIFF:
		outPath = filepath.Join(outputDir, outName+".tiff")
	default:
		outPath = filepath.Join(outputDir, outName+".png")
	}

	// Check if file exists and force flag is not set
	if _, err := os.Stat(outPath); err == nil && !forceFlag {
		return fmt.Errorf("output file already exists: %s. Use -f or --force flag to overwrite", outPath)
	}

	// Save the resized image
	err = imaging.Save(resized, outPath)
	if err != nil {
		return fmt.Errorf("failed to save resized image: %v", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Check if destination file exists and force flag is not set
	if _, err := os.Stat(dst); err == nil && !forceFlag {
		return fmt.Errorf("output file already exists: %s. Use -f or --force flag to overwrite", dst)
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// processForBrowserExtension processes an image for browser extension requirements
func processForBrowserExtension(inputPath, outputDir string) error {
	extensionDir := filepath.Join(outputDir, "browser-extension")
	if err := os.MkdirAll(extensionDir, 0755); err != nil {
		return fmt.Errorf("failed to create browser extension directory: %v", err)
	}

	// Process each size
	for _, size := range browserExtSizes {
		if err := resizeAndSave(inputPath, extensionDir, size, imaging.PNG); err != nil {
			return err
		}
	}

	fmt.Printf("Browser extension icons generated in: %s\n", extensionDir)
	return nil
}

// processForMacApp processes an image for macOS application requirements
func processForMacApp(inputPath, outputDir string) error {
	macDir := filepath.Join(outputDir, "mac-app")
	contentsDir := filepath.Join(macDir, "Contents")
	resourcesDir := filepath.Join(contentsDir, "Resources")

	// Create standard macOS app bundle structure
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		return fmt.Errorf("failed to create Mac app resources directory: %v", err)
	}

	// Create temporary iconset directory for iconutil
	iconsetDir := filepath.Join(macDir, "AppIcon.iconset")
	if err := os.MkdirAll(iconsetDir, 0755); err != nil {
		return fmt.Errorf("failed to create Mac app iconset directory: %v", err)
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
		// Load the image
		src, err := imaging.Open(inputPath)
		if err != nil {
			return fmt.Errorf("failed to open image: %v", err)
		}

		// Resize the image
		resized := imaging.Resize(src, convention.size, convention.size, imaging.Lanczos)

		// Create output path
		outPath := filepath.Join(iconsetDir, convention.filename)

		// Check if file exists and force flag is not set
		if _, err := os.Stat(outPath); err == nil && !forceFlag {
			return fmt.Errorf("output file already exists: %s. Use -f or --force flag to overwrite", outPath)
		}

		// Save the resized image
		err = imaging.Save(resized, outPath)
		if err != nil {
			return fmt.Errorf("failed to save resized image: %v", err)
		}
	}

	// Define paths for .icns files
	tmpIcnsPath := filepath.Join(macDir, "AppIcon.icns")
	finalIcnsPath := filepath.Join(resourcesDir, "AppIcon.icns")

	// Try to automatically convert .iconset to .icns using iconutil
	fmt.Println("Attempting to convert .iconset to .icns using iconutil...")

	cmd := exec.Command("iconutil", "-c", "icns", "-o", tmpIcnsPath, iconsetDir)
	if err := cmd.Run(); err != nil {
		fmt.Println("Automatic conversion failed. You will need to convert manually.")
		fmt.Printf("To create .icns file, run: cd \"%s\" && iconutil -c icns AppIcon.iconset\n", macDir)
		fmt.Printf("Then place the resulting .icns file at: %s\n", finalIcnsPath)
	} else {
		// If iconutil succeeded, move the .icns file to the Resources directory
		if err := os.Rename(tmpIcnsPath, finalIcnsPath); err != nil {
			fmt.Printf("Created .icns file but failed to move it to Resources directory: %v\n", err)
			fmt.Printf("Please manually move %s to %s\n", tmpIcnsPath, finalIcnsPath)
		} else {
			fmt.Printf("Successfully created and placed .icns file at: %s\n", finalIcnsPath)
		}
	}

	fmt.Printf("macOS app icons generated in: %s\n", iconsetDir)
	fmt.Println("Complete macOS app directory structure created at:", macDir)

	return nil
}

// processForWindowsApp processes an image for Windows application requirements
func processForWindowsApp(inputPath, outputDir string) error {
	winDir := filepath.Join(outputDir, "windows-app")
	if err := os.MkdirAll(winDir, 0755); err != nil {
		return fmt.Errorf("failed to create Windows app directory: %v", err)
	}

	// Process each size
	for _, size := range windowsAppSizes {
		if err := resizeAndSave(inputPath, winDir, size, imaging.PNG); err != nil {
			return err
		}
	}

	fmt.Printf("Windows app icons generated in: %s\n", winDir)
	fmt.Println("To create .ico file, use a third-party tool with these images")
	return nil
}
