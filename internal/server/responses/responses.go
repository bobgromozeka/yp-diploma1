package responses

type Register struct {
	Token string `json:"token"`
}

type Login Register
