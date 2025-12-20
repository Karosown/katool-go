package main

import (
	"context"
	"fmt"
	"github.com/karosown/katool-go/mq"
	"github.com/karosown/katool-go/mq/cmq"
	"math/rand"
	"time"
)

// OrderDetail 业务对象，将放入 Extra 中
type OrderDetail struct {
	OrderID   string
	Amount    float64
	ItemCount int
}

func main() {
	// -------------------------------------------------------
	// 1. 初始化 (切换注释即可换底层实现)
	// -------------------------------------------------------

	// Mode A: 内存版
	broker := cmq.NewChanBroker()

	// Mode B: Redis 版 (需本地运行 Redis)
	// broker := redis_mq.NewRedisClient("localhost:6379", "", 0)

	defer broker.Close()
	ctx := context.Background()
	topic := "new_orders"

	fmt.Println(">>> 订单系统启动...")

	// -------------------------------------------------------
	// 2. 消费者：审计员 (只关心金额 > 1000 的订单)
	// -------------------------------------------------------
	broker.Subscribe(ctx, topic, func(ctx context.Context, msg mq.Message) error {
		md := msg.GetMetadata()

		// 从 Extra 恢复结构体
		if detail, ok := md.Extra["detail"].(OrderDetail); ok {
			fmt.Printf("💰 [Audit-VIP] 发现大额订单! ID: %s | Amount: %.2f\n", detail.OrderID, detail.Amount)
		} else {
			// Redis JSON 反序列化回来是 map[string]interface{}，需要二次处理(生产环境通常用 JSON 库转)
			// 这里为了演示简单，直接打印 map
			fmt.Printf("💰 [Audit-VIP] 发现大额订单! Raw: %v\n", md.Extra["detail"])
		}
		return nil
	},
		mq.WithGroup("audit_group"),
		// [关键] 客户端过滤器
		mq.WithFilter(func(md mq.Metadata) bool {
			// 检查 Extra 中的金额 (兼容 Redis 的 float64/map 解析差异)
			// 内存版直接断言 OrderDetail
			if d, ok := md.Extra["detail"].(OrderDetail); ok {
				return d.Amount > 1000
			}
			// Redis 版简单检查 (实际需 json unmarshal)
			return true
		}), mq.WithFilter(func(md mq.Metadata) bool {
			return md.Key == "User_VIP"
		}))

	// -------------------------------------------------------
	// 3. 消费者：发货系统 (处理所有分区，模拟耗时)
	// -------------------------------------------------------
	broker.Subscribe(ctx, topic, func(ctx context.Context, msg mq.Message) error {
		md := msg.GetMetadata()
		fmt.Printf("🚚 [Shipping]  处理订单 OrderID:%-10s | Key: %-10s | Partition: %d\n", md.Extra["detail"].(OrderDetail).OrderID, md.Key, md.Partition)
		time.Sleep(50 * time.Millisecond) // 模拟处理耗时
		return nil
	}, mq.WithGroup("shipping_group"))

	// -------------------------------------------------------
	// 4. 生产者：模拟产生流量
	// -------------------------------------------------------
	time.Sleep(500 * time.Millisecond) // 等待订阅就绪

	users := []string{"User_A", "User_B", "User_C", "User_VIP"}

	for i := 0; i < 10; i++ {
		u := users[rand.Intn(len(users))]
		amount := rand.Float64() * 2000 // 0~2000 随机金额

		order := OrderDetail{
			OrderID:   fmt.Sprintf("ORD-%d", i),
			Amount:    amount,
			ItemCount: rand.Intn(5) + 1,
		}

		// 发送
		broker.Publish(ctx, topic, []byte("OrderPayload"),
			mq.WithKey(u), // 保证同一个用户的订单去同一个分区
			mq.WithExtra("detail", order),
			mq.WithExtra("timestamp", time.Now().Unix()),
		)

		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println(">>> 演示结束，等待处理完成...")
	time.Sleep(2 * time.Second)
}
