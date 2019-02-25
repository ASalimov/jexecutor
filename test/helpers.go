package test

import (
	"github.com/ASalimov/jexecutor/logic"
	"log"
)

func Test1() {
	//log.SetLevel(log.InfoLevel)
	log.Print("start...")
	URL := "https://jenkins.backend.pi.wuamerigo.com/"
	username := ""
	token := ""
	e := logic.NewExecutor(10, URL, username, token)
	KW := logic.Order{Job: "admin-rpm-deploy-autodeploy", Query: map[string]string{"COUNTRY": "kuwait", "ISO": "kw", "TAG": "1.61.0"}}
	BH := logic.Order{Job: "core-rpm-build-autodeploy", Query: map[string]string{"ZONE": "alpha", "TAG": "1.61.0"}}
	QA := logic.Order{Job: "gateway-rpm-build-autodeploy", Query: map[string]string{"ZONE": "alpha", "TAG": "1.61.0"}}
	JM := logic.Order{Job: "web-rpm-deploy-autodeploy", Query: map[string]string{"COUNTRY": "jamaica", "ISO": "jm", "TAG": "1.61.0"}}
	e.Execute(KW, BH, QA, JM)
}
