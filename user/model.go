package user

type User struct {
	ID           int     `gorm:"type:int(11);primary_key"`
	Email        string  `gorm:"type:char(11);not null;unique:email_UNIQUE"`
	Nickname     string  `gorm:"type:char(20);not null"`
	Gender       string  `gorm:"type:char(1);not null"`
	ProfileImage *string `gorm:"column:profile_img;type:char(100)"`
	Age          int     `gorm:"type:int(11)"`
	OAuthKakao   *string `gorm:"column:oauth.kakao;type:char(100)"`
	OAuthNaver   *string `gorm:"column:oauth.naver;type:char(100)"`
	OAuthGoogle  *string `gorm:"column:oauth.google;type:char(100)"`
}

func (User) TableName() string {
	return "user"
}
