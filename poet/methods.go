package poet

import (
	"fmt"
)

// MethodSpec represents a method, with a receiver name and type.
type MethodSpec struct {
	FuncSpec
	ReceiverName string
	Receiver     TypeReference
}

var _ CodeBlock = (*MethodSpec)(nil)

// NewMethodSpec creates a new method with the given method name, receiverName, and receiver type.
func NewMethodSpec(name, receiverName string, receiver TypeReference) *MethodSpec {
	return &MethodSpec{
		FuncSpec: FuncSpec{
			Name: name,
		},
		ReceiverName: receiverName,
		Receiver:     receiver,
	}
}

func (m *MethodSpec) GetStatements() []Statement {
	signature, args := m.FuncSpec.Signature()

	// add method receiver
	signature = fmt.Sprintf("func ($L $T) %s {", signature)
	args = append([]interface{}{m.ReceiverName, m.Receiver}, args...)

	var s []Statement
	s = append(s, Comment(m.Comment).GetStatements()...)
	s = append(s, newStatement(0, 1, signature, args...))
	s = append(s, m.Statements...)
	s = append(s, newStatement(-1, 0, "}"))
	return s
}

// String returns a string representation of the function
func (m *MethodSpec) String() string {
	w := newCodeWriter()
	w.WriteCodeBlock(m)
	return w.String()
}
