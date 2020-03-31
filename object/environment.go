package object

// NewEnvironment constructs a new Environment object to hold bindings of
// identifiers to their names.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Environment is what we use to keep track of value by associating them with a
// name. Technically, it's an object that holds a mapping of names to bound
// objets.
type Environment struct {
	store map[string]Object
}

// Get returns the object bound by name.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set stores the object with the given name.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
