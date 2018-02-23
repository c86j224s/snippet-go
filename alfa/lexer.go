package alfa

import (
	"bytes"
	"fmt"
	"regexp"
)

// TokenType 에 대한 enumeration 입니다.
const (
	ERROR = iota
	WHITESPACE
	COMMENT

	STRINGLITERAL
	NUMERICLITERAL
	BOOLEANLITERAL

	IDENTIFIER
	REFERENCE
	OPERATOR

	CURLYBRACEBEGIN
	CURLYBRACEEND

	ROUNDBRACEBEGIN
	ROUNDBRACEEND

	EOF
)

// TokenType 은 uint 숫자로 표현됩니다.
type TokenType uint

// Token 은 타입과 슬라이스로 정의됩니다.
type Token struct {
	tokenType   TokenType
	sourcePiece []byte
}

// Lexer 는 현재 렉서 머신 상태를 가지고 있습니다.
type Lexer struct {
	source []byte
	cursor uint
	Tokens []Token
}

// NewLexerFromFile : 파일을 읽어들여, 이를 해석할 렉서를 만듭니다.
func NewLexerFromFile(fileName string) *Lexer {
	source, err := FileReadAll(fileName)
	if err != nil {
		return nil
	}

	lexer := Lexer{source, 0, make([]Token, 256)}
	return &lexer
}

// NewLexerFromData : 데이터를 해석할 렉서를 만듭니다.
func NewLexerFromData(source []byte) *Lexer {
	lexer := Lexer{source, 0, make([]Token, 0, 256)}
	return &lexer
}

// Take 는 코드에서 일치하는 부분을 잘라냅니다.
func (lexer *Lexer) Take(tokenType TokenType, fn func([]byte, uint) uint) *Token {
	cur := fn(lexer.source, lexer.cursor)
	if lexer.cursor == cur {
		return nil
	}

	newToken := &Token{tokenType, lexer.source[lexer.cursor:cur]}

	lexer.Tokens = append(lexer.Tokens, *newToken)
	lexer.cursor = cur
	return newToken
}

// IsWhitespace : 화이트스페이스인지 체크해서 화이트스페이스의 슬라이스를 구합니다.
func IsWhitespace(source []byte, cursor uint) uint {
	cur := cursor
	for cur < uint(len(source)) && (source[cur] == ' ' || source[cur] == '\t' || source[cur] == '\r' || source[cur] == '\n') {
		cur++
	}

	return cur
}

// IsComment 는 커멘트인지를 판단합니다.
func IsComment(source []byte, cursor uint) uint {
	cur := cursor
	if source[cur] == '/' && source[cur+1] == '/' {
		cur++

		for true {
			if source[cur] == '\n' {
				cur++
				break
			}
			cur++
		}

		return cur
	}

	if source[cur] == '/' && source[cur+1] == '*' {
		cur++
		for true {
			if source[cur-1] == '*' && source[cur] == '/' {
				cur++
				break
			}
			cur++
		}
	}

	return cur
}

// IsStringLiteral 은 문자열 리터럴인지 확인합니다.
func IsStringLiteral(source []byte, cursor uint) uint {
	cur := cursor
	if source[cur] == '"' {
		cur++
		for true {
			if cur <= uint(len(source)) && source[cur] == '"' && source[cur-1] != '\\' {
				break
			}
			cur++
		}
	}

	return cur
}

// IsNumericLiteral 은 숫자 리터럴인지 확인합니다.
// 아직 exponential expression은 지원하지 않습니다.
func IsNumericLiteral(source []byte, cursor uint) uint {
	cur := cursor
	if re, err := regexp.Compile(`[-+]?\d*\.?\d+([eE][-+]?\d+)?`); err == nil {
		matched := re.FindIndex(source[cur:])

		if matched != nil && matched[0] == 0 {
			cur += uint(matched[1])
		}
	}

	return cur
}

// IsIdentifier 는 식별자인지 확인합니다.
// 이 이름이 적절한지, 그리고 Reference 와 분리해서 관리하는 게 맞는지는 고민이 필요합니다.
func IsIdentifier(source []byte, cursor uint) uint {
	cur := cursor
	re := regexp.MustCompile(`[a-zA-Z_][0-9a-zA-Z_]*`)

	matched := re.FindIndex(source[cur:])

	if matched != nil && matched[0] == 0 {
		cur += uint(matched[1])
	}

	return cur
}

