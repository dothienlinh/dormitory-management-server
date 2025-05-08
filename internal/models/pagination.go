package models

type Pagination struct {
	Page  int   `form:"page"`
	Limit int   `form:"limit"`
	Total int64 `form:"-"`
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) GetLimit() int {
	return p.Limit
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) Parse() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}
