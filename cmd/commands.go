package cmd

import (
	"fmt"
	"os"

	"github.com/nodewee/icon-converter/pkg/icon"

	"github.com/spf13/cobra"
)

// Command line flags
var (
	browserExtFlag bool // Generate browser extension icons
	macAppFlag     bool // Generate macOS application icons
	windowsAppFlag bool // Generate Windows application icons
	faviconFlag    bool // Generate favicon files
	overwriteFlag  bool // Force overwrite existing files
)

// rootCmd represents the base command for the icon converter application
var rootCmd = &cobra.Command{
	Use:   "icon [input image] [output directory]",
	Short: "Convert icons to various formats and sizes",
	Long: `Icon Converter is a command line tool to convert icons to various formats and sizes
for different platforms like browser extensions, macOS applications, Windows applications, and favicons.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputDir := args[1]

		// Check if input file exists
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", inputPath)
		}

		// Check if any output type is specified
		if !browserExtFlag && !macAppFlag && !windowsAppFlag && !faviconFlag {
			fmt.Println("Warning: No output type specified, no action performed")
			fmt.Println("Tip: Use -b, -m, -w, or -f to specify output type, or use --help to view help")
			return nil
		}

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}

		// Create converter configuration
		config := icon.Config{
			InputPath:     inputPath,
			OutputDir:     outputDir,
			OverwriteFlag: overwriteFlag,
		}

		// Create converter instance
		converter := icon.NewConverter(config)

		// Process based on specified flags
		if browserExtFlag {
			if err := converter.ProcessForBrowserExtension(); err != nil {
				return err
			}
		}

		if macAppFlag {
			if err := converter.ProcessForMacApp(); err != nil {
				return err
			}
		}

		if windowsAppFlag {
			if err := converter.ProcessForWindowsApp(); err != nil {
				return err
			}
		}

		if faviconFlag {
			if err := converter.ProcessForFavicon(); err != nil {
				return err
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once.
func Execute() error {
	return rootCmd.Execute()
}

// init registers command line flags
func init() {
	// Define command line flags
	rootCmd.Flags().BoolVarP(&browserExtFlag, "browser-extension", "b", false, "Convert for browser extension requirements")
	rootCmd.Flags().BoolVarP(&macAppFlag, "mac-app", "m", false, "Convert for macOS application requirements")
	rootCmd.Flags().BoolVarP(&windowsAppFlag, "windows-app", "w", false, "Convert for Windows application requirements")
	rootCmd.Flags().BoolVarP(&faviconFlag, "favicon", "f", false, "Convert for website favicon requirements")
	rootCmd.Flags().BoolVar(&overwriteFlag, "overwrite", false, "Overwrite existing files")
}
