package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createVirtualService(name string) {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	ic, err := versionedclient.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}

	virtualServicesClient := ic.NetworkingV1alpha3().VirtualServices("notebooks")

	vs := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "notebooks",
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    []string{"*"},
			Gateways: []string{"default/gateway"},
			Http: []*networkingv1alpha3.HTTPRoute{
				{
					Match: []*networkingv1alpha3.HTTPMatchRequest{
						{
							Uri: &networkingv1alpha3.StringMatch{
								MatchType: &networkingv1alpha3.StringMatch_Prefix{
									Prefix: fmt.Sprintf("/%s/", name),
								},
							},
						},
					},
					Rewrite: &networkingv1alpha3.HTTPRewrite{
						Uri: fmt.Sprintf("/%s/", name),
					},
					Route: []*networkingv1alpha3.HTTPRouteDestination{
						{
							Destination: &networkingv1alpha3.Destination{
								Port: &networkingv1alpha3.PortSelector{
									Number: 8888,
								},
								Host: name,
							},
						},
					},
				},
			},
		},
	}

	// Create VirtualService
	fmt.Println("Creating virtual service...")

	result, err := virtualServicesClient.Create(context.TODO(), vs, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created virtual service %q.\n", result.GetObjectMeta().GetName())
}
