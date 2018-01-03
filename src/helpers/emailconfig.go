package helpers

type EmailConfig struct {
	Server   string `json:"server"`
	From     string `json:"from"`
	Username string `json:"username"`
	Password string `json:"password"`
	UseTLS   bool   `json:"use_tls"`
}
