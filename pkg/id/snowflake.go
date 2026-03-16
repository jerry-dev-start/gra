package id

import (
	"fmt"
	"sync"
	"time"
)

// 雪花算法位分配（共63位，首位符号位固定0）
// 41位时间戳 | 10位机器ID | 12位序列号
const (
	epoch          int64 = 1700000000000 // 起始纪元：2023-11-14，可用约69年
	machineBits    uint8 = 10            // 机器ID位数，最大1023
	sequenceBits   uint8 = 12            // 序列号位数，每毫秒4096个
	machineMax     int64 = -1 ^ (-1 << machineBits)
	sequenceMax    int64 = -1 ^ (-1 << sequenceBits)
	machineShift         = sequenceBits
	timestampShift       = machineBits + sequenceBits
)

// Node 雪花算法节点
type Node struct {
	mu        sync.Mutex
	timestamp int64
	machineID int64
	sequence  int64
}

var defaultNode *Node

// Init 初始化全局雪花节点
func Init(machineID int64) error {
	if machineID < 0 || machineID > machineMax {
		return fmt.Errorf("machine ID must be between 0 and %d", machineMax)
	}
	defaultNode = &Node{machineID: machineID}
	return nil
}

// Generate 生成雪花ID（使用全局节点）
func Generate() int64 {
	return defaultNode.Generate()
}

// Generate 生成雪花ID
func (n *Node) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli() - epoch

	if now == n.timestamp {
		n.sequence = (n.sequence + 1) & sequenceMax
		if n.sequence == 0 {
			// 当前毫秒序列号用尽，等待下一毫秒
			for now <= n.timestamp {
				now = time.Now().UnixMilli() - epoch
			}
		}
	} else {
		n.sequence = 0
	}

	n.timestamp = now

	return now<<timestampShift | n.machineID<<int64(machineShift) | n.sequence
}
