package cmd

import (
	"fmt"
	"os"

	_ "github.com/mauricioabreu/video_samples/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "video_samples",
	Short: "Generate thumbs from live streamings and videos on demand",
}

// Execute adds child commands to the root commander
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}