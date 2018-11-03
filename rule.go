package unify;
import(
	"fmt"
)

type rule struct{
	antecedent []exp
	consequent []exp
}

type state []exp

func act(e exp){
	fmt.Println("do ", e.String())
}

func userule(e exp, r rule){
	s := make(substitution)
	unify(s, e, r.antecedent[0])
	for i := 0; i < len(r.consequent); i++ {
		act(r.consequent[i])
	}
}