// IsReference 는 네임스페이스의 중첩 표현을 갖는 식별자를 확인합니다.
// 이 이름이 적절한지, 그리고 Identifier 와 분리해서 관리하는 게 맞는지는 고민이 필요합니다.
func IsReference(source []byte, cursor uint) uint {
	cur := cursor
	re := regexp.MustCompile(`[a-zA-Z_][0-9a-zA-Z_]*(\.(\*|([a-zA-Z_][0-9a-zA-Z_]*)))+`)

	matched := re.FindIndex(source[cur:])

	if matched != nil && matched[0] == 0 {
		cur += uint(matched[1])
	}

	return cur
}

// IsOperator 는 연산자인지를 확인합니다.
func IsOperator(source []byte, cursor uint) uint {
	operators := [][]byte{
		[]byte("or"),
		[]byte("and"),
		[]byte("|"),
		[]byte("&"),
		[]byte("<="),
		[]byte("<"),
		[]byte(">="),
		[]byte(">"),
		[]byte("$"),
		[]byte("=="),
		[]byte("="),
		[]byte("@"),
		[]byte("^"),
		[]byte("+"),
		[]byte("-"),
		[]byte("*"),
		[]byte("/"),
		[]byte("%"),
	}

	for _, op := range operators {
		if bytes.HasPrefix(source, op) {
			return cursor + uint(len(op))
		}
	}

	return cursor
}

// IsCurlyBracedBegin 은 블록괄호의 시작을 확인합니다.
func IsCurlyBracedBegin(source []byte, cursor uint) uint {
	if source[cursor] == '{' {
		return cursor + 1
	}

	return cursor
}

// IsCurlyBracedEnd 는 블록괄호의 끝을 확인합니다.
func IsCurlyBracedEnd(source []byte, cursor uint) uint {
	if source[cursor] == '}' {
		return cursor + 1
	}
	return cursor
}

// IsRoundBracedBegin 은 괄호의 시작을 확인합니다.
func IsRoundBracedBegin(source []byte, cursor uint) uint {
	if source[cursor] == '(' {
		return cursor + 1
	}
	return cursor
}

// IsRoundBracedEnd 는 괄호의 끝을 확인합니다.
func IsRoundBracedEnd(source []byte, cursor uint) uint {
	if source[cursor] == ')' {
		return cursor + 1
	}
	return cursor
}

// IsEOF 는 데이터의 끝인지를 확인합니다.
func IsEOF(source []byte, cursor uint) bool {
	if cursor >= uint(len(source)) {
		return true
	}
	return false
}

// NextToken : 토큰을 하나 떼어냅니다.
func (lexer *Lexer) NextToken() bool {
	var nextToken *Token

	if IsEOF(lexer.source, lexer.cursor) {
		return false
	}

	if cur := IsWhitespace(lexer.source, lexer.cursor); lexer.cursor != cur {
		lexer.cursor = cur
		return true
	}

	type Tokenizer struct {
		tok TokenType
		fn  func([]byte, uint) uint
	}
	tokenizer := []Tokenizer{
		Tokenizer{COMMENT, IsComment},
		Tokenizer{REFERENCE, IsReference},
		Tokenizer{IDENTIFIER, IsIdentifier},
		Tokenizer{STRINGLITERAL, IsStringLiteral},
		Tokenizer{NUMERICLITERAL, IsNumericLiteral},
		Tokenizer{OPERATOR, IsOperator},
		Tokenizer{CURLYBRACEBEGIN, IsCurlyBracedBegin},
		Tokenizer{CURLYBRACEEND, IsCurlyBracedEnd},
		Tokenizer{ROUNDBRACEBEGIN, IsRoundBracedBegin},
		Tokenizer{ROUNDBRACEEND, IsRoundBracedEnd},
	}
	for _, each := range tokenizer {
		if nextToken = lexer.Take(each.tok, each.fn); nextToken != nil {
			return true
		}
	}

	lexer.Tokens = append(lexer.Tokens, Token{ERROR, lexer.source[lexer.cursor : lexer.cursor+1]})
	lexer.cursor++
	return true
}

// Dump 는 현재 렉서에 저장된 토큰을 덤프합니다.
func (lexer *Lexer) Dump() {
	for _, tok := range lexer.Tokens {
		fmt.Println(tok.tokenType, " : ", string(tok.sourcePiece))
	}
}
