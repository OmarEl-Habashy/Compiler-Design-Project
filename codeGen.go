package main

import (
	"fmt"
	"strings"
)

type CodeGenerator struct {
	code          []string
	labelCounter  int
	symbols       map[string]string // variable -> type
	floatLiterals map[string]string // float value -> label (e.g., "9.9" -> "FLT_L1")
	dataSection   []string
	codeSection   []string
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		code:          []string{},
		labelCounter:  0,
		symbols:       make(map[string]string),
		floatLiterals: make(map[string]string),
		dataSection:   []string{},
		codeSection:   []string{},
	}
}

func (cg *CodeGenerator) nextLabel() string {
	cg.labelCounter++
	return fmt.Sprintf("L%d", cg.labelCounter)
}

func (cg *CodeGenerator) emit(instruction string) {
	cg.codeSection = append(cg.codeSection, "    "+instruction)
}

func (cg *CodeGenerator) emitLabel(label string) {
	cg.codeSection = append(cg.codeSection, label+":")
}

// declareFloatLiteral creates a unique label for a floating-point constant and adds it to the data section
func (cg *CodeGenerator) declareFloatLiteral(value string) string {
	// Check if we already declared this literal
	if label, ok := cg.floatLiterals[value]; ok {
		return label
	}

	// Create a new unique label for the constant
	label := fmt.Sprintf("FLT_L%d", len(cg.floatLiterals)+1)

	// Add declaration to data section (REAL8 = 64-bit double)
	cg.dataSection = append(cg.dataSection, fmt.Sprintf("    %s REAL8 %s", label, value))

	cg.floatLiterals[value] = label

	return label
}

func (cg *CodeGenerator) Generate(stmt Statement) string {
	// Generate code
	cg.genStmt(stmt)

	// Build final assembly
	var result []string

	// Header of file
	result = append(result, ".386")
	result = append(result, ".model flat, stdcall")
	result = append(result, "option casemap:none")
	result = append(result, "")

	// Data section with variables and float literals
	result = append(result, ".data")

	// Variable declarations
	for name, typ := range cg.symbols {
		if typ == "double" {
			result = append(result, fmt.Sprintf("    %s REAL8 0.0", name))
		} else {
			result = append(result, fmt.Sprintf("    %s DWORD 0", name))
		}
	}

	// Float literal declarations (FLT_L1, FLT_L2, etc.)
	result = append(result, cg.dataSection...)
	result = append(result, "")

	// Code section
	result = append(result, ".code")
	result = append(result, "start:")
	result = append(result, cg.codeSection...)
	result = append(result, "")
	result = append(result, "    ; Exit program")
	result = append(result, "    mov eax, 0")
	result = append(result, "    ret")
	result = append(result, "end start")

	return strings.Join(result, "\n")
}

func (cg *CodeGenerator) genStmt(stmt Statement) {
	switch s := stmt.(type) {
	case *BlockStatement:
		for _, st := range s.Statements {
			cg.genStmt(st)
		}

	case *AssignmentStatement:
		cg.genAssignment(s)

	case *IfStatement:
		cg.genIf(s)

	case *ExpressionStatement:
		cg.genExpr(s.Expression)
	}
}

func (cg *CodeGenerator) genAssignment(stmt *AssignmentStatement) {
	// Track variable type based on the value
	varType := "int"
	if v, ok := stmt.Value.(*NumberLiteral); ok {
		if strings.Contains(v.Value, ".") {
			varType = "double"
		}
	}
	cg.symbols[stmt.Name] = varType

	cg.emit(fmt.Sprintf("; Assignment: %s = ... (Type: %s)", stmt.Name, varType))

	cg.genExpr(stmt.Value)

	// Store result in variable based on type
	if varType == "double" {
		// Floating Point: Store from FPU stack (ST(0)) to memory
		cg.emit(fmt.Sprintf("fstp [%s]  ; store float and pop FPU stack", stmt.Name))
	} else {
		// Integer: Store from EAX to memory
		cg.emit(fmt.Sprintf("mov [%s], eax", stmt.Name))
	}
}

func (cg *CodeGenerator) genIf(stmt *IfStatement) {
	elseLabel := cg.nextLabel()
	endLabel := cg.nextLabel()

	cg.emit("; If statement begins")

	// Generate condition
	cg.genCondition(stmt.Condition, elseLabel)

	// Generate body (true branch)
	cg.emit("; Then branch")
	cg.genStmt(stmt.Body)
	cg.emit(fmt.Sprintf("jmp %s", endLabel))

	// Else label (false branch - empty)
	cg.emitLabel(elseLabel)
	cg.emit("; Else branch (empty)")

	// End label
	cg.emitLabel(endLabel)
	cg.emit("; End if statement")
}

func (cg *CodeGenerator) genCondition(cond Expression, falseLabel string) {
	if c, ok := cond.(*Condition); ok {
		cg.emit("; Compare operands (integer comparison)")

		// Load left operand
		cg.genExpr(c.Left)
		cg.emit("push eax")

		// Load right operand
		cg.genExpr(c.Right)
		cg.emit("mov ebx, eax")
		cg.emit("pop eax")

		// Compare
		cg.emit("cmp eax, ebx")

		// Jump based on operator -- if is false
		switch c.Operator {
		case "==":
			cg.emit(fmt.Sprintf("jne %s  ; jump if not equal", falseLabel))
		case "!=":
			cg.emit(fmt.Sprintf("je %s   ; jump if equal", falseLabel))
		case ">":
			cg.emit(fmt.Sprintf("jle %s  ; jump if less or equal", falseLabel))
		case "<":
			cg.emit(fmt.Sprintf("jge %s  ; jump if greater or equal", falseLabel))
		case ">=":
			cg.emit(fmt.Sprintf("jl %s   ; jump if less", falseLabel))
		case "<=":
			cg.emit(fmt.Sprintf("jg %s   ; jump if greater", falseLabel))
		}
	}
}

func (cg *CodeGenerator) genExpr(expr Expression) {
	switch e := expr.(type) {
	case *NumberLiteral:
		if strings.Contains(e.Value, ".") {
			// Floating Point Literal: Use Floating point unit
			label := cg.declareFloatLiteral(e.Value)
			cg.emit(fmt.Sprintf("; Load float %s onto FPU stack", e.Value))
			cg.emit(fmt.Sprintf("fld %s  ; load float from data section", label))
		} else {
			// Integer Literal: Use EAX
			cg.emit(fmt.Sprintf("mov eax, %s", e.Value))
		}

	case *Identifier:
		// Load variable value into EAX
		// Note: If the variable is float type, proper handling would require
		// checking the symbol table and using 'fld' instead
		cg.emit(fmt.Sprintf("mov eax, [%s]", e.Value))

	case *Condition:
		// Standalone condition (shouldn't happen in normal flow)
		cg.emit("; condition expression")
	}
}

func (cg *CodeGenerator) GetVariableDeclarations() string {
	// This method is no longer needed as we handle it in Generate()
	return ""
}
