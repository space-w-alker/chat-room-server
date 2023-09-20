package generic

type PaginationArgs struct{
	Page uint `json:"page" form:"page" binding:"omitempty,number,gt=0"`
	Limit uint `json:"limit" form:"limit" binding:"omitempty,number,gt=0"`
}