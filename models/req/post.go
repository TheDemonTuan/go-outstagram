package req

type PostMeEdit struct {
	Caption string `json:"caption" validate:"required,min=1,max=2200"`
}

type PostMeComment struct {
	Content string `json:"content" validate:"required,min=1,max=255"`
}
