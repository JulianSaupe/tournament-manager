package player

type CreatePlayerRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type DeletePlayerRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdatePlayerRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}
