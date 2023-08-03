package requests

type Register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Login Register
