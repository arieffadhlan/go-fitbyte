package dto

type UserResponse struct {
	Email    string  `json:"email"`
	ImageURI *string `json:"image_uri"`
}
