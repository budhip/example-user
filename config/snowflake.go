package config

import (
	"github.com/budhip/example-user/common"

	sx "github.com/bwmarrin/snowflake"
	logger "go.uber.org/zap"
)

type Snowflake interface {
	GetNewID() int64
}

type snowflake struct {
	node *sx.Node
}

func NewSnowflake() Snowflake {
	nodeID := common.GetNodeIDFromMachineIP()
	node, err := sx.NewNode(nodeID)

	if err != nil {
		logger.Error(err)
	}

	return &snowflake{
		node: node,
	}
}

func (i *snowflake) GetNewID() int64 {
	return i.node.Generate().Int64()
}