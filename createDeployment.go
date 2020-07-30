package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func createDeployment() string {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	replicas := int32(1)
	name := fmt.Sprintf("notebook-%s", uuid.New())
	name = "notebook-123"
	//volumeName := "vol"
	namespace := "testn"
	image := "jupyter/minimal-notebook"
	cpu := "100m"
	memory := "500Mi"
	//storageClassName := "local-path"
	storageSize := "2Gi"

	//createPVC(clientset, name, namespace)
	statefulSetClient := clientset.AppsV1().StatefulSets(namespace)
	// deploymentsClient := clientset.AppsV1().Deployments(namespace)

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Env: []corev1.EnvVar{
								{
									Name:  "NB_PREFIX",
									Value: fmt.Sprintf("/%s", name),
								},
							},
							Command: []string{"sh", "-c", "jupyter notebook --notebook-dir=/home/jovyan --ip=0.0.0.0 --no-browser --allow-root --port=8888 --NotebookApp.token='' --NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}"},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(cpu),
									corev1.ResourceMemory: resource.MustParse(memory),
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      name,
									MountPath: "/home/jovyan",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8888,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						//StorageClassName: &storageClassName,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								v1.ResourceName(v1.ResourceStorage): resource.MustParse(storageSize),
							},
						},
					},
				},
			},
		},
	}

	// deployment := &appsv1.Deployment{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      name,
	// 		Namespace: namespace,
	// 	},
	// 	Spec: appsv1.DeploymentSpec{
	// 		Replicas: &replicas,
	// 		Selector: &metav1.LabelSelector{
	// 			MatchLabels: map[string]string{
	// 				"app":     name,
	// 				"version": "v1",
	// 			},
	// 		},
	// 		Template: corev1.PodTemplateSpec{
	// 			ObjectMeta: metav1.ObjectMeta{
	// 				Labels: map[string]string{
	// 					"app":     name,
	// 					"version": "v1",
	// 				},
	// 			},
	// 			Spec: corev1.PodSpec{
	// 				Containers: []corev1.Container{
	// 					{
	// 						Name:            name,
	// 						Image:           image,
	// 						ImagePullPolicy: corev1.PullIfNotPresent,
	// 						Env: []corev1.EnvVar{
	// 							{
	// 								Name:  "NB_PREFIX",
	// 								Value: fmt.Sprintf("/%s", name),
	// 							},
	// 						},
	// 						Command: []string{"sh", "-c", "jupyter notebook --notebook-dir=/home/jovyan --ip=0.0.0.0 --no-browser --allow-root --port=8888 --NotebookApp.token='' --NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}"},
	// 						Resources: corev1.ResourceRequirements{
	// 							Requests: corev1.ResourceList{
	// 								corev1.ResourceCPU:    resource.MustParse(cpu),
	// 								corev1.ResourceMemory: resource.MustParse(memory),
	// 							},
	// 						},
	// 						VolumeMounts: []v1.VolumeMount{
	// 							{
	// 								Name:      name,
	// 								MountPath: "/home/jovyan",
	// 							},
	// 						},
	// 						Ports: []corev1.ContainerPort{
	// 							{
	// 								Name:          "http",
	// 								Protocol:      corev1.ProtocolTCP,
	// 								ContainerPort: 8888,
	// 							},
	// 						},
	// 					},
	// 				},
	// 				Volumes: []corev1.Volume{
	// 					{
	// 						Name: name,
	// 						VolumeSource: corev1.VolumeSource{
	// 							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	// 								ClaimName: name,
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// // Create Deployment
	// fmt.Println("Creating deployment...")
	// result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// deploymentName := result.GetObjectMeta().GetName()
	// fmt.Printf("Created deployment %q.\n", deploymentName)

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := statefulSetClient.Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	statefulSetName := result.GetObjectMeta().GetName()
	fmt.Printf("Created deployment %q.\n", statefulSetName)

	return statefulSetName
}

func createPVC(clientset *kubernetes.Clientset, name, namespace string) {
	storageClassName := "local-path"
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			StorageClassName: &storageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse("2Gi"),
				},
			},
		},
	}
	pvcClient := clientset.CoreV1().PersistentVolumeClaims(namespace)
	// Create PVC
	fmt.Println("Creating pvc...")
	result, err := pvcClient.Create(context.TODO(), pvc, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created pvc %q.\n", result.GetObjectMeta().GetName())
}
