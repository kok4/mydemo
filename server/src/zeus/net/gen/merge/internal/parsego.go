package internal

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
	FunBeginPos int // point to "func"
	FunEndPos   int // point to '}'
}

type GoParser struct {
	fset *token.FileSet
	f    *ast.File
}

type FuncInfoMap map[string]*FunInfo

func NewGoParserFromFile(filePath string) *GoParser {
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

func NewGoParserFromStr(content string) *GoParser {
	obj := &GoParser{
		fset: token.NewFileSet(), // positions are relative to fset
	}

	var err error
	obj.f, err = parser.ParseFile(obj.fset, "", content, parser.AllErrors)
	if err != nil {
		panic(fmt.Sprintf("parse from string error: %s", err))
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
		obj.FunBeginPos = pg.fset.Position(ret.FunBeginPos).Offset
		obj.FunEndPos = pg.fset.Position(ret.FunEndPos).Offset

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

func (pg *GoParser) GetFuncInfoMap() FuncInfoMap {
	result := make(map[string]*FunInfo)
	pg.GetFunctions(func(retKey, retVal interface{}) bool {
		bi := new(FunInfo)
		rv := retVal.(*FunInfo)
		bi.FunBeginPos = rv.FunBeginPos
		bi.FunEndPos = rv.FunEndPos
		result[retKey.(string)] = bi
		return true
	})
	return result
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
	FunBeginPos token.Pos
	FunEndPos   token.Pos
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

		// Save info `func ... { ... }`
		v.SaveInfo(ident.Name, funcType.Func, blockStmt.Rbrace)
	}
	return v
}

func (v visitorFun) SaveInfo(name string, funcBegin, rbrace token.Pos) {
	bi := new(bodyInfo)
	bi.FunBeginPos = funcBegin
	bi.FunEndPos = rbrace
	v.locals[name] = bi
}

func (v visitorFun) PrintInfo() {
	for k, val := range v.locals {
		println(k, " ", val.FunBeginPos, ":", val.FunEndPos)
	}
}

func (v visitorFun) Range(f func(key, value interface{}) bool) {
	for k, e := range v.locals {
		if !f(k, e) {
			break
		}
	}
}
