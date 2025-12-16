package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SemanticAnalyzer struct {
	// variable -> type ("int" / "double")
	symbols map[string]string
	// errors/warnings
	msg []string
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		symbols: map[string]string{},
		msg:     []string{},
	}
}

func (s *SemanticAnalyzer) Errors() []string { return s.msg }

// SemanticNode = same tree shape but with inferredType added
type SemanticNode struct {
	Type         string          `json:"type"`
	Value        string          `json:"value,omitempty"`
	Operator     string          `json:"operator,omitempty"`
	InferredType string          `json:"inferredType"`
	Children     []*SemanticNode `json:"children,omitempty"`
}

func (s *SemanticAnalyzer) Analyze(stmt Statement) *SemanticNode {
	return s.semStmt(stmt)
}

func (s *SemanticAnalyzer) semStmt(st Statement) *SemanticNode {
	switch n := st.(type) {

	case *IfStatement:
		cond := s.semExpr(n.Condition)
		if cond.InferredType != "bool" && cond.InferredType != "unknown" {
			s.msg = append(s.msg, "semantic error: IF condition must be bool")
		}
		body := s.semStmt(n.Body)
		return &SemanticNode{
			Type:         "IfStatement",
			InferredType: "unknown",
			Children:     []*SemanticNode{cond, body},
		}

	case *BlockStatement:
		out := &SemanticNode{Type: "BlockStatement", InferredType: "unknown"}
		for _, st2 := range n.Statements {
			out.Children = append(out.Children, s.semStmt(st2))
		}
		return out

	case *AssignmentStatement:
		rhs := s.semExpr(n.Value)

		// define variable on first assignment
		if rhs.InferredType != "unknown" {
			if old, ok := s.symbols[n.Name]; ok {
				if old != rhs.InferredType {
					// ERROR: Don't allow type changes
					s.msg = append(s.msg,
						fmt.Sprintf("semantic error: cannot change type of '%s' from %s to %s",
							n.Name, old, rhs.InferredType))
				}
			} else {
				s.symbols[n.Name] = rhs.InferredType
			}
		}

		return &SemanticNode{
			Type:         "AssignmentStatement",
			Value:        n.Name,
			InferredType: "unknown",
			Children:     []*SemanticNode{rhs},
		}

	case *ExpressionStatement:
		ex := s.semExpr(n.Expression)
		return &SemanticNode{
			Type:         "ExpressionStatement",
			InferredType: "unknown",
			Children:     []*SemanticNode{ex},
		}

	default:
		s.msg = append(s.msg, "semantic error: unsupported statement type")
		return &SemanticNode{Type: "UnknownStatement", InferredType: "unknown"}
	}
}

func (s *SemanticAnalyzer) semExpr(e Expression) *SemanticNode {
	switch n := e.(type) {

	case *Identifier:
		t, ok := s.symbols[n.Value]
		if !ok {
			// warning (keeps your demo working even if x not assigned)
			s.msg = append(s.msg, fmt.Sprintf("semantic warning: '%s' used before assignment; assuming int", n.Value))
			s.symbols[n.Value] = "int"
			t = "int"
		}
		return &SemanticNode{Type: "Identifier", Value: n.Value, InferredType: t}

	case *NumberLiteral:
		t := "int"
		if strings.Contains(n.Value, ".") {
			t = "double"
		}
		return &SemanticNode{Type: "NumberLiteral", Value: n.Value, InferredType: t}

	case *Condition:
		left := s.semExpr(n.Left)
		right := s.semExpr(n.Right)

		// require numeric operands
		if !isNumeric(left.InferredType) || !isNumeric(right.InferredType) {
			s.msg = append(s.msg, fmt.Sprintf("semantic error: non-numeric operands for '%s'", n.Operator))
		}

		return &SemanticNode{
			Type:         "Condition",
			Operator:     n.Operator,
			InferredType: "bool",
			Children:     []*SemanticNode{left, right},
		}

	default:
		s.msg = append(s.msg, "semantic error: unsupported expression type")
		return &SemanticNode{Type: "UnknownExpression", InferredType: "unknown"}
	}
}

func isNumeric(t string) bool {
	return t == "int" || t == "double"
}

func GetSemanticTreeJSON(root *SemanticNode) string {
	b, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}

	// Keep operators readable
	out := string(b)
	out = strings.ReplaceAll(out, "\\u003e", ">")
	out = strings.ReplaceAll(out, "\\u003c", "<")
	out = strings.ReplaceAll(out, "\\u0026", "&")
	return out
}
