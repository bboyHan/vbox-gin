package request

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page      int    `json:"page" form:"page"`           // 页码
	PageSize  int    `json:"pageSize" form:"pageSize"`   // 每页大小
	StartTime int64  `json:"startTime" form:"startTime"` // 开始时间
	EndTime   int64  `json:"endTime" form:"endTime"`     // 结束时间
	Keyword   string `json:"keyword" form:"keyword"`     // 关键字 ("cnt"-数量,"sum"-总金额,"avg"-平均金额等等)
}

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}
