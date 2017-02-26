package poet

import (
	"bytes"
	"strings"
)

// codeWriter keeps track of the current indentation and writes code to a buffer
type codeWriter struct {
	buffer        bytes.Buffer
	currentIndent int
	importPath    string // the go import package this writer is writing to
}

// newCodeWriter constructs a new codeWriter writing to pkg.
func newCodeWriter() *codeWriter {
	return &codeWriter{
		buffer: bytes.Buffer{},
	}
}

// newCodeWriterWithImport constructs a new scodeWriter writing to pkg.
func newCodeWriterWithImport(importPath string) *codeWriter {
	return &codeWriter{
		buffer:     bytes.Buffer{},
		importPath: importPath,
	}
}

// WriteCode writes code at the given indentation
func (c *codeWriter) WriteCode(code string) {
	c.buffer.WriteString(strings.Repeat("\t", c.currentIndent))
	c.buffer.WriteString(code)
}

// WriteCodeBlock writes a code block at the given indentation
func (c *codeWriter) WriteCodeBlock(block CodeBlock) {
	for _, s := range block.GetStatements() {
		c.WriteStatement(s)
	}
}

// WriteStatement writes a new line of code with the current indentation and augments
// the indentation per the statement. A newline is appended at the end of the statement.
func (c *codeWriter) WriteStatement(s Statement) {
	c.currentIndent += s.BeforeIndent
	c.WriteCode(template(s.Format, s.Arguments...) + "\n")
	c.currentIndent += s.AfterIndent
}

// String gives a string with the code
func (c *codeWriter) String() string {
	return c.buffer.String()
}

func newStatement(beforeIndent, afterIndent int, format string, args ...interface{}) Statement {
	return Statement{
		BeforeIndent: beforeIndent,
		AfterIndent:  afterIndent,
		Format:       format,
		Arguments:    args,
	}
}

// Appends two Statements without a newline.
func appendStatements(first, second Statement) Statement {
	args := append(first.Arguments, second.Arguments...)
	return Statement{
		BeforeIndent: first.BeforeIndent,
		AfterIndent:  second.AfterIndent,
		Format:       first.Format + second.Format,
		Arguments:    args,
	}
}
