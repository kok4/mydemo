package parsego

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

/*
\brief: 分析Go文件的语法
\参考:
	https://medium.com/justforfunc/understanding-go-programs-with-go-parser-c4e88a6edb87
*/

type FunInfo struct {
	Lbrace int
	Rbrace int
}

type GoParser struct {
	fset *token.FileSet
	f    *ast.File
}

func NewParseGo(filePath string) *GoParser {
	obj := &GoParser{
		fset: token.NewFileSet(), // positions are relative to fset
	}
	var err error
	obj.f, err = parser.ParseFile(obj.fset, filePath, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	return obj
}

func NewPaseGoByStr(content string) *GoParser {
	obj := &GoParser{
		fset: token.NewFileSet(), // positions are relative to fset
	}

	var err error
	obj.f, err = parser.ParseFile(obj.fset, "", content, parser.AllErrors)
	if err != nil {
		fmt.Println("Err:", err)
	}
	return obj
}

func (pg *GoParser) GetFunctions(f func(retKey, retVal interface{}) bool) error {
	if pg.f == nil {
		return fmt.Errorf("PrintAST; ast.File is nil")
	}
	v := NewVisitorFunc(pg.fset)
	ast.Walk(v, pg.f)
	v.Range(func(key, value interface{}) bool {
		obj := &FunInfo{}
		ret := value.(*bodyInfo)
		obj.Lbrace = pg.fset.Position(ret.Lbrace).Offset
		obj.Rbrace = pg.fset.Position(ret.Rbrace).Offset

		return f(key, obj)
	})
	return nil
}

func (pg *GoParser) PrintAST() {
	if pg.f == nil {
		panic(fmt.Errorf("PrintAST; ast.File is nil"))
	}
	var v visitorAST
	ast.Walk(v, pg.f)
}

/*
	visitorAST 用于打印AST(语法树)
*/
type visitorAST int

func (v visitorAST) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	fmt.Printf("%s%T\n", strings.Repeat("\t", int(v)), n)
	return v + 1
}

/*
	visitorFun 用于获取函数
	Note:这里简化操作，只遍历函数名和函数体，进行互换操作.
*/
type bodyInfo struct {
	Lbrace token.Pos
	Rbrace token.Pos
}

type visitorFun struct {
	locals map[string]*bodyInfo
}

func NewVisitorFunc(fset *token.FileSet) *visitorFun {
	obj := &visitorFun{
		locals: make(map[string]*bodyInfo),
	}
	return obj
}

func (v visitorFun) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.FuncDecl:
		ident := d.Name // Function's name
		if ident == nil {
			fmt.Println("VisitorFun::Visit d.Name fail")
			return v
		}

		funcType := d.Type //Function's type
		if funcType == nil {
			fmt.Println("VisitorFun::Visit d.type fail")
			return v
		}

		blockStmt := d.Body // Function's body
		if blockStmt == nil {
			fmt.Println("VisitorFun::Visit d.Body fail")
			return v
		}

		//v.SaveInfo(ident.Name, blockStmt.Lbrace, blockStmt.Rbrace)
		// Save info `func ..... { xxxxx }`
		v.SaveInfo(ident.Name, funcType.Func, blockStmt.Rbrace)
	}
	return v
}

func (v visitorFun) SaveInfo(name string, lbrace, rbrace token.Pos) {
	bi := new(bodyInfo)
	bi.Lbrace = lbrace
	bi.Rbrace = rbrace
	v.locals[name] = bi
}

func (v visitorFun) PrintInfo() {
	for k, val := range v.locals {
		println(k, " ", val.Lbrace, ":", val.Rbrace)
	}
}

func (v visitorFun) Range(f func(key, value interface{}) bool) {
	for k, e := range v.locals {
		if !f(k, e) {
			break
		}
	}
}
