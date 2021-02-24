package user

type User struct {
	ID           int     `gorm:"type:int(11);primary_key"`
	Email        string  `gorm:"type:varchar(50);not null;unique:email_UNIQUE"`
	Nickname     string  `gorm:"type:varchar(20);not null"`
	Gender       string  `gorm:"type:char(1);not null"`
	ProfileImage *string `gorm:"column:profile_img;type:char(100)"`
	Age          int     `gorm:"type:int(11);not null"`
	OAuthVendor  string  `gorm:"column:oauth_vendor;type:varchar(20);not null"`
	OAuthID      int     `gorm:"column:oauth_id;type:varchar(100);not null"`
}

func (User) TableName() string {
	return "user"
}
