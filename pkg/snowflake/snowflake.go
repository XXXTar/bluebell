package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// 单例模式？
var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time //指定时间因子
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
