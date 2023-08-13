package responses

type Register struct {
	Token string `json:"token"`
}

type Login Register

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
