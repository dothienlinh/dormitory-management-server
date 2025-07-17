package entity

type BaseModel interface {
	TableName() string
}

type Pagination struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
	Total int `json:"total"`
}

func (p *Pagination) Parse() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 10
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
