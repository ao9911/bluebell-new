package snowflake

import (
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func init() {
	if err := InitSF("2020-01-01", 1); err != nil {
		panic(fmt.Sprintf("init snowflake: %v", err))
	}
}

func InitSF(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}
