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

func rulesnest(rules[]*rule, a []exp, ai int, ss substitution, maxdepth int, after func(ns substitution)){
	// match the antecedents of a[i] with all the consequents of the rules to find some that might extablish the preconditions for r.
	// For each success, recurse to the remaining elements of a, gradually accumulating a sunstitution as we go.
	// If there are no more, call after with the cumulative substitution as argument.
	if ai >= len(a){
		fmt.Println(maxdepth, "AFTER!!")
		after(ss)
		return
	}
	for i :=0; i < len(rules); i++ {
		for j := 0; j < len(rules[i].consequent); j++ {
			// fmt.Println(maxdepth, "utrying to unify", a[ai], " with ", rules[i].consequent[j], " after ", ss)
			ns := unify(ss.copy(), a[ai], rules[i].consequent[j])
			if(ns != nil){ // have a possibiity
				fmt.Println(maxdepth, strings.Repeat("  ", ai), "subst is", ns.String())
				rulesnest(rules, a, ai + 1, ns, maxdepth,
					func(ns substitution){
						fmt.Println(maxdepth, "action ", ai+1, "of", subst(ns, a[ai]))})
			// } else { fmt.Println(maxdepth, "unify filed.")
			}
		}
		
	}
}


func tryrules(rules[]*rule, r *rule, depth int, success func()){
	// match the antecedents of r with all the consequents of the rules to find some that might extablish the preconditions for r.
	// This means nested iteration, nested as deeply as the number of antecedents.
	// It also means recursive calle to tryrules with the antecedents of those rules.  This will be done as far as depth recursions. 
	if depth == 0 {
		fmt.Println("depth limit reached.")
		return 
		//act(r.consequent[i]) // probably wrong action
	}
	fmt.Println(depth, "process antecedents of", r)
	for i :=0; i < len(rules); i++ {
		for j := 0; j < len(rules[i].consequent); j++ {
			s := make(substitution)
			fmt.Println(depth, "use rule ", r, "on" , rules[i].consequent[j])
			rulesnest(rules, r.antecedent, 0, s, depth,
				(func(ns substitution){
					if len(r.consequent) == 0 { // is this the right test for success?
						success()
						fmt.Println(depth, "success?", r)
					}
					for i := 0; i < len(r.consequent); i++ {
						fmt.Println(depth, "work on consequent", r.consequent[i], "of rule", r)
						tryrules(rules, rules[i], depth-1,
							func(){
								fmt.Println(depth, "success at depth", depth)
								success()
							})
						// act(subst(ns, r.consequent[i]))
					}
				}))
		}	
	}
}
