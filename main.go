package main

import "fmt"

func main() {
	fmt.Println("Hello")
	deploymentName := createDeployment()
	fmt.Println(deploymentName)
	//serviceName := createService(deploymentName)
	//createVirtualService(serviceName)
	//deleteDeployment("notebook-0b9e383f-1b3c-43f9-b61e-fa0c66e5f530", "notebooks")
}
