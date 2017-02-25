package poet

import (
	"bytes"
	"fmt"
	"io"
)

// FuncSpec represents information needed to write a function
type FuncSpec struct {
	Name             string
	Comment          string
	Parameters       []IdentifierParameter
	ResultParameters []IdentifierParameter
	Statements       []Statement
}

var _ CodeBlock = (*FuncSpec)(nil)

// NewFuncSpec returns a FuncSpec with the given name
func NewFuncSpec(name string) *FuncSpec {
	return &FuncSpec{
		Name:             name,
		Parameters:       []IdentifierParameter{},
		ResultParameters: []IdentifierParameter{},
		Statements:       []Statement{},
	}
}

// String returns a string representation of the function
func (f *FuncSpec) String() string {
	w := newCodeWriter()
	w.WriteCodeBlock(f)
	return w.String()
}

// Signature returns a format string and slice of arguments for the function's signature, not
// including the starting "func", any receiver, or opening curly brace
func (f *FuncSpec) Signature() (string, []interface{}) {
	// create a buffer for the format string and a slice for the arguments to the format string
	b := &bytes.Buffer{}
	var args []interface{}

	fmt.Fprint(b, "$L(")
	args = append(args, f.Name)
	args = append(args, writeParameters(b, f.Parameters)...)
	fmt.Fprint(b, ")")

	switch {
	case len(f.ResultParameters) == 0:
		break
	case len(f.ResultParameters) == 1 && f.ResultParameters[0].Name == "":
		// if there is only one parameter and the parameter is unnamed, do not wrap it in parens
		fmt.Fprint(b, " ")
		args = append(args, writeParameters(b, f.ResultParameters)...)
	default:
		fmt.Fprint(b, " (")
		args = append(args, writeParameters(b, f.ResultParameters)...)
		fmt.Fprint(b, ")")
	}

	return b.String(), args
}

// writeParamters writes a format to w and returns arguments for the format.
func writeParameters(w io.Writer, params []IdentifierParameter) []interface{} {
	var args []interface{}

	for i, p := range params {
		// if the argument is named, add its name to the format string
		if p.Name != "" {
			fmt.Fprint(w, "$L ")
			args = append(args, p.Name)
		}

		// add its type
		fmt.Fprint(w, "$T")
		args = append(args, p.Type)

		// if the argument is variadic, add the '...', will never happen for
		// result parameters
		if p.Variadic {
			fmt.Fprint(w, "...")
		}

		// if its not the last parameter, add a comma
		if i != len(params)-1 {
			fmt.Fprint(w, ", ")
		}
	}

	return args
}

// GetImports returns a slice of imports that this function needs, including
// parameters, result parameters, and statements within the function
func (f *FuncSpec) GetImports() []Import {
	packages := []Import{}

	for _, st := range f.Statements {
		for _, arg := range st.Arguments {
			if asTypeRef, ok := arg.(TypeReference); ok {
				packages = append(packages, asTypeRef.GetImports()...)
			}
		}
	}

	for _, param := range f.Parameters {
		packages = append(packages, param.Type.GetImports()...)
	}

	for _, param := range f.ResultParameters {
		packages = append(packages, param.Type.GetImports()...)
	}

	return packages
}

// GetStatements returns the Statements that make up the function.
func (f *FuncSpec) GetStatements() []Statement {
	signature, args := f.Signature()
	sigFormat := fmt.Sprintf("func %s {", signature)

	var s []Statement
	s = append(s, Comment(f.Comment).GetStatements()...)
	s = append(s, newStatement(0, 1, sigFormat, args...))
	s = append(s, f.Statements...)
	s = append(s, newStatement(-1, 0, "}"))
	return s
}

// Statement is a convenient method to append a statement to the function
func (f *FuncSpec) Statement(format string, args ...interface{}) *FuncSpec {
	f.Statements = append(f.Statements, newStatement(0, 0, format, args...))

	return f
}

// BlockStart is a convenient method to append a statement that marks the start of a
// block of code.
func (f *FuncSpec) BlockStart(format string, args ...interface{}) *FuncSpec {
	f.Statements = append(f.Statements, newStatement(0, 1, format+" {", args...))

	return f
}

// BlockEnd is a convenient method to append a statement that marks the end of a
// block of code.
func (f *FuncSpec) BlockEnd() *FuncSpec {
	f.Statements = append(f.Statements, newStatement(-1, 0, "}"))

	return f
}

// Parameter is a convenient method to append a parameter to the function
func (f *FuncSpec) Parameter(name string, spec TypeReference) *FuncSpec {
	f.Parameters = append(f.Parameters, IdentifierParameter{
		Identifier: Identifier{
			Name: name,
			Type: spec,
		},
	})

	return f
}

// VariadicParameter is a convenient method to append a parameter to the function
func (f *FuncSpec) VariadicParameter(name string, spec TypeReference) *FuncSpec {
	f.Parameters = append(f.Parameters, IdentifierParameter{
		Identifier: Identifier{
			Name: name,
			Type: spec,
		},
		Variadic: true,
	})

	return f
}

// ResultParameter is a convenient method to append a result parameter to the function
func (f *FuncSpec) ResultParameter(name string, spec TypeReference) *FuncSpec {
	f.ResultParameters = append(f.ResultParameters, IdentifierParameter{
		Identifier: Identifier{
			Name: name,
			Type: spec,
		},
	})

	return f
}

// FunctionComment adds a comment to the function
func (f *FuncSpec) FunctionComment(comment string) *FuncSpec {
	f.Comment = comment

	return f
}
