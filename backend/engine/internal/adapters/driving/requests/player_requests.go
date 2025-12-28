package requests

type CreatePlayerRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type DeletePlayerRequest struct {
	Id string `path:"playerId" validate:"required"`
}

type UpdatePlayerRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}
