package ldap

import (
	"fmt"
	"testing"

	"github.com/go-ldap/ldap"
)

func TestNewLDAP(t *testing.T) {

	url := "18.18.2.2:389"
	l, err := ldap.Dial("tcp", url)
	fmt.Println("dial===", l, err)
	err = l.Bind("cn=admin,dc=example,dc=org", "admin")
	if err != nil {
		fmt.Println("ldap Bind err:=", err)
	}

}
