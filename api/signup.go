package api

type SignupRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Repassword string `json:"repassword" binding:"required,eqfield=Password"`
}

type SignupResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}
