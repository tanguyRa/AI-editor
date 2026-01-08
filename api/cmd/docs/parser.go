package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// parseGoFile parses a Go file and extracts types and functions
func parseGoFile(filePath string, fileNode *FileNode) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			typeInfo := TypeInfo{
				Name:       x.Name.Name,
				Properties: make([]PropertyInfo, 0),
			}

			// Extract struct fields
			switch t := x.Type.(type) {
			case *ast.StructType:
				if t.Fields != nil {
					for _, field := range t.Fields.List {
						if len(field.Names) == 0 {
							// Handle embedded/anonymous fields
							prop := PropertyInfo{
								Name: fmt.Sprintf("*%s", exprToString(field.Type)), // Mark as embedded
								Type: "embedded",
							}
							typeInfo.Properties = append(typeInfo.Properties, prop)
						} else {
							// Handle named fields
							for _, name := range field.Names {
								prop := PropertyInfo{
									Name: name.Name,
									Type: exprToString(field.Type),
								}
								typeInfo.Properties = append(typeInfo.Properties, prop)
							}
						}
					}
				}
			default:
				// For non-struct types, store the underlying type as a special property
				typeInfo.Properties = append(typeInfo.Properties, PropertyInfo{
					Name: "__underlying", // Special marker for underlying type
					Type: exprToString(x.Type),
				})
			}

			fileNode.Types = append(fileNode.Types, typeInfo)
		case *ast.FuncDecl:
			funcSig := buildFunctionSignature(x)
			fileNode.Funcs = append(fileNode.Funcs, funcSig)
		}
		return true
	})

	return nil
}

// buildFunctionSignature creates a function signature string from an AST function declaration
func buildFunctionSignature(x *ast.FuncDecl) string {
	funcSig := "func "
	if x.Recv != nil {
		recv := x.Recv.List[0]
		recvType := exprToString(recv.Type)
		funcSig += fmt.Sprintf("(%s) ", recvType)
	}
	funcSig += x.Name.Name
	funcSig += "("
	
	// Add parameters
	if x.Type.Params != nil {
		params := make([]string, 0)
		for _, p := range x.Type.Params.List {
			paramType := exprToString(p.Type)
			if len(p.Names) == 0 {
				params = append(params, paramType)
			} else {
				for _, name := range p.Names {
					params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
				}
			}
		}
		funcSig += strings.Join(params, ", ")
	}
	funcSig += ")"

	// Add return types
	if x.Type.Results != nil {
		returns := make([]string, 0)
		for _, r := range x.Type.Results.List {
			returnType := exprToString(r.Type)
			if len(r.Names) == 0 {
				returns = append(returns, returnType)
			} else {
				for _, name := range r.Names {
					returns = append(returns, fmt.Sprintf("%s %s", name.Name, returnType))
				}
			}
		}
		if len(returns) == 1 && len(x.Type.Results.List[0].Names) == 0 {
			funcSig += " " + returns[0]
		} else {
			funcSig += " (" + strings.Join(returns, ", ") + ")"
		}
	}
	
	return funcSig
}