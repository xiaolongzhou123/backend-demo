package ldap

import (
	"errors"
	"fmt"
	"sso/pkg"
	"strings"

	"github.com/go-ldap/ldap"
)

type LdapClient struct {
	Name   string
	Passwd string
	Conn   *ldap.Conn
}

func NewLDAP() (*LdapClient, error) {
	conf := pkg.Conf()
	url := fmt.Sprintf("%s", conf.Ldap.Host)
	fmt.Println(url)
	l, err := ldap.Dial("tcp", url)
	if err != nil {
		fmt.Println("ldap Dial err:=", err)
		return nil, err
	}
	err = l.Bind(conf.Ldap.BindDN, conf.Ldap.BinPass)
	if err != nil {
		fmt.Println("ldap Bind err:=", err)
		return nil, err
	}

	return &LdapClient{
		Conn: l,
	}, nil

}
func (this *LdapClient) Close() {
	this.Conn.Close()
}
func (this *LdapClient) GetUser(name string) (map[string]interface{}, error) {
	// config := g.Config()
	conf := pkg.Conf()
	searchdn := fmt.Sprintf("%s,%s", conf.Ldap.UserDN["signup"], conf.Ldap.BaseDN)
	request := ldap.NewSearchRequest(
		searchdn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf("(&(cn=%s))", name),
		conf.Ldap.LdapAttr,
		nil)

	sr, err := this.Conn.Search(request)
	if err != nil {
		fmt.Println("user search err", err)
		return nil, err
	}

	list := GetResult(sr)
	m := make(map[string]interface{}) //这个是空的数据

	if len(list) == 0 {
		return make(map[string]interface{}), errors.New("未查找到用户")
	}

	for _, v := range list {
		t := make(map[string]interface{})
		for _, vv := range v.Data {
			t[vv.Name] = vv.Value
		}
		return t, nil
	}
	return m, nil
}

func (this *LdapClient) GetUsers() ([]map[string]interface{}, error) {
	conf := pkg.Conf()
	mm := make([]map[string]interface{}, 0)
	searchdn := fmt.Sprintf("%s,%s", conf.Ldap.UserDN["signup"], conf.Ldap.BaseDN)
	request := ldap.NewSearchRequest(
		searchdn,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		"(&(objectClass=inetOrgPerson))",
		conf.Ldap.LdapAttr,
		nil)

	sr, err := this.Conn.Search(request)
	if err != nil {
		fmt.Println("user search err", err)
		return mm, err
	}
	list := GetResult(sr)
	if len(list) == 0 {
		return mm, nil
	}
	for _, v := range list {
		m := make(map[string]interface{})
		for _, vv := range v.Data {
			m[vv.Name] = vv.Value
		}
		mm = append(mm, m)
	}
	return mm, nil
}

type Result struct {
	Name  string
	Value string
}
type List struct {
	Index int
	Data  []*Result
}

func GetResult(sr *ldap.SearchResult) []*List {
	// config := g.Config()
	arrList := make([]*List, 0)
	for k, one := range sr.Entries {
		arr := make([]*Result, 0)
		// for _, name := range config.GlobalConfig.Ldap.LdapAttr {
		for _, atr := range one.Attributes {
			rs := &Result{}
			rs.Value = one.GetAttributeValue(atr.Name)
			rs.Name = atr.Name
			arr = append(arr, rs)
		}
		list := &List{}
		list.Index = k
		list.Data = arr
		arrList = append(arrList, list)
	}
	return arrList
}
func (l *LdapClient) GetUserNameCN(name, searchdn string) (string, error) {
	req := ldap.NewSearchRequest(searchdn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(cn=%s))", name), // The filter to apply
		[]string{"entryDN"}, nil)
	res, err := l.Conn.Search(req)
	if err != nil {
		return "", err
	}

	var entry *ldap.Entry
	for _, entry = range res.Entries {
		break
	}
	if entry == nil {
		return "", nil
	}
	return entry.DN, nil
}
func (this *LdapClient) Login(user, pass string) error {
	// config := g.Config()
	conf := pkg.Conf()
	searchdn := fmt.Sprintf("%s,%s", conf.Ldap.UserDN["signup"], conf.Ldap.BaseDN)
	fmt.Println(searchdn)
	dn, err := this.GetUserNameCN(user, searchdn)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if dn == "" {
		return errors.New("user isn't exist.")
	}
	err = this.Conn.Bind(dn, strings.TrimSpace(pass))
	return err
}
