package model

type Nationality struct {
	NationalityID   int    `json:"nationality_id"`
	NationalityName string `json:"nationality_name"`
	NationalityCode string `json:"nationality_code"`
}
