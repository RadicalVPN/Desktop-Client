package webapi

type Signin struct {
	Email              string `json:"email"`
	Password           string `json:"password"`
	TotpToken          string `json:"totpToken"`
	RememberMe         bool   `json:"rememberMe"`
	TurnstileChallenge string `json:"turnstileChallenge"`
}

type VpnConnect struct {
	NodeLocation    string `json:"nodeLocation"`
	PrivacyFirewall string `json:"privacyFirewall"`
}

type Server struct {
	Id          string  `json:"id"`
	Hostname    string  `json:"hostname"`
	Country     string  `json:"country"`
	CountryName string  `json:"country_name"`
	City        string  `json:"city"`
	Latitude    string  `json:"latitude"`
	Longitude   string  `json:"longitude"`
	Internalp   string  `json:"internal_ip"`
	ExternaIp   string  `json:"external_ip"`
	PublicKey   string  `json:"public_key"`
	Online      bool    `json:"online"`
	Latency     int     `json:"latency"`
	Load        float64 `json:"load"`
	Location    string  `json:"location"`
}

type VpnTransfer struct {
	Rx int `json:"rx"`
	Tx int `json:"tx"`
}

type CurrentVpnTransfer struct {
	Rx int `json:"rx"`
	Tx int `json:"tx"`
}

type VpnStatus struct {
	AllowedIps          []string           `json:"allowedIps"`
	LatestHandshakeAt   string             `json:"latestHandshakeAt"`
	Transfer            VpnTransfer        `json:"transfer"`
	Current             CurrentVpnTransfer `json:"current"`
	PersistentKeepalive string             `json:"persistentKeepalive"`
}

type Vpn struct {
	Id        string    `json:"id"`
	Active    bool      `json:"active"`
	Alias     string    `json:"alias"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	Node      string    `json:"node"`
	Status    VpnStatus `json:"status"`
}

type DynamicVpnCreation struct {
	NodeLocation    string `json:"nodeLocation"`
	PrivacyFirewall string `json:"privacyFirewall"`
}
