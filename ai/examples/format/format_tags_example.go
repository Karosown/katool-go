package format

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/karosown/katool-go/ai"
	"github.com/karosown/katool-go/ai/types"
)

// User 用户信息（展示各种 tag 的使用）
//type User struct {
//	Name   string `json:"name" description:"用户姓名" example:"张三" minLength:"2" maxLength:"50"`
//	Age    int    `json:"age" description:"用户年龄" minimum:"0" maximum:"150"`
//	Email  string `json:"email" description:"用户邮箱地址" format:"email" example:"user@example.com"`
//	Gender string `json:"gender,omitempty" description:"用户性别" enum:"male,female,other"`
//	Phone  string `json:"phone,omitempty" description:"手机号码" pattern:"^1[3-9]\\d{9}$" example:"13800138000"`
//}

// Product 产品信息
type Product struct {
	ID          string   `json:"id" title:"产品ID" description:"产品唯一标识符" example:"PROD-001"`
	Name        string   `json:"name" title:"产品名称" description:"产品的显示名称" minLength:"1" maxLength:"100"`
	Price       float64  `json:"price" title:"价格" description:"产品价格（元）" minimum:"0"`
	Stock       int      `json:"stock" title:"库存" description:"当前库存数量" minimum:"0"`
	Tags        []string `json:"tags,omitempty" title:"标签" description:"产品标签列表" minItems:"0" maxItems:"10"`
	Status      string   `json:"status" title:"状态" description:"产品状态" enum:"active,inactive,discontinued"`
	Description string   `json:"description,omitempty" title:"描述" description:"产品详细描述" maxLength:"1000"`
}

// Order 订单信息（嵌套结构体）
type Order struct {
	OrderID   string    `json:"order_id" description:"订单号" example:"ORD-20240101-001"`
	UserInfo  User      `json:"user_info" description:"下单用户信息"`
	Products  []Product `json:"products" description:"订单商品列表" minItems:"1"`
	TotalCost float64   `json:"total_cost" description:"订单总金额" minimum:"0"`
	Status    string    `json:"status" description:"订单状态" enum:"pending,paid,shipped,completed,cancelled"`
}

func main() {
	fmt.Println("=== 结构体 Tag 增强示例 ===\n")

	// 创建客户端
	client, err := ai.NewClient()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	fmt.Println("示例1: 基本类型约束")
	fmt.Println("=" + "==================================================\n")
	example1(client)

	fmt.Println("\n示例2: 枚举和格式约束")
	fmt.Println("=" + "==================================================\n")
	example2(client)

	fmt.Println("\n示例3: 嵌套结构体")
	fmt.Println("=" + "==================================================\n")
	example3(client)
}

// 示例1: 基本类型约束
func example1(client *ai.Client) {
	// 生成 schema
	schema, err := ai.FormatFromType[User]()
	if err != nil {
		log.Fatal(err)
	}

	// 打印生成的 schema（查看所有 tag 信息）
	fmt.Println("生成的 Schema（包含所有 tag 信息）:")
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(schemaJSON))

	// 使用 schema 发送请求
	req := &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: "user", Content: "生成一个中国用户的信息"},
		},
		Format: schema,
	}

	fmt.Println("\n发送请求...")
	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	// 提取数据
	var user User
	if err := ai.UnmarshalStructuredData(response, &user, "extract_structured_data"); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("\n生成的用户信息:\n")
	fmt.Printf("  姓名: %s\n", user.Name)
	fmt.Printf("  年龄: %d\n", user.Age)
	fmt.Printf("  邮箱: %s\n", user.Email)
	fmt.Printf("  性别: %s\n", user.Gender)
	fmt.Printf("  电话: %s\n", user.Phone)
}

// 示例2: 枚举和格式约束
func example2(client *ai.Client) {
	schema, err := ai.FormatFromType[Product]()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("产品 Schema（注意 enum 和 format 约束）:")
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(schemaJSON))

	req := &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: "user", Content: "生成一个电子产品的信息"},
		},
		Format: schema,
	}

	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	var product Product
	if err := ai.UnmarshalStructuredData(response, &product, "extract_structured_data"); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("\n生成的产品信息:\n")
	fmt.Printf("  ID: %s\n", product.ID)
	fmt.Printf("  名称: %s\n", product.Name)
	fmt.Printf("  价格: %.2f 元\n", product.Price)
	fmt.Printf("  库存: %d\n", product.Stock)
	fmt.Printf("  状态: %s\n", product.Status)
	fmt.Printf("  标签: %v\n", product.Tags)
}

// 示例3: 嵌套结构体
func example3(client *ai.Client) {
	schema, err := ai.FormatFromType[Order]()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("订单 Schema（包含嵌套的 User 和 Product）:")
	schemaJSON, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(schemaJSON))

	req := &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: "user", Content: "生成一个包含2个商品的订单"},
		},
		Format: schema,
	}

	response, err := client.Chat(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	var order Order
	if err := ai.UnmarshalStructuredData(response, &order, "extract_structured_data"); err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("\n生成的订单信息:\n")
	fmt.Printf("  订单号: %s\n", order.OrderID)
	fmt.Printf("  用户: %s (%s)\n", order.UserInfo.Name, order.UserInfo.Email)
	fmt.Printf("  商品数量: %d\n", len(order.Products))
	for i, p := range order.Products {
		fmt.Printf("    %d. %s - %.2f 元\n", i+1, p.Name, p.Price)
	}
	fmt.Printf("  总金额: %.2f 元\n", order.TotalCost)
	fmt.Printf("  状态: %s\n", order.Status)
}
