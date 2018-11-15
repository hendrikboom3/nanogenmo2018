package plotcalculus;

import (
	"bytes"
)

/*
This file contains a unification package -- that is a function taht takes two expressions and sees if there is a substitution for their variables that makes the two identical.

In formal descriptions of the process, the substitution is usually treated as a separate object form the two expressions.

Each distinct variable is represented once in memory.  If it occurs more often in the expression, all those occurrences point to the one object.  So a test for variable equality is smply a pointer equality test.

For applicatios, you will often want the variables in two expressions to be unified to be disjoint.  There is, therefore, a rename function provided that replaces all the variable in an expression with new variables that hapen to have the same names.
  THe unify function completely ifgnores variaabe names.  THey are there only for generating a slightly more readable form of debug output than might otherwiase be the case,
*/

type exp interface{
	String() string
	subst(substitution) exp
}

type exps []exp

type operator string

func (o vars)op(w string)operator { // TODO: an actual environment for oerators
	return operator(w)
}

type apply struct{
	fn operator  // TODO: not just integer, unless it's an index into some kind of operator table.  In which cse it shoule be a pointer to an operator struct -- and found in a table during reading.
	arg []exp
}

type variable struct{value exp
	name string
}


type substitution map[*variable]exp

func (s substitution) String() string{
	if s == nil { return "nosubst"
	} else {// https://www.calhoun.io/concatenating-and-building-strings-in-go/
		/* In later versions of Go, tey have introduced a string builder.
                   For our purposes, it behaves exacly like a bute buffer, which is already avaailable in earlier versions of Go.
		var sb strings.Builder
                */
		
		var sb bytes.Buffer
		sb.WriteString("[")
		first := true
		for key, value := range s {
			if ! first {sb.WriteString("; ")}
			first = false
			sb.WriteString(key.String())
			sb.WriteString("->")
			sb.WriteString(value.String())
		}
		sb.WriteString("]")
		return sb.String()
	}
}

/* An example gleaned from the net
   Ours has to be different in detail because we have static types and type tests in the impementation language, and eventually also in the expressions.

function UNIFY(x, y, theta) returns a substitution to make x and y identical
  inputs: x, a variable, constant, list, or compound expression
          y, a variable, constant, list, or compound expression
          theta, the substitution built up so far (optional, defaults to empty)

  if theta = failure then return failure
  else if x = y the return theta
  else if VARIABLE?(x) then return UNIFY-VAR(x, y, theta)
  else if VARIABLE?(y) then return UNIFY-VAR(y, x, theta)
  else if COMPOUND?(x) and COMPOUND?(y) then
      return UNIFY(x.ARGS, y.ARGS, UNIFY(x.OP, y.OP, theta))
  else if LIST?(x) and LIST?(y) then
      return UNIFY(x.REST, y.REST, UNIFY(x.FIRST, y.FIRST, theta))
  else return failure

---------------------------------------------------------------------------------------------------

function UNIFY-VAR(var, x, theta) returns a substitution

  if {var/val} E theta then return UNIFY(val, x, theta)
  else if {x/val} E theta then return UNIFY(var, val, theta)
  else if OCCUR-CHECK?(var, x) then return failure
  else return add {var/x} to theta

*/

func (e *apply)String()string{
	return string(e.fn) + format(e.arg)
}

func (e *variable)String()string{return e.name}

func format(e []exp)string{
	var sb bytes.Buffer
	sb.WriteString("[")
	for i := 0; i<len(e); i++ {
		if i > 0 { sb.WriteString(",") }
		sb.WriteString(e[i].String())
	}
	sb.WriteString("]")
	return sb.String()
}

func occurs(v *variable, e exp) bool {
	// Does v occur in e?
	switch e := e.(type) {
	case *apply:
		for i := 0; i < len(e.arg); i++ {
			if occurs(v, e.arg[i]) { return true }
		}
		return false;
	case *variable:
		return e == v
	}
	return false
}

func (e *variable)subst(s substitution) exp{
		val, ok := s[e]
		if ok {return val} else {return e}
}

func (e *apply)subst(s substitution) exp{
		var changed = false;
		nargs := make([]exp, len(e.arg))
		for i := 0; i < len(e.arg); i++ {
			nargs[i] = e.arg[i].subst(s) // subst(s, e.arg[i])
			changed = changed || nargs[i] != e.arg[i] 
		}
		if changed {
			return &apply{e.fn, nargs}
		} else {
			return e
			// If noting ie e has been substituted, we return e itself, instead of a copy of e.
		}
}

func subst(s substitution, e exp)exp {
	return e.subst(s)
}

func (s substitution)set(v *variable, val exp)substitution{
	// Set s[e] to val if e does not occur in val.
	if occurs(v, val) {
		return nil
	} else {
		s[v] = val  // TODO: What if val contins a variable that will substitute to v?
		return s
	}
}

func (s *substitution)get(v *variable)exp{
	return v.value
}

func unifyVar(s substitution, e *variable, f exp) substitution {
	eval, eOK := s[e]
	if eOK {
		return unify(s, eval, f)
	}
	switch f := f.(type) {
	case *variable:
		if e == f { return s } // trivial case, perhaps not needed.
		fval, fOK := s[f]
		if fOK {
			return unify(s, e, fval)   // TODO: What if fval contains a variable that will substitute to e?
		}
	}
	return s.set(e, f)
}

func unifies(s substitution, e []exp, f []exp) substitution{
	// Unify an array of expressions with another
	if len(e) != len(f) {return nil;}
	for i := 0; i < len(e); i++ {
		s = unify(s, e[i], f[i])
		if s == nil {
			return nil;
		}
	}
	return s
}

func unify(s substitution, e exp, f exp) substitution{
	// Distinct cases up to symmetry:
	// two so-far undefined variables
	// two variables, one with a value
	// Two variables, each with a value
	// an undefined variable and an application
	// a defined variable and an application
	// two applications
	if e == f { return s }
	switch e := e.(type) {
	case *apply:
		switch f := f.(type) {
		case *apply: // e and f both applications
			if e.fn == f.fn {
				return unifies(s, e.arg, f.arg)
			} else { return nil }
			
		case *variable: // e is application and f is a variable
			return unifyVar(s, f, e)
		}
	case *variable: // e is a variable
		return unifyVar(s, e, f)
	}
	return s // should never e exxecuted, but go copiler coplains
}

// TODO: name field no longer relevant, except for disgnostic printouts.
func newvariable(name string)*variable{
	return &variable{
		value: nil,
		name: name,
	}
}

func rename_variables(s substitution,  e exp)exp{
	// When called with an empty sunstitution, replace all variables in e by new variables.
	// This substitution is altered during the execution of this function.
	// and records tha changes made;  You usually don't want the substition afterward.
	switch e := e.(type){
	case *variable:
		eval, eok := s[e]
		if ! eok {
			s[e] = newvariable("?")
			return eval
		}
	case *apply:
		for _, a := range e.arg {
			rename_variables(s, a)
		}
	}
	return e // never executed byt go insists on it.
}

func tequal(e exp, f exp) bool {
	switch e := e.(type) {
	case *apply: 
		switch f := f.(type) {
		case *apply:
			if e.fn != f.fn { return false }
			if len(e.arg) != len(f.arg) { return false }
			for i := 0; i < len(e.arg); i++ {
				if ! tequal(e.arg[i], f.arg[i]) { return false }
			}
		
		case *variable:
			return false
		}
	
	case *variable:
		switch f := f.(type) {
		case *apply:
			return false
		case *variable:
			if e != f { return false }
		}
	
	}
	return true
}

