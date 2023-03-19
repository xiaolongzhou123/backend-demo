package login

type Login struct {
	User string `json:"user" binding:"required,usercheck" `
	Pass string `json:"pass" binding:"required,passcheck" `
}

type ModifyLogin struct {
	User    string `json:"user" binding:"required,usercheck" `
	OldPass string `json:"old_pass" binding:"required,passcheck" `
	CurPass string `json:"cur_pass" binding:"required,passcheck" `
}
