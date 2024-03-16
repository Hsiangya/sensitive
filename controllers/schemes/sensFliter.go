package schemes

type SensitiveStringCreate struct {
	Text          []string `json:"text" binding:"required,min=1"`
	SensitiveType string   `json:"sensitive_type"`
}

type SensitiveStringQuery struct {
	Text string `form:"text" binding:"required,min=1"`
}
