package ldap

import (
	"errors"
	"fmt"
	"log"
	"sso/pkg"
	"sso/pkg/typing"
	"sso/pkg/utils"
	"strings"

	"github.com/go-ldap/ldap/v3"
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
	err = l.Bind(conf.Ldap.BindDN, conf.Ldap.BindPass)
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

func (this *LdapClient) GetUserValue(name string) (map[string]interface{}, error) {
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
func (l *LdapClient) GetUserNameCN(name string) (string, error) {

	conf := pkg.Conf()
	searchdn := fmt.Sprintf("%s,%s", conf.Ldap.UserDN["signup"], conf.Ldap.BaseDN)
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

func (this *LdapClient) LoginV1(userdn, pass string) error {
	// config := g.Config()
	err := this.Conn.Bind(userdn, strings.TrimSpace(pass))
	return err
}

func (this *LdapClient) Login(user, pass string) (error, string) {
	conf := pkg.Conf()
	searchdn := fmt.Sprintf("%s,%s", conf.Ldap.UserDN["signup"], conf.Ldap.BaseDN)
	//主要为了查询userdn,cn=abc,cn=dev,ou=tech,dc=example,dc=org
	req := ldap.NewSearchRequest(searchdn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(cn=%s))", user), // The filter to apply
		[]string{"entryDN"}, nil)

	res, err := this.Conn.Search(req)
	if err != nil {
		return err, ""
	}
	if len(res.Entries) == 0 {
		return fmt.Errorf("用户不存在"), ""
	}
	userdn := res.Entries[0].DN
	fmt.Println("====", userdn)

	// userdn := "cn=abc,cn=dev,ou=tech,dc=example,dc=org"

	controls := []ldap.Control{}
	controls = append(controls, ldap.NewControlBeheraPasswordPolicy())
	bindRequest := ldap.NewSimpleBindRequest(userdn, pass, controls)

	r, err := this.Conn.SimpleBind(bindRequest)
	ppolicyControl := ldap.FindControl(r.Controls, ldap.ControlTypeBeheraPasswordPolicy)

	var ppolicy *ldap.ControlBeheraPasswordPolicy
	if ppolicyControl != nil {
		ppolicy = ppolicyControl.(*ldap.ControlBeheraPasswordPolicy)
	} else {
		log.Printf("ppolicyControl response not avaliable.\n")
	}
	if err != nil {
		errStr := "ERROR: Cannot bind: " + err.Error()
		if ppolicy != nil && ppolicy.Error >= 0 {
			errStr += ":" + ppolicy.ErrorString
		}
		return err, errStr
	} else {
		logStr := "Login Ok"
		if ppolicy != nil {
			if ppolicy.Expire >= 0 {
				logStr += fmt.Sprintf(". Password expires in %d seconds\n", ppolicy.Expire)
			} else if ppolicy.Grace >= 0 {
				logStr += fmt.Sprintf(". Password expired, %d grace logins remain\n", ppolicy.Grace)
			}
		}
		return nil, logStr
	}
	return errors.New("unknow error"), ""
}
func (l *LdapClient) ChangePasswd(dn, old_pass, cur_pass string) error {
	pmr := ldap.NewPasswordModifyRequest(dn, old_pass, cur_pass)
	_, err := l.Conn.PasswordModify(pmr)
	if err != nil {
		return err
	}
	return nil
}

func (l *LdapClient) AddUser(userdn string, user typing.LdapUser) error {
	conf := pkg.Conf()
	req := &ldap.AddRequest{
		DN:       userdn,
		Controls: nil,
	}

	req.Attribute("objectClass", []string{"myObjectClass", "top"})
	req.Attribute("cn", []string{user.Cn})
	req.Attribute("sn", []string{user.Sn})
	req.Attribute("myId", []string{user.MyId})
	req.Attribute("myName", []string{user.MyName})
	req.Attribute("myPhone", []string{user.MyPhone})
	req.Attribute("myTel", []string{user.MyTel})
	req.Attribute("myEmail", []string{user.MyEmail})
	req.Attribute("myGoogle", []string{utils.CreatePrivKey()})
	req.Attribute("myLeader", []string{user.MyLeader})
	req.Attribute("myPostion", []string{user.MyPostion})
	req.Attribute("myReg", []string{user.MyReg})
	req.Attribute("myDep", []string{user.MyDep})
	req.Attribute("myCompanyGroup", []string{user.MyCompanyGroup})
	if user.MyGender {
		req.Attribute("myCustom1", []string{"true"})
	} else {
		req.Attribute("myCustom1", []string{"false"})
	}
	req.Attribute("UserPassword", []string{conf.Ldap.DefaultPass})

	fmt.Println(req.Attributes)
	err := l.Conn.Add(req)
	if err != nil {
		fmt.Println("ldap adduser,conn add err:", err)
		return err
	}

	return nil
}

func (l *LdapClient) ModifyUser(userdn string, user typing.LdapUser) error {
	conf := pkg.Conf()
	req := &ldap.ModifyRequest{
		DN:       userdn,
		Controls: nil,
	}

	req.Replace("cn", []string{user.Cn})
	req.Replace("sn", []string{user.Sn})
	req.Replace("myName", []string{user.MyName})
	req.Replace("myPhone", []string{user.MyPhone})
	req.Replace("myTel", []string{user.MyTel})
	req.Replace("myEmail", []string{user.MyEmail})
	req.Replace("myGoogle", []string{user.MyGoogle})
	req.Replace("myLeader", []string{user.MyLeader})
	req.Replace("myPostion", []string{user.MyPostion})
	req.Replace("myReg", []string{user.MyReg})
	req.Replace("myDep", []string{user.MyDep})
	req.Replace("myCompanyGroup", []string{user.MyCompanyGroup})
	if user.MyGender {
		req.Replace("myCustom1", []string{"true"})
	} else {
		req.Replace("myCustom1", []string{"false"})
	}
	req.Replace("UserPassword", []string{conf.Ldap.DefaultPass})

	// fmt.Println(req.Attributes)
	err := l.Conn.Modify(req)
	if err != nil {
		fmt.Println("ldap adduser,conn add err:", err)
		return err
	}

	return nil
}

func (l *LdapClient) DelUser(userdn string) error {
	req := &ldap.DelRequest{
		DN:       userdn,
		Controls: nil,
	}

	// fmt.Println(req.Attributes)
	err := l.Conn.Del(req)
	if err != nil {
		fmt.Println("ldap adduser,conn add err:", err)
		return err
	}

	return nil
}
