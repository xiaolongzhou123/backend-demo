package typing

type LdapUser struct {
	Cn             string `json:"cn"`
	Sn             string `json:"sn"`
	MyId           string `json:"myId"`
	MyName         string `json:"myName"`
	MyPhone        string `json:"myPhone"`
	MyTel          string `json:"myTel"`
	MyEmail        string `json:"myEmail"`
	MyGoogle       string `json:"myGoogle"`
	MyLeader       string `json:"myLeader"`
	MyPostion      string `json:"myPostion"`
	MyReg          string `json:"myReg"`
	MyDep          string `json:"myDep"`
	MyCompanyGroup string `json:"myCompanyGroup"`
	MyGender       bool   `json:"myGender"`
}
