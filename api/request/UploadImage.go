package request

type UploadImageRequest struct {
	ImageFile string `json:"imageFile" binding:"required"`
	SportId   int    `json:"sportId" binding:"required"`
}
