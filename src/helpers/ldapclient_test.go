package helpers

import "testing"

var (
	ldapClient = LDAPClient{
		Address:      "ldap.forumsys.com:389",
		Base:         "dc=example,dc=com",
		BindDN:       "cn=read-only-admin,dc=example,dc=com",
		BindPassword: "password",
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		UseSSL:       false,
		SkipTLS:      true,
	}
	testLogin = "newton"
	testPass  = "password"
)

func TestLDAPClient_Connect(t *testing.T) {
	defer ldapClient.Close()
	err := ldapClient.Connect()
	if err != nil {
		t.Error(err)
	}
}

func TestLDAPClient_Authenticate(t *testing.T) {
	defer ldapClient.Close()
	ok, _, err := ldapClient.Authenticate(testLogin, testPass)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("auth failed")
	}
}

func TestLDAPClient_GetGroupsOfUser(t *testing.T) {
	defer ldapClient.Close()
	_, err := ldapClient.GetGroupsOfUser(testLogin)
	if err != nil {
		t.Error(err)
	}
}
