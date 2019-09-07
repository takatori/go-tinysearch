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

// TODO: 名前考える
func converter(r rune) rune {
	// 英数字以外だったら捨てる
	if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && !unicode.IsNumber(r) {
		return -1
	}
	// 大文字を小文字に変換する
	return unicode.ToLower(r)
}

// io.Readerから読んだデータをトークンに分割する関数
func (t *Tokenizer) SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
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
