package web_crawler

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// CaptchaSolver 验证码解决器
// CaptchaSolver captcha solver
type CaptchaSolver struct {
	Page *rod.Page
}

// NewCaptchaSolver 创建验证码解决器
// NewCaptchaSolver creates a captcha solver
func NewCaptchaSolver(page *rod.Page) *CaptchaSolver {
	return &CaptchaSolver{Page: page}
}

// HumanDrag 模拟人类拖拽行为
// HumanDrag simulates human drag behavior
func (cs *CaptchaSolver) HumanDrag(box *rod.Element, x, y float64) error {
	// 获取元素的初始位置
	// Get element initial position
	boxRect, err := box.Shape()
	if err != nil {
		return err
	}
	startX := boxRect.Box().X + boxRect.Box().Width/2
	startY := boxRect.Box().Y + boxRect.Box().Height/2

	// 移动鼠标到起点
	// MoveTo mouse to start
	if err := cs.Page.Mouse.MoveTo(proto.NewPoint(startX, startY)); err != nil {
		return err
	}
	if err := cs.Page.Mouse.Down(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	// 模拟人类移动轨迹
	// Simulate human movement trajectory
	// 目标位置
	targetX := startX + x
	targetY := startY + y

	// 使用贝塞尔曲线或物理模拟生成轨迹
	// Use Bézier curve or physics simulation to generate trajectory
	// 这里使用简单的物理模拟算法（模拟手部抖动和加速减速）

	currentX := startX
	currentY := startY

	// 总距离
	totalDist := math.Sqrt(math.Pow(currentY-currentX, 2) + math.Pow(targetY-startY, 2))
	// 步数，距离越长步数越多
	steps := int(totalDist / 2)
	if steps < 10 {
		steps = 10
	}

	// 随机生成控制点
	controlX1 := startX + (targetX-startX)/3 + (rand.Float64()-0.5)*50
	controlY1 := startY + (targetY-startY)/3 + (rand.Float64()-0.5)*50
	controlX2 := startX + 2*(targetX-startX)/3 + (rand.Float64()-0.5)*50
	controlY2 := startY + 2*(targetY-startY)/3 + (rand.Float64()-0.5)*50

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		// 三阶贝塞尔曲线
		// Cubic Bezier Curve
		// B(t) = (1-t)^3 P0 + 3(1-t)^2 t P1 + 3(1-t) t^2 P2 + t^3 P3
		cx := math.Pow(1-t, 3)*startX +
			3*math.Pow(1-t, 2)*t*controlX1 +
			3*(1-t)*math.Pow(t, 2)*controlX2 +
			math.Pow(t, 3)*targetX

		cy := math.Pow(1-t, 3)*startY +
			3*math.Pow(1-t, 2)*t*controlY1 +
			3*(1-t)*math.Pow(t, 2)*controlY2 +
			math.Pow(t, 3)*targetY

		// 添加随机抖动
		// Add random jitter
		jitter := (rand.Float64() - 0.5) * 2
		cx += jitter
		cy += jitter

		// 移动
		if err := cs.Page.Mouse.MoveTo(proto.NewPoint(cx, cy)); err != nil {
			return err
		}

		// 随机停顿，模拟思考
		// Random pause, simulating thinking
		if rand.Float64() < 0.05 {
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		}

		// 变速移动：开始慢，中间快，结束慢
		// Variable speed: slow start, fast middle, slow end
		// 通过调整Sleep时间来实现
		sleepTime := time.Duration(rand.Intn(5)+2) * time.Millisecond
		if t < 0.1 || t > 0.9 {
			sleepTime += time.Duration(rand.Intn(10)) * time.Millisecond
		}
		time.Sleep(sleepTime)
	}

	// 稍微过冲一点再回退（模拟真实操作）
	// Overshoot slightly and move back (simulate real operation)
	if rand.Float64() > 0.5 {
		overshootX := targetX + (rand.Float64() * 10)
		if err := cs.Page.Mouse.MoveTo(proto.NewPoint(overshootX, targetY)); err != nil {
			return err
		}
		time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)
		if err := cs.Page.Mouse.MoveTo(proto.NewPoint(targetX, targetY)); err != nil {
			return err
		}
	}

	// 释放鼠标
	// Release mouse
	if err := cs.Page.Mouse.Up(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	return nil
}

