package types

type Meta struct {
	Page      int `json:"page"`
	Paginate  int `json:"paginate"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func (r *Meta) CountTotalPage(page, paginate, totalData int) {
	r.Page = page
	r.Paginate = paginate
	r.TotalData = totalData

	if totalData == 0 {
		r.TotalPage = 0
		return
	}

	r.TotalPage = totalData / r.Paginate
	if totalData%r.Paginate > 0 {
		r.TotalPage++
	}
}
