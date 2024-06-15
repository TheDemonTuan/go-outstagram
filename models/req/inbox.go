package req

type InboxSendMessage struct {
	Message string `json:"message" validate:"required,min=1,max=2000"`
}
