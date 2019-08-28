package tinysearch

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
)

type Tokenizer struct{}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

// io.Readerから読んだデータをトークンに分割する関数
func (t *Tokenizer) SplitFunc(data []byte, atEOF bool) (advance int,
	token []byte, err error) {

	advance, token, err = bufio.ScanWords(data, atEOF)

	converter := func(r rune) rune {
		if (r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			!unicode.IsNumber(r) {
			return -1
		}
		return unicode.ToLower(r)
	}

	if err == nil && token != nil {
		token = bytes.Map(converter, token)
	}

	return
}

// 文字列を分解する処理
func (t *Tokenizer) TextToWordSequence(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(t.SplitFunc)
	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}
