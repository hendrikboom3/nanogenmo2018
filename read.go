package unify
import(
	"io/ioutil"
	"fmt"
//	"os"
)

type lex struct{
	s string
	pos int
}

type vars map[string]*variable // TODO: peroper name for this type
// TODO: proper implementation of operator lookup, too
// vars shuold probably be defined in unify.

func (l *lex) more() bool { return l.pos < len(l.s) }

func (l *lex) peek() byte{
	if l.pos < len(l.s) {
		return l.s[l.pos]
	} else {
		return '\xff'
	}
}

func err(m string) {
	fmt.Println(m);
}

func (l *lex) advance(){
	// DEBUG fmt.Printf("advance at %d %c %d.\n", l.pos, l.s[l.pos], len(l.s))
	if l.pos < len(l.s) { l.pos++ } else { err("eof"); panic("eof")}
	for l.pos < len(l.s) &&
		(l.s[l.pos] == ' ' ||
		l.s[l.pos] == '\t' ||
		l.s[l.pos] == '\n' ||
		l.s[l.pos] == '\r'){
		// DEBUG fmt.Printf("advance now at %d %c %d.\n", l.pos, l.s[l.pos], len(l.s))
		l.pos++ 
		// DEBUG fmt.Printf("advanced to %d %c %d.\n", l.pos, l.s[l.pos], len(l.s))		// todo: extract runes, not chars??
	}
}

func isletter(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

func isdigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// do something about reserved words?

func (ll *lex)readword() (string, bool){
	start := ll.pos;
	if ll.more() && isletter(ll.peek()) {
		for ll.more() && (isletter(ll.peek()) || isdigit(ll.peek())) {
			ll.advance()
		}
		return ll.s[start:ll.pos], true // actually, accumulates string
	}
	return "[ERROR]", false
}

func(env vars)lookup(name string) *variable { // TODO: an environment for functions andpredicates as well
	v := env[name]
	if v == nil {
		v = newvariable(name)
		env[name] = v
	}
	return v 
}

func (ll *lex) readform(env vars)(exp, bool){
	w, OK := ll.readword()
	if OK {
		if ll.more() && ll.peek() == '[' {
			ll.advance()
			a, OK := ll.readforms(env)
			if ! OK { return nil, false }
				if ll.more() && ll.peek() == ']'{
				ll.advance()
					return &apply{env.op(w), a}, true //TODO use proper operator query
			}
		}
		v := env.lookup(w)
		return v, true // TODO: variable
	} else { return nil, false }
}

func (ll *lex)readforms(env vars)([]exp, bool){
	allOK := true
	var list []exp
	a, OK := ll.readform (env) 
	if OK {
		list = append(list, a)
		for ll.peek() == ',' {
			ll.advance()
			a, OK:= ll.readform(env)
			if !OK {
				err("foo")
				allOK = false
				break;
			}
			list = append(list, a)
		}
	} else { allOK = true /* empty list */  }
	return list, allOK
}

func (ll *lex)readantecedents(env vars)([]exp, bool){
	return ll.readforms(env)
}

/*
func (ll *lex)readantecedents() (exp, bool){
	allOK := true
	a, OK := ll.readantecedent()
	if OK {
		for ll.peek() == ',' {
			a, OK:= ll.readantecedent()
			if !OK {
				err("foo")
				allOK := false
				break;
			}
		}
		return nil, allOK // actually, catenate them
	}
	return nil, false
}
*/

func (ll *lex)readconsequents(env vars) ([]exp, bool) {
	return ll.readforms(env)
}

func (ll *lex)readrule() *rule{
	env := make(vars)
	a, OK := ll.readantecedents (env)
	if OK {
		// DEBUG fmt.Println("Now at char ", ll.pos, "namely", string(rune(ll.peek()))) 
		if ll.peek() == ':' {
			ll.advance();
			// DEBUG fmt.Println("Next at char ", ll.pos, "namely", string(rune(ll.peek())))
 			c, OK := ll.readconsequents(env)
			fmt.Println("OK?", OK, "pos", ll.pos)
			if OK {
				return &rule{a, c}
			}
		}
	}
	return nil
}

func (ll *lex)readrules()[]*rule{
	// TODO: readrules should have an end marker
	var rr []*rule
	r := ll.readrule()
	if r == nil {
		fmt.Println("No rule.")
		return rr
	}
	rr = append(rr, r)
	for r != nil && ll.peek() == ';' {
		ll.advance()
		r = ll.readrule()
		rr = append(rr, r)
		//  fmt.Println("Expected rule.")
	}
	fmt.Printf("got to position %d finding %d rules.\n", ll.pos, len(rr))
	
	return rr
}

func read(filename string){
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fileStr := string(fileBytes)
	ll := lex{fileStr, 0}
	ll.readrules()
}
