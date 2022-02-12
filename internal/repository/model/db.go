package model

type Session struct {
	ID        int    `gorm:"primary_key;auto_increment;unique"`
	SessionID string `gorm:"type:varchar(255);not null;unique"`
}

type Shorten struct {
	ID        int    `gorm:"primary_key;auto_increment;unique"`
	URLID     string `gorm:"type:varchar(255);not null;unique"`
	ShortURL  string `gorm:"type:varchar(255);not null"`
	LongURL   string `gorm:"type:varchar(255);not null;unique"`
	SessionID int    `gorm:"type:int;not null"`
	IsDeleted bool   `gorm:"type:bool; default:false"`
}
