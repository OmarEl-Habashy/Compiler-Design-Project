# Compiler Design Project

A simple compiler implementation in Go that demonstrates the fundamental phases of compilation: lexical analysis, syntax analysis (parsing), and semantic analysis.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Components](#components)
- [Installation & Usage](#installation--usage)
- [Example](#example)
- [Supported Language Features](#supported-language-features)
- [Error Handling](#error-handling)

## 🎯 Overview

This compiler processes a simple C-like language and generates:
1. **Tokens** - Lexical analysis output
2. **Abstract Syntax Tree (AST)** - Parse tree in JSON format
3. **Semantic Tree** - Type-annotated AST with semantic analysis

## ✨ Features

- **Lexical Analysis**: Tokenizes input code into keywords, identifiers, constants, operators, and special symbols
- **Syntax Analysis**: Builds an Abstract Syntax Tree (AST) from tokens
- **Semantic Analysis**: Performs type checking and inference
- **Error Detection**: Comprehensive error messages for syntax and semantic errors
- **JSON Output**: Pretty-printed JSON representation of parse and semantic trees

## 📁 Project Structure

```
Compiler/
├── main.go          # Entry point and orchestration
├── lexer.go         # Lexical analyzer (tokenizer)
├── parser.go        # Syntax analyzer (parser)
├── semantic.go      # Semantic analyzer (type checker)
├── token.go         # Token type definitions
├── tree.go          # AST node definitions and JSON serialization
├── test_code.txt    # Sample input code
└── go.mod           # Go module file
```

## 🔧 Components

### 1. Lexer (`lexer.go`)

**Purpose**: Converts raw source code into a sequence of tokens.

**Token Types**:
- `KEYWORD`: Reserved words (if, else, while, for, int, double, return)
- `IDENTIFIER`: Variable names (x, y, myVar)
- `CONSTANT`: Numeric literals (5, 10, 3.14)
- `OPERATOR`: Comparison and assignment operators (==, !=, >, <, >=, <=, =)
- `SPECIAL_SYMBOL`: Punctuation and delimiters ((, ), {, }, ;)

**Key Features**:
- Regex-based pattern matching
- Automatic whitespace skipping
- Error reporting for unrecognized characters

### 2. Parser (`parser.go`)

**Purpose**: Builds an Abstract Syntax Tree (AST) from the token stream.

**Supported Constructs**:
- **If Statements**: `if (condition) { ... }`
- **Block Statements**: `{ stmt1; stmt2; ... }`
- **Assignment Statements**: `x = 5;`
- **Expression Statements**: `x;`
- **Conditions**: `x > 5`, `y != 10`, etc.

**Key Features**:
- Recursive descent parsing
- Comprehensive syntax error detection
- Special error messages for common mistakes (using `=` instead of `==` in conditions)

**Parser Methods**:
- `ParseIfStatement()`: Parses if statements
- `parseCondition()`: Parses boolean conditions
- `parseBlockStatement()`: Parses code blocks
- `parseAssignmentStatement()`: Parses variable assignments
- `parseExpressionStatement()`: Parses expression statements
- `parsePrimary()`: Parses identifiers and literals

### 3. Semantic Analyzer (`semantic.go`)

**Purpose**: Performs type checking and inference on the AST.

**Features**:
- **Type Inference**: Automatically infers types (int, double) from literals
- **Symbol Table**: Tracks variable types across the program
- **Type Checking**: Ensures type consistency in assignments and operations
- **Condition Validation**: Verifies that if-conditions evaluate to boolean

**Type System**:
- `int`: Integer literals (e.g., 5, 10)
- `double`: Floating-point literals (e.g., 3.14)
- `bool`: Result of comparison operations
- `unknown`: Unresolved or error types

**Error Detection**:
- Type mismatches in assignments
- Non-numeric operands in comparisons
- Variables used before assignment (warning)
- Non-boolean if conditions

### 4. AST Structure (`tree.go`)

**Node Types**:

**Statements**:
- `IfStatement`: If-then control flow
- `BlockStatement`: Group of statements
- `AssignmentStatement`: Variable assignment
- `ExpressionStatement`: Standalone expression

**Expressions**:
- `Identifier`: Variable reference
- `NumberLiteral`: Numeric constant
- `Condition`: Binary comparison operation

**JSON Serialization**:
- `exprToJSON()`: Converts expressions to JSON
- `stmtToJSON()`: Converts statements to JSON
- `GetParseTreeJSON()`: Generates formatted JSON output

## 🚀 Installation & Usage

### Prerequisites
- Go 1.16 or higher

### Running the Compiler

1. **Clone or navigate to the project directory**:
   ```bash
   cd "d:\Semester 5\Systems Prog\Compiler"
   ```

2. **Create or edit the input file** (`test_code.txt`):
   ```
   if (x != 5) { 
       y = 10;
       z = 9; 
   }
   ```

3. **Run the compiler**:
   ```bash
   go run .
   ```

### Output

The compiler produces three main outputs:

1. **Tokens**: List of lexical tokens
2. **Parse Tree**: JSON representation of the AST
3. **Semantic Tree**: Type-annotated AST with inferred types

## 📝 Example

### Input Code (`test_code.txt`)
```
if (x != 5) { 
    y = 10;
    z = 9; 
}
```

### Output

**Tokens**:
```
Token(KEYWORD, 'if')
Token(SPECIAL_SYMBOL, '(')
Token(IDENTIFIER, 'x')
Token(OPERATOR, '!=')
Token(CONSTANT, '5')
Token(SPECIAL_SYMBOL, ')')
Token(SPECIAL_SYMBOL, '{')
Token(IDENTIFIER, 'y')
Token(OPERATOR, '=')
Token(CONSTANT, '10')
Token(SPECIAL_SYMBOL, ';')
...
```

**Parse Tree (AST)**:
```json
{
  "type": "IfStatement",
  "children": [
    {
      "type": "Condition",
      "operator": "!=",
      "children": [
        {
          "type": "Identifier",
          "value": "x"
        },
        {
          "type": "NumberLiteral",
          "value": "5"
        }
      ]
    },
    {
      "type": "BlockStatement",
      "children": [
        {
          "type": "AssignmentStatement",
          "value": "y",
          "children": [
            {
              "type": "NumberLiteral",
              "value": "10"
            }
          ]
        },
        {
          "type": "AssignmentStatement",
          "value": "z",
          "children": [
            {
              "type": "NumberLiteral",
              "value": "9"
            }
          ]
        }
      ]
    }
  ]
}
```

**Semantic Tree**:
```json
{
  "type": "IfStatement",
  "inferredType": "unknown",
  "children": [
    {
      "type": "Condition",
      "operator": "!=",
      "inferredType": "bool",
      "children": [
        {
          "type": "Identifier",
          "value": "x",
          "inferredType": "int"
        },
        {
          "type": "NumberLiteral",
          "value": "5",
          "inferredType": "int"
        }
      ]
    },
    {
      "type": "BlockStatement",
      "inferredType": "unknown",
      "children": [...]
    }
  ]
}
```

## 🔤 Supported Language Features

### Keywords
- `if`, `else`, `while`, `for`
- `int`, `double`, `return`

### Operators
- **Comparison**: `==`, `!=`, `>`, `<`, `>=`, `<=`
- **Assignment**: `=`

### Data Types
- **Integer**: Whole numbers (e.g., `5`, `100`)
- **Double**: Floating-point numbers (e.g., `3.14`, `2.5`)

### Syntax Rules
- Statements must end with semicolons (`;`)
- Code blocks must be enclosed in braces (`{` and `}`)
- If-conditions must be enclosed in parentheses (`(` and `)`)
- Assignments use single `=`, comparisons use `==`

## ⚠️ Error Handling

### Lexer Errors
- **Unknown characters**: Reports character and position
  ```
  Lexer Error: Unknown character @ at position 15
  ```

### Parser Errors
- **Missing tokens**: Expected vs. actual token
  ```
  expected SPECIAL_SYMBOL, got IDENTIFIER at position 5
  ```

- **Assignment in condition**: Helpful suggestion
  ```
  syntax error: cannot use assignment '=' in if condition, did you mean '=='?
  ```

- **Comparison as statement**: Clear error message
  ```
  syntax error: 'x >' - comparison cannot be used as a statement (did you mean 'x = ...'?)
  ```

- **Missing semicolons**: Specific location
  ```
  expected ';' after assignment 'y = ...'
  ```

### Semantic Errors
- **Type mismatch**: Variable type conflicts
  ```
  semantic error: type mismatch for 'x' (was int, now double)
  ```

- **Non-numeric operands**: Invalid operation types
  ```
  semantic error: non-numeric operands for '>'
  ```

- **Non-boolean condition**: Invalid if-condition type
  ```
  semantic error: IF condition must be bool
  ```

### Semantic Warnings
- **Undefined variables**: Used before assignment
  ```
  semantic warning: 'x' used before assignment; assuming int
  ```

## 🛠️ Extending the Compiler

### Adding New Token Types
1. Add constant to `token.go`
2. Add regex pattern to `lexer.go` rules
3. Update parser methods as needed

### Adding New Statement Types
1. Define struct in `tree.go` implementing `Statement` interface
2. Add parsing method in `parser.go`
3. Add case in `semStmt()` in `semantic.go`
4. Add case in `stmtToJSON()` in `tree.go`

### Adding New Operators
1. Update operator regex in `lexer.go`
2. Add to valid operators map in relevant parser methods
3. Update semantic analysis if needed

## 📚 Technical Details

### Parser Strategy
- **Type**: Recursive Descent Parser
- **Lookahead**: 1-2 tokens
- **Error Recovery**: Returns nil on error with error message

### AST Design
- **Interface-based**: `Statement` and `Expression` interfaces
- **Type-safe**: Go's type system ensures correctness
- **Serializable**: JSON export for visualization

### Type System
- **Static inference**: Types inferred at compile-time
- **Simple types**: int, double, bool
- **Symbol table**: Tracks variable types across scope

## 👥 Contributors

- Omar El-Habashy
- Nariman Adel
- Ahmed Morsy

## 📄 License

This project is for educational purposes as part of a Systems Programming course.

---

**Note**: This is a simplified compiler for educational purposes. It demonstrates core compiler concepts but is not production-ready.
