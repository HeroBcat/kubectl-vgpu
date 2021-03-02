package main

import (
	"fmt"
	"log"
)

func init() {
	kubeInit()
}

func main() {

	for _, node := range getNodes() {

		fmt.Printf("NODE NAME: %s\n", node.Name)
		fmt.Printf("NODE IP: %s\n", node.InternalIP)
		fmt.Printf("GPU COUNT: %d\n", node.GPUCount)
		fmt.Printf("GPU MEMORY (Mib): %d * %d\n", node.EachGpuMemory, node.GPUCount)

		pods, err := getActivePodsByNode(node)
		if err != nil {
			log.Fatal(err)
		}

		usingGPUCore := make([]int64, node.GPUCount)
		usingGPUMemory := make([]int64, node.GPUCount)

		for _, pod := range pods {
			fmt.Println("")
			fmt.Printf("POD NAME: %s\n", pod.Name)
			fmt.Printf("POD NAMESPACE: %s\n", pod.NameSpace)
			fmt.Printf("REQUEST GPU CORE: %d\n", pod.RequestGPUCore)
			fmt.Printf("REUQEST GPU MEMORY: %d\n", pod.RequestGPUMemory)

			if pod.RequestGPUCore < 100 {
				gpuIdx := pod.PredicateGPUIdx[0]
				requestMemory := pod.RequestGPUMemory * 256
				fmt.Printf("GPU%d: CORE = %d, MEMORY (Mib) = %d\n", gpuIdx, pod.RequestGPUCore, requestMemory)
				usingGPUCore[gpuIdx] += pod.RequestGPUCore
				usingGPUMemory[gpuIdx] += requestMemory
			} else {
				for _, gpuIdx := range pod.PredicateGPUIdx {
					fmt.Printf("GPU%d: CORE = 100, MEMORY (Mib) = %d\n", gpuIdx, node.EachGpuMemory)
					usingGPUCore[gpuIdx] += 100
					usingGPUMemory[gpuIdx] += node.EachGpuMemory
				}
			}

			fmt.Println("--------------------")
		}

		fmt.Println()

		var totalUsingMemory int64 = 0
		for _, m := range usingGPUMemory {
			totalUsingMemory += m
		}
		totalMemory := node.EachGpuMemory * int64(node.GPUCount)

		var totalUsingGPUCore int64 = 0
		for _, c := range usingGPUCore {
			totalUsingGPUCore += c
		}

		for i := 0; i < node.GPUCount; i++ {
			fmt.Printf("GPU%d: CORE = %d, MEMORY (Mib) = %d / %d (%.0f%%)\n", i, usingGPUCore[i], usingGPUMemory[i], node.EachGpuMemory, float64(usingGPUMemory[i])/float64(node.EachGpuMemory)*100)
		}

		fmt.Printf("Allocated GPU Core In Node %s: %d / %d\n", node.Name, totalUsingGPUCore, 100*node.GPUCount)
		fmt.Printf("Allocated GPU Memory In Node %s: %d / %d (%.0f%%)\n", node.Name, totalUsingMemory, totalMemory, float64(totalUsingMemory)/float64(totalMemory)*100)

		fmt.Println("----------------------------------------")
	}

}
