package pager

// Pager 分页器结构体
// Pager represents a pagination structure
type Pager struct {
	Page     int `json:"page" form:"page"`           // 当前页码 / Current page number
	PageSize int `json:"page_size" form:"page_size"` // 每页大小 / Page size
	Total    int `json:"total"`                      // 总记录数 / Total record count
}

// NewPager 创建新的分页器
// NewPager creates a new pager
func NewPager(page, pageSize int) *Pager {
	return &Pager{
		Page:     page,
		PageSize: pageSize,
		Total:    0,
	}
}

// Offset 计算偏移量
// Offset calculates the offset
func (p *Pager) Offset() int {
	if p.Page <= 0 || p.Page >= 1000 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Limit()
}

// Limit 获取限制数量
// Limit gets the limit count
func (p *Pager) Limit() int {
	if p.PageSize <= 0 || p.PageSize >= 1000 {
		p.PageSize = 20
	}
	return p.PageSize
}
