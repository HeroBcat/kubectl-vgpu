package main

import (
	"context"
	"strconv"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeInfo struct {
	Name          string
	InternalIP    string
	GPUCount      int
	EachGpuMemory int64
}

func getNodes() []NodeInfo {

	result := make([]NodeInfo, 0)

	nodeList, err := clientSet.CoreV1().Nodes().List(context.Background(), meta.ListOptions{})
	if err != nil {
		return result
	}

	for _, node := range nodeList.Items {

		if isNodeLabelExist(node, LabelGPUManagerKey, LabelGPUManagerValue) || isNodeLabelExist(node, LabelVCudaKey, "") {
			info := NodeInfo{
				Name:          node.Name,
				InternalIP:    getInternalIPFromNode(node),
				GPUCount:      getGPUCountFromNode(node),
				EachGpuMemory: getGPUMemoryFromNode(node),
			}
			result = append(result, info)
		}
	}

	return result
}

func isNodeLabelExist(node core.Node, key, value string) bool {

	labels := node.GetLabels()

	if v, ok := labels[key]; ok {
		if v == value {
			return true
		}
	}

	return false
}

func getInternalIPFromNode(node core.Node) string {
	for _, address := range node.Status.Addresses {
		if address.Type == StatusAddressInternalIP {
			return address.Address
		}
	}
	return ""
}

func getGPUCountFromNode(node core.Node) int {
	labels := node.GetLabels()
	if value, ok := labels[LabelGPUCount]; ok {
		v, err := strconv.Atoi(value)
		if err != nil {
			return 0
		}
		return v
	}
	return 0
}

func getGPUMemoryFromNode(node core.Node) int64 {
	labels := node.GetLabels()
	if value, ok := labels[LabelGPUMemory]; ok {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0
		}
		return v
	}
	return 0
}
