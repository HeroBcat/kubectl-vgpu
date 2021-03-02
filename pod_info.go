package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type PodInfo struct {
	Name             string
	NameSpace        string
	RequestGPUCore   int64
	RequestGPUMemory int64
	PredicateGPUIdx  []int
	UsingGPUIdx      []bool
}

var (
	inCluster = false
	clientSet *k8s.Clientset
)

func kubeInit() {

	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{}).ClientConfig()
	}

	if err != nil {
		log.Fatalf("could not get kubernetes config: %s", err)
	}

	clientSet, err = k8s.NewForConfig(config)
	if err != nil {
		log.Fatalf("could not get kubernetes config: %s", err)
	}

}

func getActivePodsByNode(node NodeInfo) ([]PodInfo, error) {

	result := make([]PodInfo, 0)

	selector := fields.SelectorFromSet(fields.Set{"spec.nodeName": node.Name})

	podList, err := clientSet.CoreV1().Pods(meta.NamespaceAll).List(context.Background(), meta.ListOptions{
		FieldSelector: selector.String(),
		LabelSelector: labels.Everything().String(),
	})

	if err != nil {
		return result, err
	}

	for _, pod := range podList.Items {
		if hasPodPredicated(pod) && isActivePod(pod) {
			predicateGPUIdx := getPredicateGPUIdx(pod)
			info := PodInfo{
				Name:             pod.Name,
				NameSpace:        pod.Namespace,
				RequestGPUCore:   getRequestGPUCore(pod),
				RequestGPUMemory: getRequestGPUMemory(pod),
				PredicateGPUIdx:  predicateGPUIdx,
				UsingGPUIdx:      getUsingGPUIdx(node.GPUCount, predicateGPUIdx),
			}
			result = append(result, info)
		}
	}

	return result, nil
}

func hasPodPredicated(pod core.Pod) bool {
	for k := range pod.Annotations {
		if strings.Contains(k, GPUAssigned) ||
			strings.Contains(k, PredicateTimeAnnotation) ||
			strings.Contains(k, PredicateGPUIndexPrefix) {
			return true
		}
	}
	return false
}

func isActivePod(pod core.Pod) bool {
	if pod.Status.Phase == core.PodSucceeded || pod.Status.Phase == core.PodFailed {
		return false
	}
	return true
}

func getPredicateGPUIdx(pod core.Pod) []int {
	result := make([]int, 0)
	for k, v := range pod.Annotations {
		if strings.HasPrefix(k, PredicateGPUIndexPrefix) {
			splits := strings.Split(v, ",")
			for _, v := range splits {
				idx, err := strconv.Atoi(v)
				if err != nil {
					continue
				}
				result = append(result, idx)
			}
		}
	}
	return result
}

func getRequestGPUCore(pod core.Pod) int64 {
	var result int64 = 0
	for _, container := range pod.Spec.Containers {
		if v, ok := container.Resources.Requests[VCore]; ok {
			result += v.Value()
		}
	}
	return result
}

func getRequestGPUMemory(pod core.Pod) int64 {
	var result int64 = 0
	for _, container := range pod.Spec.Containers {
		if v, ok := container.Resources.Requests[VMemory]; ok {
			result += v.Value()
		}
	}
	return result
}

func getUsingGPUIdx(gpuCount int, gpuIdx []int) []bool {
	result := make([]bool, gpuCount)
	for _, idx := range gpuIdx {
		if idx >= gpuCount {
			continue
		}
		result[idx] = true
	}
	return result
}
