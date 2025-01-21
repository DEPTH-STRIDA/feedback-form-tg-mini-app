package request

type FormRequest struct {
	Name     string `json:"name" validate:"required,max=128"`
	Feedback string `json:"feedback" validate:"max=256"`
	Comment  string `json:"comment" validate:"max=512"`
}