// FindSliderGap 简单的滑块缺口识别（基于像素差异）
// FindSliderGap simple slider gap detection (based on pixel difference)
// bgElem: 背景图片元素 / background image element
// sliderElem: 滑块元素（用于确定起始搜索位置，避免识别到滑块本身） / slider element (to determine start search position, avoid detecting slider itself)
func (cs *CaptchaSolver) FindSliderGap(bgElem *rod.Element, sliderElem *rod.Element) (int, error) {
	// 获取背景图截图
	bgBytes, err := bgElem.Screenshot(proto.PageCaptureScreenshotFormatPng, 0)
	if err != nil {
		return 0, err
	}

	img, _, err := image.Decode(bytes.NewReader(bgBytes))
	if err != nil {
		return 0, err
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 获取滑块位置，确定搜索起点
	sliderRect, err := sliderElem.Shape()
	if err != nil {
		return 0, err
	}
	// 从滑块右侧开始搜索
	startScanX := int(sliderRect.Box().Width) + 10
	if startScanX >= width {
		startScanX = width / 5 // 默认从1/5处开始
	}

	// 简单的边缘检测算法
	// Simple edge detection algorithm
	// 遍历像素，寻找颜色突变的垂直线

	// 为了简化，我们取中间几行进行扫描
	// For simplicity, we scan a few middle rows
	targetX := 0
	found := false

	// 阈值
	threshold := 5000 // 颜色差异阈值，需要根据实际情况调整

	for x := startScanX; x < width-10; x++ {
		diffSum := 0
		// 扫描垂直线上的多个点
		scanRows := []int{height / 3, height / 2, height * 2 / 3}

		for _, y := range scanRows {
			// 获取当前像素和右侧像素
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r2, g2, b2, _ := img.At(x+1, y).RGBA()

			// 计算差异 (RGBA返回的是0-65535)
			diff := abs(int(r1)-int(r2)) + abs(int(g1)-int(g2)) + abs(int(b1)-int(b2))
			diffSum += diff
		}

		// 平均差异
		avgDiff := diffSum / len(scanRows)

		if avgDiff > threshold {
			// 找到缺口左边缘
			targetX = x
			found = true
			break
		}
	}

	if !found {
		// 如果没找到，尝试返回中间偏右的位置，碰碰运气
		return width / 2, nil
	}

	return targetX, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// SolveSlideToVerify 解决“滑动验证”类型（拖到最右侧或指定缺口）
// SolveSlideToVerify solves "slide to verify" type (drag to rightmost or specific gap)
// sliderSelector: 滑块按钮选择器 / slider button selector
// bgSelector: 背景图选择器（可选，如果是缺口验证则需要） / background image selector (optional, needed if gap verification)
// containerSelector: 滑块容器选择器（用于计算拖拽距离） / slider container selector (used to calculate drag distance)
func (cs *CaptchaSolver) SolveSlideToVerify(sliderSelector, containerSelector string, bgSelector ...string) error {
	// 等待元素加载
	slider, err := cs.Page.Element(sliderSelector)
	if err != nil {
		return err
	}
	container, err := cs.Page.Element(containerSelector)
	if err != nil {
		return err
	}

	// 计算拖拽距离
	// Calculate drag distance
	sliderRect, err := slider.Shape()
	if err != nil {
		return err
	}
	containerRect, err := container.Shape()
	if err != nil {
		return err
	}

	distance := 0.0

	if len(bgSelector) > 0 && bgSelector[0] != "" {
		// 缺口验证模式
		// Gap verification mode
		bg, err := cs.Page.Element(bgSelector[0])
		if err == nil {
			gapX, err := cs.FindSliderGap(bg, slider)
			if err == nil && gapX > 0 {
				// 调整距离：目标位置 - 滑块当前相对背景的X偏移
				// 这里的计算可能需要根据实际页面调整，通常缺口X就是需要拖动的距离
				// 假设背景图和滑块容器对齐
				distance = float64(gapX) - (sliderRect.Box().X - containerRect.Box().X)

				// 修正：有时需要减去滑块自身的一半宽度或微调
				distance -= 5
			}
		}
	}

	if distance <= 0 {
		// 默认模式：拖到最右侧
		// Default mode: drag to rightmost
		distance = containerRect.Box().Width - sliderRect.Box().Width
	}

	// 执行拖拽
	// Execute drag
	return cs.HumanDrag(slider, distance, 0)
}

// SolveSimpleSlider 通用简单滑块解决入口
// SolveSimpleSlider generic simple slider solver entry point
func (c *Client) SolveSimpleSlider(sliderSelector, containerSelector string) error {
	chrome := c.getChrome()
	if chrome == nil {
		return nil
	}
	// 获取当前页面
	// Get current page
	// 注意：这里假设Client操作的是最近活跃的页面，或者需要从Chrome实例获取页面
	// 由于Client结构体目前没有持有当前Page引用，这是一个设计上的限制。
	// 我们尝试获取最新打开的页面。
	// Note: assuming Client operates on the most recently active page, or need to get page from Chrome instance.
	// Since Client struct currently doesn't hold current Page reference, this is a design limitation.
	// We try to get the latest opened page.

	pages, err := chrome.Browser.Pages()
	if err != nil || len(pages) == 0 {
		return err
	}
	page := pages[len(pages)-1] // 取最后一个页面

	solver := NewCaptchaSolver(page)
	return solver.SolveSlideToVerify(sliderSelector, containerSelector)
}

// SolvePuzzleSlider 通用缺口滑块解决入口
// SolvePuzzleSlider generic puzzle slider solver entry point
func (c *Client) SolvePuzzleSlider(sliderSelector, containerSelector, bgSelector string) error {
	chrome := c.getChrome()
	if chrome == nil {
		return nil
	}
	pages, err := chrome.Browser.Pages()
	if err != nil || len(pages) == 0 {
		return err
	}
	page := pages[len(pages)-1]

	solver := NewCaptchaSolver(page)
	return solver.SolveSlideToVerify(sliderSelector, containerSelector, bgSelector)
}
