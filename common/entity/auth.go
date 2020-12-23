package entity

type UserRegister struct {
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email,omitempty"`
	Uuid     string `json:"-"`
	Status   string `json:"-"`
}

type UserLogin struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool	`json:"remember"`
}

type Follow struct {
	UserID int `json:"user_id"`
}
