package snowflake

import (
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

type Config struct {
	StartTime string `toml:"start_time"`
	MachineID int64  `toml:"machine_id"`
}

func init() {
	err := NewNode(&Config{
		StartTime: "2023-01-01",
		MachineID: 1,
	})
	if err != nil {
		panic(fmt.Errorf("snowflake: failed to initialize node: %v", err))
	}
}

func NewNode(c *Config) error {
	if c == nil {
		return fmt.Errorf("snowflake: nil config")
	}
	if c.StartTime == "" {
		return fmt.Errorf("snowflake: empty start_time")
	}
	if c.MachineID < 0 || c.MachineID > 1023 {
		return fmt.Errorf("snowflake: invalid machine_id: %d", c.MachineID)
	}
	st, err := time.Parse("2006-01-02", c.StartTime)
	if err != nil {
		return err
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(c.MachineID)
	if err != nil {
		return err
	}
	return nil
}

func GenID() int64 {
	return node.Generate().Int64()
}
