package clusters

import (
	"clustercloner/clusters/util"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

// ClusterInfo ...
type ClusterInfo struct {
	Cloud         string //GCP, Azure, AWS, or Hub (for a standard neutral format)
	Scope         string //Project in GKE, Subscription in AKS, blank in EKS
	Location      string //Zone or Region in GKE, Region in others
	Name          string
	K8sVersion    string
	GeneratedBy   string
	Labels        map[string]string
	NodePools     []NodePoolInfo
	SourceCluster *ClusterInfo
}

// AddNodePool ...
func (ci *ClusterInfo) AddNodePool(pool NodePoolInfo) {
	if ci.NodePools == nil {
		ci.NodePools = make([]NodePoolInfo, 0)
	}
	ci.NodePools = append(ci.NodePools, pool)
}

// MachineType ...
type MachineType struct {
	Name  string
	CPU   int
	RAMMB int
}

// NodePoolInfo ...
type NodePoolInfo struct {
	Name        string
	NodeCount   int
	K8sVersion  string
	MachineType MachineType
	DiskSizeGB  int
	Preemptible bool
}

var (
	// Mock Created for testing
	Mock = "Mock"
	// Read from the cloud
	Read = "Read"
	// Created using this tool
	Created = "Created"
	// Transformation  by this tool; can be an intermediate transformation step (Hub) or the output of the transformation, which can be optionally created in the cloud
	Transformation = "Transformation"

	// SearchTemplate used for searching for Clusters by label
	SearchTemplate = "SearchTemplate"

	// InputFile read from a JSON file. Even if the Cluster in the JSON file does not say InputFile, the value will be replaced
	InputFile = "InputFile"
)

var (
	// Hub ...
	Hub = "Hub"
	// GCP ...
	GCP = "GCP"
	// Azure ...
	Azure = "Azure"
	// AWS ...
	AWS = "AWS"
)

// LoadFromFile ...
func LoadFromFile(inputFile string) (ret []*ClusterInfo, err error) {
	if inputFile[0:1] == "/" {
		inputFile = inputFile[1:]
	}
	fn := util.RootPath() + "/" + inputFile
	jsonBytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load input file "+inputFile)
	}

	err = json.Unmarshal(jsonBytes, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshall input file "+inputFile)
	}
	for _, ci := range ret {
		labelsNilToEmptyMap(ci)
	}
	for _, ci := range ret {
		ci.GeneratedBy = InputFile
	}

	return ret, nil
}

func labelsNilToEmptyMap(ci *ClusterInfo) {
	if ci.Labels == nil {
		ci.Labels = make(map[string]string)
	}
	if ci.SourceCluster != nil {
		labelsNilToEmptyMap(ci.SourceCluster)
	}
}
