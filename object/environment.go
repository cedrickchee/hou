package object

// NewEnclosedEnvironment returns a new Environment with the outer set to the
// current environment (enclosing environment).
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// NewEnvironment constructs a new Environment object to hold bindings of
// identifiers to their names.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// Environment is what we use to keep track of value by associating them with a
// name. Technically, it's an object that holds a mapping of names to bound
// objets.
type Environment struct {
	store map[string]Object
	// outer is a reference to another Environment, which is the enclosing
	// environment, the one it’s extending.
	outer *Environment
}

// Get returns the object bound by name.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		// Check the enclosing environment for the given name.
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set stores the object with the given name.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
