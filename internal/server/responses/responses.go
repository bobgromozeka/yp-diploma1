package responses

type Register struct {
	Token string `json:"token"`
}

type Login Register

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
