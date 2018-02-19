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
	tokenType   TokenType
	sourcePiece []byte
}

// Lexer 는 현재 렉서 머신 상태를 가지고 있습니다.
type Lexer struct {
	source []byte
	cursor uint
	tokens []Token
}

// NewLexer : 파일을 읽어들여, 이를 해석할 렉서를 만듭니다.
func NewLexer(fileName string) *Lexer {
	source, err := FileReadAll(fileName)
	if err != nil {
		return nil
	}

	lexer := Lexer{source, 0, make([]Token, 256)}
	return &lexer
}

// Take 는 코드에서 일치하는 부분을 잘라냅니다.
func (lexer *Lexer) Take(tokenType TokenType, fn func([]byte, uint) uint) *Token {
	cur := fn(lexer.source, lexer.cursor)
	if lexer.cursor == cur {
		return nil
	}

	newToken := &Token{tokenType, lexer.source[lexer.cursor:cur]}
	lexer.cursor = cur
	return newToken
}

// IsWhitespace : 화이트스페이스인지 체크해서 화이트스페이스의 슬라이스를 구합니다.
func IsWhitespace(source []byte, cursor uint) uint {
	cur := cursor
	for true {
		switch source[cur] {
		case ' ', '\t', '\r', '\n':
			cur++
		default:
			break
		}
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
			if source[cur] == '*' && source[cur+1] == '/' {
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
			if source[cur] == '"' && source[cur-1] != '\\' {
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

		if matched[0] == 0 {
			cur += uint(matched[1])
		}
	}

	return cur
}

// IsIdentifier 는 식별자인지 확인합니다.
func IsIdentifier(source []byte, cursor uint) uint {
	cur := cursor
	if re, err := regexp.Compile(`[a-zA-Z_][0-9a-zA-Z_]*`); err == nil {
		matched := re.FindIndex(source[cur:])

		if matched[0] == 0 {
			cur += uint(matched[1])
		}
	}

	return cur
}

// NextToken : 토큰을 하나 떼어냅니다.
func (lexer *Lexer) NextToken() {
	var nextToken *Token

	tokenizer := map[TokenType]func([]byte, uint) uint{
		WHITESPACE:     IsWhitespace,
		COMMENT:        IsComment,
		IDENTIFIER:     IsIdentifier,
		STRINGLITERAL:  IsStringLiteral,
		NUMERICLITERAL: IsNumericLiteral,
	}
	for tok, fn := range tokenizer {
		if nextToken = lexer.Take(tok, fn); nextToken != nil {
			return
		}
	}

	lexer.tokens = append(lexer.tokens, Token{ERROR, lexer.source[lexer.cursor : lexer.cursor+1]})
	lexer.cursor++
}
