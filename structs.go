package election

import (
	"github.com/samuel/go-zookeeper/zk"
	"github.com/uol/funks"
)

//
// All structs and constants used by this library.
// author: rnojiri
//

// Master - signals for master role acquisition
const Master = 1

// Slave - signals for slave role acquisition
const Slave = 2

// ClusterChanged - signals for cluster change
const ClusterChanged = 3

// Disconnected - int signal for disconnection
const Disconnected = 4

// Config - configures the election
type Config struct {
	ZKURL                  []string       `json:"zkURL"`
	ZKElectionNodeURI      string         `json:"zkElectionNodeURI"`
	ZKSlaveNodesURI        string         `json:"zkSlaveNodesURI"`
	ReconnectionTimeout    funks.Duration `json:"reconnectionTimeout"`
	SessionTimeout         funks.Duration `json:"sessionTimeout"`
	ClusterChangeCheckTime funks.Duration `json:"clusterChangeCheckTime"`
	ClusterChangeWaitTime  funks.Duration `json:"clusterChangeWaitTime"`
}

// Cluster - has cluster info
type Cluster struct {
	IsMaster bool
	Master   string
	Slaves   []string
	Nodes    []string
	NumNodes int
}

const (
	// EventDisconnected - specifies a custom event for disconnection
	EventDisconnected zk.EventType = 99
)
