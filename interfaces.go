package gopoet

type InterfaceSpec struct {
	// CodeBlock

	Name               string
	Comment            string
	EmbeddedInterfaces []TypeReference
	Methods            []FuncSpec
}

func NewInterfaceSpec(name string) *InterfaceSpec {
	return &InterfaceSpec{
		Name: name,
	}
}

func (i *InterfaceSpec) String() string {
	writer := NewCodeWriter()
	writer.WriteStatement(Statement{
		Format:      "type $L interface {",
		Arguments:   []interface{}{i.Name},
		AfterIndent: 1,
	})

	for _, interf := range i.EmbeddedInterfaces {
		writer.WriteStatement(Statement{
			Format:    "$L",
			Arguments: []interface{}{interf.GetName()},
		})
	}

	for _, method := range i.Methods {
		signature, args := method.Signature()
		writer.WriteStatement(Statement{
			Format:    signature,
			Arguments: args,
		})
	}

	writer.WriteStatement(Statement{
		Format:       "}",
		BeforeIndent: -1,
	})

	return writer.String()
}

func (i *InterfaceSpec) Packages() []Import {
	packages := []Import{}

	for _, method := range i.Methods {
		packages = append(packages, method.Packages()...)
	}

	for _, embedded := range i.EmbeddedInterfaces {
		packages = append(packages, embedded.GetImports()...)
	}

	return packages
}

func (i *InterfaceSpec) Method(spec FuncSpec) *InterfaceSpec {
	i.Methods = append(i.Methods, spec)
	return i
}

func (i *InterfaceSpec) EmbedInterface(interfaceType TypeReference) *InterfaceSpec {
	i.EmbeddedInterfaces = append(i.EmbeddedInterfaces, interfaceType)
	return i
}