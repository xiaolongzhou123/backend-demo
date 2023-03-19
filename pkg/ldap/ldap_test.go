package ldap

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-ldap/ldap"
)

func TestNewLDAP(t *testing.T) {

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "18.18.2.2", 389))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	controls := []ldap.Control{}
	controls = append(controls, ldap.NewControlBeheraPasswordPolicy())
	bindRequest := ldap.NewSimpleBindRequest("cn=admin,dc=example,dc=org", "admin", controls)

	r, err := l.SimpleBind(bindRequest)
	ppolicyControl := ldap.FindControl(r.Controls, ldap.ControlTypeBeheraPasswordPolicy)

	var ppolicy *ldap.ControlBeheraPasswordPolicy
	if ppolicyControl != nil {
		ppolicy = ppolicyControl.(*ldap.ControlBeheraPasswordPolicy)
	} else {
		log.Printf("ppolicyControl response not available.\n")
	}
	if err != nil {
		errStr := "ERROR: Cannot bind: " + err.Error()
		if ppolicy != nil && ppolicy.Error >= 0 {
			errStr += ":" + ppolicy.ErrorString
		}
		log.Print(errStr)
	} else {
		logStr := "Login Ok"
		if ppolicy != nil {
			if ppolicy.Expire >= 0 {
				logStr += fmt.Sprintf(". Password expires in %d seconds\n", ppolicy.Expire)
			} else if ppolicy.Grace >= 0 {
				logStr += fmt.Sprintf(". Password expired, %d grace logins remain\n", ppolicy.Grace)
			}
		}
		log.Print(logStr)
		log.Println(ppolicy)
	}

}
