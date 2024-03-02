package dto

type RespImage struct {
	ID string `json:"imageID"`
}

func RespImageFromID(id string) *RespImage {
	return &RespImage{
		ID: id,
	}
}
