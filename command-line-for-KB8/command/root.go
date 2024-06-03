package command

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kb8lctl",
	Short: "Unifying command-line interface for interacting with the Kubernetes Serverless Stack",
	Long:  `The Kubernetes Serverless CLI (kslctl) Project provides a unifying command-line interface for interacting with the Kubernetes Serverless Stack which constitutes of Knative, Tekton and ArgoCD.`,
	// Uncomment the following line if your bare application
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
