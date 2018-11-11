package unify

import (
	"testing"
	"fmt"
)

/*
func TestRule(t *testing.T){
	fmt.Println("testing rule.");
	ruletxt := "a[b]:c[b]"
	formtxt := "a[d]"
	rule := (&lex{ruletxt, 0}).readrule()
	if rule == nil { t.Error("norule"); return }
	form, OK := (&lex{formtxt, 0}).readform(make(vars))
	if ! OK { t.Error("noform"); return }
	userule(form, rule)
}
*/

func testRules(t *testing.T, ruletxt string){
	rules := (&lex{ruletxt, 0}).readrules()
	fmt.Println("read", len(rules), "rules")
	if len(rules) == 0 { t.Error("norule"); return }
	// form, OK := (&lex{formtxt, 0}).readform(make(vars))
	// if ! OK { t.Error("noform"); return }
	tryrules(rules, rules[0], 1)
}

func TestRules(t *testing.T){
	fmt.Println("testing rules.")
	testRules(t, "loves[b,c],loves[c,d]:kills[b,d];:loves[john[], mary[]];:loves[mary[],jim[]].")
	testRules(t, "loves[b,c],loves[c,d]:kills[b,d];:loves[john[], mary[]];:loves[mary[],john]].")
}
