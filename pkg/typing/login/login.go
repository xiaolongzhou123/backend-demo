package login

type Login struct {
	User string `json:"user" binding:"required,usercheck" `
	Pass string `json:"pass" binding:"required,passcheck" `
}
