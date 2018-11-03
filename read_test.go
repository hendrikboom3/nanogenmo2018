package unify;

import (
	"testing"
	"fmt"
)

func TestRead(t *testing.T){
	fmt.Println("testing read.")
	rr := (&lex{"a[],e[]:b[],c[]", 0}).readrules();
	if len(rr) == 0 {t.Error("morules")}
}
