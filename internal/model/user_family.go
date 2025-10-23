package model

type User struct {
	UserID        int      `json:"user_id"`
	Name          string   `json:"name"`
	Dob           string   `json:"dob"`
	NationalityID int      `json:"national_id"`
	Families      []Family `json:"families"`
}

type Family struct {
	FamilyID int    `json:"family_id"`
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Dob      string `json:"dob"`
}

type FamiliesJSON struct {
	Families []byte `json:"families"`
}

type UserDetailResponse struct {
	UserID        int         `json:"user_id"`
	Name          string      `json:"name"`
	Dob           string      `json:"dob"`
	NationalityID int         `json:"national_id"`
	Nationality   Nationality `json:"nationality"`
	Families      []Family    `json:"families"`
}
