package dto

type QueryParams struct {
	Query string
	Order string
	Sort  string
}

func (qp *QueryParams) isSortValid() bool {
	return qp.Sort == "asc" || qp.Sort == "desc"
}

type PaginationParam struct {
	Page   int
	Offset int
	Limit  int
}

type PaginationQuery struct {
	Page int
	Take int
	Skip int
}

type Paging struct {
	Page        int
	RowsPerPage int
	TotalRows   int
	TotalPages  int
}

type RequestQueryParam struct {
	QueryParams
	PaginationParam
}