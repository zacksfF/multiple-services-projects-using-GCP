package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// InstallCMD represent teh install command
var installCmd = &cobra.Command{
	Use:   "Install",
	Short: "Helps to install serverless project on kuberenets",
	Long: `Example:
	kb8lctl install --name=argocd
	kb8lctl install --name=knative
	kb8lctl install --name=tekton`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		defaultNamespace := make(map[string]string)
		defaultNamespace["tekton"] = "tekton-pipelines"
		defaultNamespace["knative"] = "knative-serving"
		defaultNamespace["argocd"] = "argocd"

		if kubectlPresentCheck() {
			fmt.Println("kubectl is installed.")
			fmt.Println("Installing project:", name, "in namescape", defaultNamespace[name])
			installProject(name, defaultNamespace[name])
		} else {
			fmt.Println("kubectl is not installed. Please try again")
			os.Exit(1)
		}
	},
}

func kubectlPresentCheck() bool {
	topCammand := "kubectl"
	arg0 := "version"

	cmd := exec.Command(topCammand, arg0)
	_, err := cmd.Output()
	return err == nil
}

func installProject(name string, namespace string) {
	switch name {
	case "argocd":
		installArgoCD(namespace)
	case "knative":
		installKnative(namespace)
	case "tekton":
		installTektonPipelines(namespace)
	case "all":
		installArgoCD(namespace)
		installKnative(namespace)
		installTektonPipelines(namespace)
	default:
		fmt.Println("Tool not found.")
	}
}

func installArgoCD(namespace string) {
	fmt.Println("Installing ArgoCD Started")
	fmt.Println("Creating namespace:", namespace)
	createNamespace(namespace)
	fmt.Println("Namespace created.")
	c := exec.Command("bash")
	c.Stdin = strings.NewReader(argocdInstallationScript)
	b, e := c.Output()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(string(b))
	fmt.Println("ArgoCD installation complete.")
}

func installKnative(namespace string) {
	fmt.Println("Installing Knative Serving Started...")
	c := exec.Command("bash")
	c.Stdin = strings.NewReader(knativeInstallationScript)
	b, e := c.Output()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(string(b))
	fmt.Println("Knative installation complete.")
}

func installTektonPipelines(namespace string) {
	fmt.Println("Installing TektonPipelines Started...")
	c := exec.Command("bash")
	c.Stdin = strings.NewReader(tektonInstallationScript)
	b, e := c.Output()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(string(b))
	fmt.Println("TektonPipelines installation complete.")
}

func createNamespace(namespace string) {
	topCommand := "kubectl"
	arg0 := "create"
	arg1 := "namespace"
	arg2 := namespace
	// temp := "version"

	cmd := exec.Command(topCommand, arg0, arg1, arg2)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
}


////////
var knativeInstallationScript = `
#!/bin/bash

echo "Installing Knative Serving CRDs..."
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.3.0/serving-crds.yaml
echo "Finished Installing Knative Serving CRDs"
echo "Installing core components of Knative Serving..."
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.3.0/serving-core.yaml
echo "Finished Installing core components of Knative Serving"
echo "Installing networking layer (Knative Kourier controller)..."
kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.3.0/kourier.yaml
echo "Finished Installing networking layer"
echo "Patching networking layer settings..."
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
echo "Installing Magic DNS (sslip.io)..."
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.3.0/serving-default-domain.yaml
echo "Finished Installing Magic DNS"
`

var tektonInstallationScript = `
#!/bin/bash

echo "Installing Tekton Pipelines..."
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.26.0/release.yaml
echo "Finished Installing Tekton Pipelines"
echo "Installing Tekton Dashboard..."
kubectl apply --filename https://storage.googleapis.com/tekton-releases/dashboard/latest/tekton-dashboard-release.yaml
echo "Finished Installing Tekton Dashboard"
echo "Installing Tekton Triggers..."
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
echo "Finished Installing Tekton Triggers"
echo "Installing Tekton Core Interceptors..."
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml
echo "Finished Installing Tekton Core Interceptors"
echo "Installing Nginx-Ingress Controller..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/aws/deploy.yaml
echo "Finished Installing Nginx-Ingress Controller"
`

var argocdInstallationScript = `
#!/bin/bash

echo "Installing ArgoCD..."
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
echo "Finished Installing ArgoCD"
`