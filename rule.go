package unify;
import(
	"fmt"
	"strings"
)

type rule struct{
	antecedent []exp
	consequent []exp
}

type state []exp

func act(e exp){
	fmt.Println("do ", e.String())
}

func (r rule)String()string{
	return format(r.antecedent) + ":" + format(r.consequent)
}

func (s substitution) copy ()substitution{
	n := make(substitution)
	for k,v := range s {
		n[k] = v
	}
	return n
}

func rulesnest(rules[]*rule, a []exp, ai int, ss substitution, depth int, c []exp){
	for i :=0; i < len(rules); i++ {
		for j := 0; j < len(rules[i].consequent); j++ {
			ns := unify(ss.copy(), a[ai], rules[i].consequent[j])
			if(ns != nil){ // have a possibiity
				fmt.Println(strings.Repeat("  ", ai), "subst is", ns.String())
				if ai+1 < len(a) {
					rulesnest(rules, a, ai + 1, ns, depth, c)
				} else {
					if len(rules[i].antecedent) == 0 { // achieved
						for i := 0; i < len(c); i++ {
							act(subst(ns, c[i]))
						}
					}
				}
			}	
		}
		
	}
}

func tryrules(rules[]*rule, r *rule, depth int){
	// match the antecedents of r with all the consequents of the rules to find some that might extablish the preconditions for r.
	// Still the case of only one antecedent in r... We nwwd o find rules set that work compatibly for All the antecedents, not just one.
	// This means nested iteration, nested as deeply as the number of antecedents.
	
	for i :=0; i < len(rules); i++ {
		for j := 0; j < len(rules[i].consequent); j++ {
			s := make(substitution)
			fmt.Println("use rule ", r, "on" , rules[i].consequent[j])
			rulesnest(rules, r.antecedent, 0, s, 0, r.consequent)
		}	
	}
}
