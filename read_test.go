package plotcalculus;

import (
	"testing"
	"fmt"
)

func TestRead(t *testing.T){
	str := "a[k],e[k]:b[u],c[k]"
	fmt.Println("testing read.", str)
	rr := (&lex{str, 0}).readrules();
	if len(rr) == 0 {t.Error("morules")
	} else{ fmt.Println(rr) }
}
