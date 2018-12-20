package istiocertmerger

type Port struct {
	Name     string `json:"name"`
	Number   int32  `json:"number"`
	Protocol string `json:"protocol"`
}

type TLS struct {
	Mode              string `json:"mode,omitempty"`
	PrivateKey        string `json:"privateKey,omitempty"`
	ServerCertificate string `json:"serverCertificate,omitempty"`
	HttpsRedirect     bool   `json:"httpsRedirect"`
}

type GatewayServer struct {
	Hosts []string `json:"hosts"`
	Port  *Port    `json:"port"`
	TLS   *TLS     `json:"tls"`
}
