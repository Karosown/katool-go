package mongo_util

import (
	"fmt"
	"sort"

	"github.com/karosown/katool-go/container/optional"
	"github.com/twmb/murmur3"
)

//const (
//	ConversationName        = "conversation"
//	ConversationMessageName = "conversation_message"
//)
//
//// TableNumber 真实表的数量增加 必须是当前的倍数 这样迁移数据代价才会最小
//// 原因：现在环的构建方式是 012012012012
//const (
//	TableNumber = 20
//)

type DocumentPartitionHelper struct {
	// 总虚拟节点数
	replicas int
	// hash函数
	hash func(data []byte) uint32
	// 哈希环，按照节点哈希值排序
	ring []int
	// 节点哈希值到真实节点字符串，哈希映射的逆过程
	nodes map[int]string
	// ring的MaxSize
	ringSize int
	// tableBaseName
	tableBaseName string
}

func NewDefPartitionHelper(tableName string, sizes ...int) *DocumentPartitionHelper {
	replicas := optional.IsTrue(len(sizes) == 1 && sizes[0] != 0, sizes[0], 100000)
	ringSize := optional.IsTrue(len(sizes) == 2 && sizes[1] != 0, sizes[1], 20)
	return NewDocumentPartitionHelper(replicas, ringSize, tableName)
}
func NewDocumentPartitionHelper(replicas, ringSize int, tableBaseName string) *DocumentPartitionHelper {
	cli := &DocumentPartitionHelper{
		replicas:      replicas, // 100000
		hash:          murmur3.Sum32,
		nodes:         make(map[int]string),
		tableBaseName: tableBaseName,
		ringSize:      ringSize,
	}
	// 构建虚拟节点
	// 每个节点创建多个虚拟节点
	for i := 0; i < cli.replicas; i++ {
		// 每个虚拟节点计算哈希值
		hash := int(cli.hash([]byte(fmt.Sprintf("node%05d", i))))
		// 加入哈希环
		cli.ring = append(cli.ring, hash)
	}
	// 哈希环排序
	sort.Ints(cli.ring)
	for i, v := range cli.ring {
		key := i % cli.ringSize
		cli.nodes[v] = fmt.Sprintf("%s_%d", cli.tableBaseName, key)
	}
	return cli
}

// GetMessageName 获取Key对应的节点
func (r *DocumentPartitionHelper) GetCollName(key string) string {
	// 计算Key哈希值
	hash := int(r.hash([]byte(key)))

	// 二分查找第一个大于等于Key哈希值的节点
	idx := sort.Search(len(r.ring), func(i int) bool { return r.ring[i] >= hash })

	// 这里是特殊情况，也就是数组没有大于等于Key哈希值的节点
	// 但是逻辑上这是一个环，因此第一个节点就是目标节点
	if idx == len(r.ring) {
		idx = 0
	}

	// 返回哈希值对应的真实节点字符串
	return r.nodes[r.ring[idx]]
}
