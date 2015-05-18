package c6

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import "fmt"
import "c6/ast"
import "path/filepath"
import "io/ioutil"

const (
	UnknownFileType = iota
	ScssFileType
	SassFileType
	EcssFileType
)

func debug(format string, args ...interface{}) {
	if debugParser {
		fmt.Printf(format+"\n", args...)
	}
}

func getFileTypeByExtension(extension string) uint {
	switch extension {
	case "scss":
		return ScssFileType
	case "sass":
		return SassFileType
	case "ecss":
		return EcssFileType
	}
	return UnknownFileType
}

type Parser struct {
	Context *Context
	Input   chan *ast.Token

	// integer for counting token
	Pos         int
	RollbackPos int
	Tokens      []*ast.Token
}

func NewParser(context *Context) *Parser {
	return &Parser{context, nil, 0, 0, []*ast.Token{}}
}

func (parser *Parser) ParseFile(path string) error {
	ext := filepath.Ext(path)
	filetype := getFileTypeByExtension(ext)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var code string = string(data)
	switch filetype {
	case ScssFileType:
		parser.ParseScss(code)
		break
	}
	return nil
}

func (self *Parser) backup() {
	self.Pos--
}

func (self *Parser) restore(pos int) {
	self.Pos = pos
}

func (self *Parser) remember() {
	self.RollbackPos = self.Pos
}

func (self *Parser) rollback() {
	self.Pos = self.RollbackPos
}

func (self *Parser) accept(tokenType ast.TokenType) *ast.Token {
	var tok = self.next()
	if tok.Type == tokenType {
		return tok
	}
	self.backup()
	return nil
}

func (self *Parser) acceptAny(tokenTypes ...ast.TokenType) *ast.Token {
	var tok = self.next()
	for _, tokType := range tokenTypes {
		if tok.Type == tokType {
			return tok
		}
	}
	self.backup()
	return nil
}

func (self *Parser) acceptAnyOf2(tokType1, tokType2 ast.TokenType) *ast.Token {
	var tok = self.next()
	if tok.Type == tokType1 || tok.Type == tokType2 {
		return tok
	}
	self.backup()
	return nil
}

func (self *Parser) acceptAnyOf3(tokType1, tokType2, tokType3 ast.TokenType) *ast.Token {
	var tok = self.next()
	if tok.Type == tokType1 || tok.Type == tokType2 || tok.Type == tokType3 {
		return tok
	}
	self.backup()
	return nil
}

func (self *Parser) expect(tokenType ast.TokenType) *ast.Token {
	var tok = self.next()
	if tok.Type != tokenType {
		self.backup()
		panic(SyntaxError{
			Expecting:   tokenType.String(),
			ActualToken: tok,
		})
		panic(fmt.Errorf("Expecting %s, Got %s", tokenType, tok))
	}
	return tok
}

func (self *Parser) next() *ast.Token {
	var p = self.Pos
	self.Pos++

	// if we've appended the token
	if p < len(self.Tokens) {
		return self.Tokens[p]
	}
	return nil
}

func (self *Parser) peekBy(offset int) *ast.Token {
	var i = 0
	var tok *ast.Token = nil
	for offset > 0 {
		tok = self.next()
		offset--
		i++
		if tok == nil {
			break
		}
	}
	self.Pos -= i
	return tok
}

func (self *Parser) advance() {
	self.Pos++
}

func (self *Parser) current() *ast.Token {
	return self.Tokens[self.Pos]
}

func (self *Parser) peek() *ast.Token {
	if self.Pos < len(self.Tokens) {
		return self.Tokens[self.Pos]
	}
	return nil
}

func (self *Parser) eof() bool {
	var tok = self.next()
	self.backup()
	return tok == nil
}
