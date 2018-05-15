package template

import (
	_ "ainit"
	"log"
	"testing"
)

func TestTemplateList(t *testing.T) {

	if l := GetTemplateTags(); len(l) <= 0 {
		t.Fail()
	}
	GetTemplateContent("default")
	log.Println("TestTemplateList:OK")
}

func TestTemplateDesc(t *testing.T) {

	if l := GetTemplateTags(); len(l) <= 0 {
		t.Fail()
	}
	log.Println("TestTemplateCount:OK")
	log.Println(MergeContent(55, 1))
}
