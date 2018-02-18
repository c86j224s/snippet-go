package alfa

import "regexp"

// TokenType 에 대한 enumeration 입니다.
const (
	ERROR = iota
	WHITESPACE
	COMMENT

	STRINGLITERAL
	NUMERICLITERAL
	BOOLEANLITERAL

	KEYWORD
	IDENTIFIER

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
	tokenType TokenType
	codePiece []byte
}

// Lexer 는 현재 렉서 머신 상태를 가지고 있습니다.
type Lexer struct {
	code   []byte
	cursor uint
	tokens []Token
}

// NewLexer : 파일을 읽어들여, 이를 해석할 렉서를 만듭니다.
func NewLexer(fileName string) *Lexer {
	code, err := FileReadAll(fileName)
	if err != nil {
		return nil
	}

	lexer := Lexer{code, 0, make([]Token, 256)}
	return &lexer
}

// Take 는 코드에서 일치하는 부분을 잘라냅니다.
func (lexer *Lexer) Take(tokenType TokenType, fn func([]byte, uint) uint) *Token {
	cur := fn(lexer.code, lexer.cursor)
	if lexer.cursor == cur {
		return nil
	}

	newToken := &Token{tokenType, lexer.code[lexer.cursor:cur]}
	lexer.cursor = cur
	return newToken
}

// IsWhitespace : 화이트스페이스인지 체크해서 화이트스페이스의 슬라이스를 구합니다.
func IsWhitespace(code []byte, cursor uint) uint {
	cur := cursor
	for true {
		switch code[cur] {
		case ' ', '\t', '\r', '\n':
			cur++
		default:
			break
		}
	}

	return cur
}

// IsComment 는 커멘트인지를 판단합니다.
func IsComment(code []byte, cursor uint) uint {
	cur := cursor
	if code[cur] == '/' && code[cur+1] == '/' {
		cur++
		for true {
			if code[cur] == '\n' {
				cur++
				break
			}
			cur++
		}

		return cur
	}

	if code[cur] == '/' && code[cur+1] == '*' {
		cur++
		for true {
			if code[cur] == '*' && code[cur+1] == '/' {
				cur++
				break
			}
			cur++
		}
	}

	return cur
}

// IsStringLiteral 은 문자열 리터럴인지 확인합니다.
func IsStringLiteral(code []byte, cursor uint) uint {
	cur := cursor
	if code[cur] == '"' {
		cur++
		for true {
			if code[cur] == '"' && code[cur-1] != '\\' {
				break
			}
			cur++
		}
	}

	return cur
}

// IsNumericLiteral 은 숫자 리터럴인지 확인합니다.
// 아직 exponential expression은 지원하지 않습니다.
func IsNumericLiteral(code []byte, cursor uint) uint {
	cur := cursor
	re, err := regexp.Compile(`[-+]?\d*\.?\d+([eE][-+]?\d+)?`)
	if err != nil {
		return cur
	}

	matched := re.FindIndex(code[cur:])

	if matched[0] == 0 {
		cur += uint(matched[1])
	}

	return cur
}

func IsKeyword(code []byte, cursor uint) uint {

}

// NextToken : 토큰을 하나 떼어냅니다.
func (lexer *Lexer) NextToken() {
	var nextToken *Token

	if nextToken = lexer.Take(WHITESPACE, IsWhitespace); nextToken != nil {
		return
	}
	if nextToken = lexer.Take(COMMENT, IsComment); nextToken != nil {
		return
	}
	if nextToken = lexer.Take(STRINGLITERAL, IsStringLiteral); nextToken != nil {
		return
	}
	if nextToken = lexer.Take(NUMERICLITERAL, IsNumericLiteral); nextToken != nil {
		return
	}

	lexer.tokens = append(lexer.tokens, Token{ERROR, lexer.code[lexer.cursor : lexer.cursor+1]})
	lexer.cursor++
}
