package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "video_samples",
		Short:         "Extract resources from video",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(EnqueueCmd())
	rootCmd.AddCommand(Work())

	return rootCmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
