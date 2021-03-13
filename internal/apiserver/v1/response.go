package v1

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ImageUploadResponse struct {
	ImageURI string `json:"image_uri"`
}
