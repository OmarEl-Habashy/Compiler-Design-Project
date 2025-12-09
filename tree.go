package main

import (
    "encoding/json"
    "fmt"
    "strings" 
)



type ASTNode interface{ node() }

type Statement interface {
    ASTNode
    statementNode()
}

type Expression interface {
    ASTNode
    expressionNode()
}


type Program struct {
    Statements []Statement
}

func (p *Program) node() {}


type IfStatement struct {
    Condition   Expression
    Body Statement
}

func (is *IfStatement) node()          {}
func (is *IfStatement) statementNode() {}

type BlockStatement struct {
    Statements []Statement
}

func (bs *BlockStatement) node()          {}
func (bs *BlockStatement) statementNode() {}

// x+5
type ExpressionStatement struct {
    Expression Expression
}

func (es *ExpressionStatement) node()          {}
func (es *ExpressionStatement) statementNode() {}

// x=5
type AssignmentStatement struct {
    Name  string     //x  
    Value Expression  //5
}

func (as *AssignmentStatement) node()          {}
func (as *AssignmentStatement) statementNode() {}



// Variable name 
type Identifier struct {
    Value string
}

func (i *Identifier) node()           {}
func (i *Identifier) expressionNode() {}

// Number
type NumberLiteral struct {
    Value string
}

func (nl *NumberLiteral) node()           {}
func (nl *NumberLiteral) expressionNode() {}

type Condition struct {
    Left     Expression //x
    Operator string // ==
    Right    Expression //5
}

func (c *Condition) node()           {}
func (c *Condition) expressionNode() {}


type ParseTreeNode struct {
    Type     string           `json:"type"`
    Value    string           `json:"value,omitempty"`
    Operator string           `json:"operator,omitempty"`
    Children []*ParseTreeNode `json:"children,omitempty"`
}



func exprToJSON(expr Expression) *ParseTreeNode {
    switch node := expr.(type) {
    case *Identifier:
        return &ParseTreeNode{
            Type:  "Identifier",
            Value: node.Value,
        }
    case *NumberLiteral:
        return &ParseTreeNode{
            Type:  "NumberLiteral",
            Value: node.Value,
        }
    case *Condition:
        return &ParseTreeNode{
            Type:     "Condition",
            Operator: node.Operator,
            Children: []*ParseTreeNode{
                exprToJSON(node.Left),
                exprToJSON(node.Right),
            },
        }
    default:
        return &ParseTreeNode{Type: "UnknownExpression"}
    }
}

func stmtToJSON(stmt Statement) *ParseTreeNode {
    switch node := stmt.(type) {
    case *IfStatement:
        return &ParseTreeNode{
            Type: "IfStatement",
            Children: []*ParseTreeNode{
                exprToJSON(node.Condition),      
                stmtToJSON(node.Body),    
            },
        }
    case *BlockStatement:
        children := make([]*ParseTreeNode, 0, len(node.Statements))
        for _, s := range node.Statements {
            children = append(children, stmtToJSON(s))
        }
        return &ParseTreeNode{
            Type:     "BlockStatement",
            Children: children,
        }
    case *ExpressionStatement:
        return &ParseTreeNode{
            Type: "ExpressionStatement",
            Children: []*ParseTreeNode{
                exprToJSON(node.Expression),
            },
        }
    case *AssignmentStatement:
        return &ParseTreeNode{
            Type:  "AssignmentStatement",
            Value: node.Name,
            Children: []*ParseTreeNode{
                exprToJSON(node.Value),
            },
        }
    default:
        return &ParseTreeNode{Type: "UnknownStatement"}
    }
}

func GetParseTreeJSON(stmt Statement) string {
    tree := stmtToJSON(stmt)
    b, err := json.MarshalIndent(tree, "", "  ")
    if err != nil {
        return fmt.Sprintf(`{"error": "%s"}`, err.Error())
    }
    
    output := string(b)
    output = strings.ReplaceAll(output, "\\u003e", ">")
    output = strings.ReplaceAll(output, "\\u003c", "<")
    output = strings.ReplaceAll(output, "\\u0026", "&")
    
    return output
}
