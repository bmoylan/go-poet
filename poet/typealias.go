package poet

var _ CodeBlock = (*TypeAliasSpec)(nil)
var _ TypeReference = (*TypeAliasSpec)(nil)

// TypeAliasSpec represents a type alias. *AliasSpec implements CodeBlock and TypeReference.
type TypeAliasSpec struct {
	Name           string
	UnderlyingType TypeReference
	Comment        string
}

// NewTypeAliasSpec returns a new spec representing a type alias.
func NewTypeAliasSpec(name string, typeRef TypeReference) *TypeAliasSpec {
	return &TypeAliasSpec{
		Name:           name,
		UnderlyingType: typeRef,
	}
}

// AliasComment adds a comment to a type alias.
func (a *TypeAliasSpec) AliasComment(comment string) *TypeAliasSpec {
	a.Comment = comment
	return a
}

// GetName returns the alias for this Type Alias.
func (a *TypeAliasSpec) GetName() string {
	return a.Name
}

// GetImports returns a slice of imports that the aliased type requires.
func (a *TypeAliasSpec) GetImports() []Import {
	return a.UnderlyingType.GetImports()
}

func (a *TypeAliasSpec) String() string {
	w := newCodeWriter()
	w.WriteCodeBlock(a)
	return w.String()
}

// GetStatements returns the comment and the type definition.
func (a *TypeAliasSpec) GetStatements() []Statement {
	var s []Statement
	s = append(s, Comment(a.Comment).GetStatements()...)
	s = append(s, newStatement(0, 0, "type $T $T", a, a.UnderlyingType))
	return s
}
