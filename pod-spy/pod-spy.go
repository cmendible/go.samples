package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	// get K8s Configuration
	config, checkEvery := getConfig()

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		for _, namespace := range getNamespaces() {
			pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
	
			for _, pod := range pods.Items {
				for _, label := range getRequiredLabels(){
					_, present := pod.Labels[label]
					if !present {
						fmt.Printf("Pod %v does not have the %v label\n", pod.Name, label)	
					}
				}
			}
		}

		time.Sleep(time.Duration(*checkEvery) * time.Second)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getNamespaces() []string {
	return []string {""}
}

func getRequiredLabels() []string {
	return []string {"app"}
}

func getConfig() (config *rest.Config, checkEvery *int) {
	// Check if running outside K8s
	clusterMode := flag.Bool("clustermode", true, "(optional) run in cluster mode")
	checkEvery = flag.Int("checkevery", 60, "(optional) run checks every X seconds")
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	
	if *clusterMode {
		// creates the in-cluster config
		conf, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		config = conf
	} else {
		// use the current context in kubeconfig
		conf, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		config = conf
	}

	return config, checkEvery
}
