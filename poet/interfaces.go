package poet

var _ CodeBlock = (*InterfaceSpec)(nil)
var _ TypeReference = (*InterfaceSpec)(nil)

// InterfaceSpec represents an interface
type InterfaceSpec struct {
	Name               string
	Comment            string
	EmbeddedInterfaces []TypeReference
	Methods            []*FuncSpec
}

// NewInterfaceSpec constructs a new interface with the given name
func NewInterfaceSpec(name string) *InterfaceSpec {
	return &InterfaceSpec{
		Name: name,
	}
}

// InterfaceComment sets the comment for the interface.
func (i *InterfaceSpec) InterfaceComment(comment string) *InterfaceSpec {
	i.Comment = comment
	return i
}

// Method adds a new method to the interface
func (i *InterfaceSpec) Method(spec *FuncSpec) *InterfaceSpec {
	i.Methods = append(i.Methods, spec)
	return i
}

// EmbedInterface specifies an interface to embed in the interface
func (i *InterfaceSpec) EmbedInterface(interfaceType TypeReference) *InterfaceSpec {
	i.EmbeddedInterfaces = append(i.EmbeddedInterfaces, interfaceType)
	return i
}

// GetImports returns Imports used by the interface
func (i *InterfaceSpec) GetImports() []Import {
	packages := []Import{}

	for _, method := range i.Methods {
		packages = append(packages, method.GetImports()...)
	}

	for _, embedded := range i.EmbeddedInterfaces {
		packages = append(packages, embedded.GetImports()...)
	}

	return packages
}

// GetName returns the name and fulfills TypeReference.
func (i *InterfaceSpec) GetName() string {
	return i.Name
}

// String returns a string representation of the interface
func (i *InterfaceSpec) String() string {
	w := newCodeWriter()
	w.WriteCodeBlock(i)
	return w.String()
}

// GetStatements returns the statements representing the interface declaration and
// all of its methods.
func (i *InterfaceSpec) GetStatements() []Statement {
	var s []Statement
	s = append(s, Comment(i.Comment).GetStatements()...)
	s = append(s, newStatement(0, 1, "type $L interface {", i.Name))

	for _, interf := range i.EmbeddedInterfaces {
		s = append(s, newStatement(0, 0, "$L", interf.GetName()))
	}

	for _, m := range i.Methods {
		s = append(s, Comment(m.Comment).GetStatements()...)
		signature, args := m.Signature()
		s = append(s, newStatement(0, 0, signature, args...))
	}

	s = append(s, newStatement(-1, 0, "}"))
	return s
}
