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
	if l.pos < len(l.s) { l.pos++ } else { err("eof"); panic("eof")}
	for l.pos < len(l.s) && (l.s[l.pos] == ' ' || l.s[l.pos] == '\t' || l.pos == '\n' || l.s[l.pos] == '\r'){ l.pos++ }
//	fmt.Printf("advanced to %d %c %d.\n", l.pos, l.s[l.pos], len(l.s))
			// todo: extract runes, not chars??
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
		return ll.s[start:ll.pos-1], true // actually, accumulates string
	}
	return "[ERROR]", false
}

func (ll *lex) readform()(exp, bool){
	w, OK := ll.readword()
	if OK {
		if ll.more() && ll.peek() == '[' {
			ll.advance()
			a, OK := ll.readforms()
			if ! OK { return nil, false }
				if ll.more() && ll.peek() == ']'{
				ll.advance()
				return &apply{7, a}, true //TODO 7 should be operator
			}
		}
		return &variable{nil, w}, true // TODO: variable
	} else { return nil, false }
}

func (ll *lex)readforms()([]exp, bool){
	allOK := true
	var list []exp
	a, OK := ll.readform () 
	if OK {
		list = append(list, a)
		for ll.peek() == ',' {
			ll.advance()
			a, OK:= ll.readform()
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

func (ll *lex)readantecedents()([]exp, bool){
	return ll.readforms()
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

func (ll *lex)readconsequents() ([]exp, bool) {
	return ll.readantecedents()
}

func (ll *lex)readrule() *rule{
	a, OK := ll.readantecedents ()
	if OK {
		if ll.peek() == ':' {
			ll.advance();
			c, OK := ll.readconsequents()
			if OK {
				return &rule{a, c}
			}
		}
	}
	return nil
}

func (ll lex)readrules()[]*rule{
	var rr []*rule
	r := (&ll).readrule()
	for r != nil {
		rr = append(rr, r)
		r = ll.readrule()
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
	l := lex{fileStr, 0}
	l.readrule() // should be rules.  Need end marker.
}
