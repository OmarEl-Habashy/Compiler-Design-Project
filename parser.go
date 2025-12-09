package main

import (
    "fmt"
)

type Parser struct {
    tokens []Token
    pos    int
    errors []string
}

func NewParser(tokens []Token) *Parser {
    return &Parser{
        tokens: tokens,
        pos:    0,
        errors: []string{},
    }
}

func (p *Parser) currentToken() Token {
    if p.pos < len(p.tokens) {
        return p.tokens[p.pos]
    }
    return Token{Type: "EOF", Value: ""}
}

func (p *Parser) next() {
    p.pos++
}

func (p *Parser) expect(tokenType TokenType) bool {
    if p.currentToken().Type == tokenType {
        p.next()
        return true
    }
    p.errors = append(p.errors,
        fmt.Sprintf("expected %v, got %v at position %d",
            tokenType, p.currentToken().Type, p.pos))
    return false
}

func (p *Parser) Errors() []string {
    return p.errors
}



// Parse: if (x > 5) { ... }
func (p *Parser) ParseIfStatement() *IfStatement {
    stmt := &IfStatement{}

    // Expect "if"
    if !p.expect(KEYWORD) || p.tokens[p.pos-1].Value != "if" {
        return nil
    }

    // Expect "("
    if !p.expect(SPECIAL_SYMBOL) || p.tokens[p.pos-1].Value != "(" {
        return nil
    }

    // Parse condition: x > 5
    stmt.Condition = p.parseCondition()
    if stmt.Condition == nil {
        return nil
    }

    // Expect ")"
    if !p.expect(SPECIAL_SYMBOL) || p.tokens[p.pos-1].Value != ")" {
        return nil
    }

    // Parse body:
    if p.currentToken().Value == "{" {
        stmt.Body = p.parseBlockStatement()
        if stmt.Body == nil {
            return nil
        }
    } else {
        p.errors = append(p.errors, "expected '{' after if condition")
        return nil
    }

    return stmt
}

// Parse condition: x > 5
func (p *Parser) parseCondition() Expression {
    left := p.parsePrimary()
    if left == nil {
        return nil
    }

   
    if p.currentToken().Type == OPERATOR {
        op := p.currentToken().Value
        
        
        if op == "=" {
            p.errors = append(p.errors, 
                "syntax error: cannot use assignment '=' in if condition, did you mean '=='?")
            return nil
        }
        
       
        validOps := map[string]bool{
            "==": true, "!=": true, ">": true, 
            "<": true, ">=": true, "<=": true,
        }
        
        if !validOps[op] {
            p.errors = append(p.errors, 
                fmt.Sprintf("invalid operator '%s' in condition", op))
            return nil
        }
        
        p.next()

        right := p.parsePrimary()
        if right == nil {
            return nil
        }

        return &Condition{
            Left:     left,
            Operator: op,
            Right:    right,
        }
    }

    return left
}

// Parse primary: identifier or number
func (p *Parser) parsePrimary() Expression {
    tok := p.currentToken()

    switch tok.Type {
    case IDENTIFIER:
        p.next()
        return &Identifier{Value: tok.Value}
    case CONSTANT:
        p.next()
        return &NumberLiteral{Value: tok.Value}
    default:
        p.errors = append(p.errors,
            fmt.Sprintf("unexpected token: %v (%s)", tok.Type, tok.Value))
        return nil
    }
}

// Parse block: { stmt1; stmt2; ... }
func (p *Parser) parseBlockStatement() *BlockStatement {
    block := &BlockStatement{Statements: []Statement{}}

    if !p.expect(SPECIAL_SYMBOL) || p.tokens[p.pos-1].Value != "{" {
        return nil
    }

   
    for p.currentToken().Type != "EOF" && p.currentToken().Value != "}" {
        tok := p.currentToken()

      
        if tok.Type == IDENTIFIER && p.pos+1 < len(p.tokens) {
            nextTok := p.tokens[p.pos+1]
            
            
            comparisonOps := map[string]bool{
                "==": true, "!=": true, ">": true, 
                "<": true, ">=": true, "<=": true,
            }
            
            if comparisonOps[nextTok.Value] {
                p.errors = append(p.errors, 
                    fmt.Sprintf("syntax error: '%s %s' - comparison cannot be used as a statement (did you mean '%s = ...'?)",
                        tok.Value, nextTok.Value, tok.Value))
                return nil
            }
            
          
            if nextTok.Value == "=" {
                stmt := p.parseAssignmentStatement()
                if stmt == nil {
                    return nil
                }
                block.Statements = append(block.Statements, stmt)
                continue
            }
        }

        stmt := p.parseExpressionStatement()
        if stmt == nil {
            return nil
        }
        block.Statements = append(block.Statements, stmt)
    }

 
    if !p.expect(SPECIAL_SYMBOL) || p.tokens[p.pos-1].Value != "}" {
        return nil
    }

    return block
}

// Parse expression statement: x;
func (p *Parser) parseExpressionStatement() *ExpressionStatement {
    expr := p.parsePrimary()
    if expr == nil {
        return nil
    }

    if p.currentToken().Value != ";" {
        p.errors = append(p.errors, 
            fmt.Sprintf("expected ';' after expression"))
        return nil
    }
    p.next()

    return &ExpressionStatement{Expression: expr}
}

// Parse assignment: y = 10;
func (p *Parser) parseAssignmentStatement() *AssignmentStatement {
    stmt := &AssignmentStatement{}

 
    stmt.Name = p.currentToken().Value
    p.next()

    if p.currentToken().Value != "=" {
        p.errors = append(p.errors, "expected '=' in assignment")
        return nil
    }
    p.next()
    stmt.Value = p.parsePrimary()
    if stmt.Value == nil {
        return nil
    }
    if p.currentToken().Value != ";" {
        p.errors = append(p.errors, 
            fmt.Sprintf("expected ';' after assignment '%s = ...'", stmt.Name))
        return nil
    }
    p.next()

    return stmt
}