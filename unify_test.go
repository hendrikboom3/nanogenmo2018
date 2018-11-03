package unify;

import ("testing"
)

func stringify(s bool) string{
	if s { return "true" } else { return "false" }
}

func tunify(t *testing.T, e exp, f exp, r bool) substitution {
	// fmt.Println("unify ", e, " and ", f, " expect ", r);
	s := make(substitution)
	s = unify(s, e, f)
	message := "unify " + e.String() + " and " + f.String() + " expecting " + stringify(r) + " gives " + s.String() + " yielding " + subst(s, e).String() + " and " + subst(s, f).String()
	// fmt.Println("message", message)
	if (s != nil) != r { t.Error("wrong result: " + message)
		// fmt.Println("   result not ", r, message)
	}
	if s != nil {
		// fmt.Println("   substitution ", s)
		if ! tequal(subst(s, e), subst(s, f)) {
			t.Error("   unequal after subst: " + message)
		}
		// fmt.Println("   giving ", subst(s, e), " or ", subst(s, f))
	}
	return s;
}

func TestUnify(t *testing.T){
	v := newvariable("v")
	w := newvariable("w")
	if v.String() != "v" {t.Error("onion")}
	if w.String() != "w" {t.Fail()}
	s := tunify(t, v, w, true)
	if s == nil {t.Error("bar")}
//	if subst(s, v) != nil {t.Error("vvalue not nil")}
//	if subst(s, w) == nil {t.Error("wvalue nil")}
	if subst(s, w) != subst(s, v) {t.Error("wvalue not v")}
	
}

func TestFun(t *testing.T){
	v := newvariable("v")
	w := newvariable("w")
	b := &apply{2, []exp{w}}
	a := &apply{2, []exp{v}}
	/* s := */ tunify(t, a, b, true)
	/* if s == nil {t.Error("a bot b")}
	_, vok := s[v]
	if ! vok {t.Error("avalue missing")}
	wval, wok := s[w]
	if ! wok  {t.Error("bvalue nil")}
	if wval != v {t.Error("awvalue not v")}
        */
}

func TestVarfun(t *testing.T){
	//	t.Error("test whether faiing works.")
	v := newvariable("v")
	w := newvariable("w")
	b := &apply{2, []exp{w}}
	a := &apply{2, []exp{v}}
	/* s := */ tunify(t, a, v, false)
/*	if s == nil {t.Error("exp np tunify with var")}
	if ! tequal(subst(s, a), subst(s, b)) {
		t.Error("substitoin does not unify")
	}
*/
	tunify(t, w, b, false)/* == nil {t.Error("var np tunify with exp")}*/
	c := &apply{3, []exp{w}}
	_ = tunify(t, a, c, false) /* != nil {t.Error("unified with different operators")} */
	
}
