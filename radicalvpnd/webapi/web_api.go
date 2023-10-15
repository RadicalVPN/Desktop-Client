package webapi

type Signin struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type Server struct {
	Id        string `json:"id"`
	Hostname  string `json:"hostname"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Internalp string `json:"internal_ip"`
	ExternaIp string `json:"external_ip"`
	PublicKey string `json:"public_key"`
	Online    bool   `json:"online"`
	Latency   int    `json:"latency"`
}
