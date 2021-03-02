package main

const (
	VCore                   = "tencent.com/vcuda-core"
	VMemory                 = "tencent.com/vcuda-memory"
	PredicateTimeAnnotation = "tencent.com/predicate-time"
	PredicateGPUIndexPrefix = "tencent.com/predicate-gpu-idx-"
	PredicateNode           = "tencent.com/predicate-node"
	GPUAssigned             = "tencent.com/gpu-assigned"
	HundredCore             = 100

	LabelVCudaKey = "node-role.kubernetes.io/vcuda"

	LabelGPUManagerKey   = "nvidia-device-enable"
	LabelGPUManagerValue = "enable"

	LabelGPUCount  = "nvidia.com/gpu.count"
	LabelGPUMemory = "nvidia.com/gpu.memory"

	StatusAddressInternalIP = "InternalIP"
)
