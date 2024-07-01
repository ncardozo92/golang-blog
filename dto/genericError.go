package dto

type GenericErrorDTO struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}
