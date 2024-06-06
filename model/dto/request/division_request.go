package request

type DivisionRequest struct {
	Id   string `json:"id"`
	Name string `json:"name" binding:"required"`
}
