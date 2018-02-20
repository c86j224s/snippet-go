package alfa_test

import (
	"testing"

	"github.com/c86j224s/snippet-go/alfa"
)

// TestComment 는 유닛테스트입니다. 3.4.1 예제를 테스트합니다.
func TestComment(t *testing.T) {
	exampleSource := []byte(`
	// This is a single line comment

/* This comment spans 3 lines
   Line 2 of 3-line comment
*/
	`)

	lexer := alfa.NewLexerFromData(exampleSource)
	for lexer.NextToken() {
	}

	if len(lexer.Tokens) != 2 {
		t.Errorf("unexpected lexer.Tokens : %d", len(lexer.Tokens))
	}
}

// TestNamespace 는 유닛테스트입니다. 3.5.1 예제를 테스트합니다.
func TestNamespace(t *testing.T) {
	exampleSource := []byte(`
		namespace documents {
			policy level1 {
			}
		}
	`)

	lexer := alfa.NewLexerFromData(exampleSource)
	for lexer.NextToken() {
	}

	if len(lexer.Tokens) != 8 {
		t.Errorf("unexpected lexer.Tokens : %d", len(lexer.Tokens))
	}
}

// TestNamespace2 는 3.6.1 예제를 테스트합니다.
func TestNamespace2(t *testing.T) {
	exampleSource := []byte(`
		namespace A {
			namespace B {
				policy P {
				}
			}
			import B.P
		}
		namespace C {
			import A.B.*
		}
	`)

	lexer := alfa.NewLexerFromData(exampleSource)
	for lexer.NextToken() {
	}

	if len(lexer.Tokens) != 20 {
		lexer.Dump()
		t.Errorf("unexpected lexer.Tokens : %d", len(lexer.Tokens))
	}
}

/*
// TestRule 은 3.7.1 예제를 테스트합니다.
func TestRule(t *testing.T) {
	exampleSource := []byte(`
		rule {
			permit
			target clause Attributes.resourceType == "document"
			condition Attributes.userClearance >= Attributes.resourceClassification
		}
	`)

	lexer := alfa.NewLexerFromData(exampleSource)
	for lexer.NextToken() {
	}

	if len(lexer.Tokens) != 12 {
		t.Errorf("unexpected lexer.Tokens : %d", len(lexer.Tokens))
	}

}

*/