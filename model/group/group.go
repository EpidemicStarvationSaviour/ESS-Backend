package group

type Group struct {
	GroupId     int    `gorm:"primaryKey"`
	GroupName   string `gorm:"uniqueIndex;size:100"`
	CreatorId   int
	Description string `gorm:"size:200"`
	Remark      string `gorm:"size:200"`
	Aid         int
	Cids        int
	RiderId     int
	Status      string
	Seen        bool
}

type GroupInfoReq struct {
	Type     int `json:"type"`
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type GroupInfoResp struct {
	Count int `json:"count"`
}
