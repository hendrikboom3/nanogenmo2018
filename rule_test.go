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
	for i := 0; i < len(rules); i++{
		fmt.Println("rule ", i, "is", rules[i])
	}
	if len(rules) == 0 { t.Error("norule"); return }
	// form, OK := (&lex{formtxt, 0}).readform(make(vars))
	// if ! OK { t.Error("noform"); return }
	for i := 0; i < len(rules); i++{
		tryrules(rules, rules[i], 2, func(){fmt.Println("done")})
	}
}

func xTestRules(t *testing.T){
	fmt.Println("testing rules.")
	testRules(t, "loves[b,c],loves[c,d]:kills[b,d];:loves[john[], mary[]];:loves[mary[],jim[]].")
	testRules(t, "loves[b,c],loves[c,d]:kills[b,d];:loves[john[], mary[]];:loves[mary[],john]].")
}

func TestMurder(t *testing.T){
	rule1 := "loves[b,c],loves[c,d],kills[b,d]:;"
        state := ":loves[john[], mary[]];:loves[mary[],jim[]];"
	intent := "findweapon[b, w]:kills[b, d];"
	state2 :=  ":findweapon[john[], gun[]]."
	testRules(t, rule1 + state + intent + state2)
}
