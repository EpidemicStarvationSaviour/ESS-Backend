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
