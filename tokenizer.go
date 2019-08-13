package tinysearch

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
	"unicode"
)

var ignoreCharsRegxp = regexp.MustCompile("['!,?.:{}()|\\-+<>\\][/_]")
var whitespaceRegexp = regexp.MustCompile("\\s+")

// ドキュメントをトークンに分割する関数
func TextToWordSequence(document string) []string {

	// 大文字を小文字に変換
	document = strings.ToLower(document)

	// 不要な文字を削除 TODO: refactoring strings.Replacer
	document = ignoreCharsRegxp.ReplaceAllString(document, "")

	// 一つ以上のスペースで分割
	terms := whitespaceRegexp.Split(document, -1)

	for i, term := range terms {
		terms[i] = strings.Trim(term, "\n")
	}

	return terms
}

// トークンに分割する関数
func Analyzer(data []byte, atEOF bool) (advance int, token []byte, err error) {

	advance, token, err = bufio.ScanWords(data, atEOF)

	myAnalyzer := func(r rune) rune {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return -1
		}
		return unicode.ToLower(r)
	}

	if err == nil && token != nil {
		token = bytes.Map(myAnalyzer, token)
	}

	return
}
