package helpers

type LdapConfig struct {
	Address      string `json:"address"`
	Base         string `json:"base"`
	BindDN       string `json:"bind_dn"`
	BindPassword string `json:"bind_password"`
	UserFilter   string `json:"user_filter"`
	GroupFilter  string `json:"group_filter,omitempty"`
	SSL          struct {
		UseSSL     bool   `json:"use_ssl"`
		SkipTLS    bool   `json:"skip_tls"`
		SkipVerify bool   `json:"skip_verify"`
		ServerName string `json:"server_name"`
	} `json:"ssl"`
}

func (lc *LdapConfig) Init() (LDAPClient, error) {
	var ldapError error
	client := LDAPClient{
		Address:            lc.Address,
		Base:               lc.Base,
		BindDN:             lc.BindDN,
		BindPassword:       lc.BindPassword,
		GroupFilter:        lc.GroupFilter,
		UserFilter:         lc.UserFilter,
		ServerName:         lc.SSL.ServerName,
		UseSSL:             lc.SSL.UseSSL,
		SkipTLS:            lc.SSL.SkipTLS,
		InsecureSkipVerify: lc.SSL.SkipVerify,
	}
	ldapError = client.Connect()

	return client, ldapError
}
