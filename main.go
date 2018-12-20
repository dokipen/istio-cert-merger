package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/dokipen/istio-cert-merger/lib"
)

func getClientset() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func main() {
	bind := os.Getenv("BIND")
	if bind == "" {
		bind = ":8080"
	}
	s := &http.Server{
		Addr:    bind,
		Handler: &istiocertmerger.Server{getClientset()},
	}
	logrus.WithFields(logrus.Fields{"bind": bind}).Info("Listening")
	logrus.WithFields(logrus.Fields{"err": s.ListenAndServe().Error()}).Fatal("Exited")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
