package code

import (
	"bytes"
	"fmt"
	"monkey-i/ast"
	"strings"
)

type NamedFunctionBlock struct {
	ast.FunctionBlock
	Name string
}

func (nfl *NamedFunctionBlock) expressionNode()      {}
func (nfl *NamedFunctionBlock) TokenLiteral() string { return nfl.Token.Literal }
func (nfl *NamedFunctionBlock) String() string {
	var out bytes.Buffer

	if nfl.Name != "" {
		out.WriteString(fmt.Sprintf("<%s>", nfl.Name))
	}

	params := []string{}
	for _, param := range nfl.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(nfl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(nfl.Body.String())

	return out.String()
}
