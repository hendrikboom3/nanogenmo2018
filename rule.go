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

func (ee exps)subst(s substitution)exps{
	changed := false
	nee   := make([]exp, len(ee))
	for i := 0; i < len(ee); i++ {
		nee[i] = ee[i].subst(s) // subst(s, e.arg[i])
			changed = changed || nee[i] != ee[i] 
	}
	if changed {return nee } else { return ee }
}

func (s substitution) copy ()substitution{
	n := make(substitution)
	for k,v := range s {
		n[k] = v
	}
	return n
}

func rulesnest(rules[]*rule, a []exp, ai int, ss substitution, maxdepth int, after func(ns substitution)){
	// match the antecedents of a[ai] with all the consequents of the rules to find some that might extablish the preconditions for r.
	// For each success, recurse to the remaining elements of a, gradually accumulating a sunstitution as we go.
	// If there are no more, call after with the cumulative substitution as argument.
	if ai >= len(a){
		// fmt.Println(maxdepth, "AFTER!!")
		after(ss)
		return
	}
	for i :=0; i < len(rules); i++ {
		for j := 0; j < len(rules[i].consequent); j++ {
			// fmt.Println(maxdepth, "trying to unify", a[ai], " with ", rules[i].consequent[j], " after ", ss)
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


func proveoneantecedent(ante exp, substsofar substitution, fanout func(func(exp, *rule, int)), myact func(substsofar substitution, report func()), depth int){
	// process one antecedent by confronting it with every consequent of every rule, as povided by fanout.  For each such success, call act.
	fanout(
		(func(c exp, rr *rule, index int){ // Maybe rule or even consequent index should also be parameters?
			// fmt.Println(depth, "confront ", ante, "with", c)
			newsubst := unify(substsofar.copy(), ante, c)
			if newsubst == nil { return }
			{ // We have to handle premisses of this new rule.
				// fmt.Println(depth, "Consider rules for antecedents for ", ante  )
				nest(rr.antecedent, 0, newsubst, fanout,
					func(endsubst substitution, report func()){
						// TODO: worry about bound rule-variable capture.
						//       I suspect we need to rename all the variables in the nested rule before unifying.

						// TODO: I suspect we want to strip endsubst of all the items that substitute names defined in the nested rule.
						//       This probably means completely substituting them out of ante
						//       and then unifying substituted ante with original ante in context ot sustitution prior to ante.
						// s := unify(substsofar, ante, newante)
						myact(endsubst, func(){
							report()
							/* fmt.Println(" ||| ", subst(endsubst, ante)) */
						})
						// TODO: process the end substitution in some fashiom to use newly discovered intofmation
					},
					depth - 1)
				// fmt.Println(depth, "finishing rules for antecedents for ", ante)
			}
		}))
}

	
func nest(aa exps, i int, substsofar substitution, fanout func(func(exp, *rule, int)), act func(substsofar substitution, report func()), depth int) { // process antecentents  aa from i on
	// by confronting every exp in aa with every consequent of every rule, as performed by fanout.
	// Call act whenever you get to the end of aa.

	if depth == 0 { fmt.Println("recursion limit reached"); return}

	if i >= len(aa) {
		act(substsofar, func(){ /* fmt.Println("++++", aa.subst(substsofar))*/ })
		return
	}

	proveoneantecedent(aa[i], substsofar, fanout,
		/*myact*/func(endsubst substitution, report1 func()){
			nest(aa, i+1, endsubst, fanout,
				func(substsofar substitution, report2 func()){
					act(substsofar, func(){
						report1();
						fmt.Println("^^^", subst(substsofar, aa[i]));
						report2()});
					
				}, depth)
		},
		depth)
}


func forcons(rules[]*rule, act func(exp, *rule, int)){ //  Apply ACT to every consequent of any rule in rules. 
	for i :=0; i < len(rules); i++ {
		for j := 0; j < len(rules[i].consequent); j++ {
			act(rules[i].consequent[j], rules[i], j)
		}
	}
}

func tryrules(rules[]*rule, aa exps, depth int, success func()){
	newsubst := make( substitution)
	nest(aa, 0, newsubst,
		/* fanout */(func(what func(c exp, r *rule, index int)){
			forcons(rules, what)}),
		/* act */ (func(substsofar substitution, report func()){
			// fmt.Println(depth, "antecedents", aa, "subst", substsofar, "yields", aa.subst(substsofar))
			report()
			success()
		}),
		depth /* arbitrary recursion depth limit */ )
}


func oldtryrules(rules[]*rule, rr *rule, depth int, success func()){
	// match the antecedents of r with all the consequents of the rules to find some that might extablish the preconditions for r.
	// This means nested iteration, nested as deeply as the number of antecedents.
	// It also means recursive calle to oldtryrules with the antecedents of those rules.  This will be done as far as depth recursions. 
	if len(rr.antecedent) == 0 {
		fmt.Println("At start?end", rr)
		success();
		return
	}
	fmt.Println(depth, "process antecedents of", rr)
	if depth == 0 {
		fmt.Println("depth limit reached.")
		return 
	}
	for i :=0; i < len(rules); i++ {
		for j := 0; j < len(rules[i].consequent); j++ {
			s := make(substitution)
			fmt.Println(depth, "use rule ", rr, "on consequent" , rules[i].consequent[j], "of rule", rules[i])
			rulesnest(rules, rr.antecedent, 0, s, depth,
				(func(ns substitution){
					for i := 0; i < len(rr.consequent); i++ {
						fmt.Println(depth, "work on consequent", rr.consequent[i], "of rule", rr)
						oldtryrules(rules, rules[i], depth-1,
							func(){
								fmt.Println(depth, "success at depth", depth)
								fmt.Println("doing consequent", rr.consequent[i], "of rule", rr)
								success() // further out
							})
					}
				}))
		}	
	}
}
