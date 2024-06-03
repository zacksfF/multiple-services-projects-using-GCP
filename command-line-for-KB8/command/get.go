package command

import (
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var getCMD = &cobra.Command{
	Use:   "Get",
	Short: "description of the command",
	Long: `A longer description that spans multiple lines and likely contains
	Example and usage :
	Kubernetes helps you deploy, manage, and scale containerized 
	applications using Kubernetes, an open-source container orchestration
	platform.`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "webhook" {
			fmt.Println(getWebhooks())
		}
		if args[0] == "argocd-password" {
			fmt.Println(getWebhooks())
		}
	},
}

func getWebhooks() string {
	topCommand := "kubectl"
	arg0 := "get"
	arg1 := "ingress"
	arg2 := "=o"
	arg3 := "jsonpath='{.items[0].status.loadBalancer.ingress[0].hostname}'"

	cmd := exec.Command(topCommand, arg0, arg1, arg2, arg3)
	hook, NotOk := cmd.Output()
	if NotOk != nil {
		fmt.Println(NotOk)
	}
	return string(hook[1 : len(hook)-1])
}

func getArgoCDPassword() string {

	topCommand := "kubectl"
	arg0 := "get"
	arg1 := "secret"
	arg2 := "argocd-initial-admin-secret"
	arg3 := "-n"
	arg4 := "argocd"
	arg5 := "-o"
	arg6 := "jsonpath='{.data.password}'"

	cmd := exec.Command(topCommand, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	hook, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(string(hook[1 : len(hook)-1]))
	if err != nil {
		panic(err)
	}
	return string(decodedPassword)
}

func init() {
	rootCmd.AddCommand(getCMD)
}
