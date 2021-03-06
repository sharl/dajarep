package dajarep

import (
	"fmt"
	"unicode"
	"math"
	"regexp"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

func init() {
	tokenizer.SysDic()
}

//単語
type word struct {
	str   string
	kana  string
	wtype string
}

//文章
type sentence struct {
	str   string
	kana  string
	words []word
}

//Dajarep :駄洒落を返す
func Dajarep(text string, debug bool) (dajares []string, debugStrs []string) {
	sentencesN := getSentences(text, tokenizer.Normal)
	sentencesS := getSentences(text, tokenizer.Search)
	for i := 0; i < len(sentencesN); i++ {
		if ok, kana := isDajare(sentencesN[i], debug); ok == true {
			dajares = append(dajares, sentencesN[i].str)
			debugStrs = append(debugStrs, kana)
		} else if ok, kana = isDajare(sentencesS[i], debug); ok == true {
			dajares = append(dajares, sentencesS[i].str)
			debugStrs = append(debugStrs, kana)
		}
	}
	return dajares, debugStrs
}

//駄洒落かどうかを評価する。
func isDajare(sen sentence, debug bool) (bool, string) {
	words := sen.words
	for i := 0; i < len(words); i++ {
		w := words[i]
		if debug {
			fmt.Println(w)
		}
		if w.wtype == "名詞" && len([]rune(w.kana)) > 1 {
			rStr := regexp.MustCompile(w.str)
			rKana := regexp.MustCompile(fixWord(w.kana))
			hitStr := rStr.FindAllString(sen.str, -1)
			hitKana := rKana.FindAllString(sen.kana, -1)
			hitKana2 := rKana.FindAllString(fixSentence(sen.kana), -1)
			//ある単語における　原文の一致文字列数<フリガナでの一致文字列数　→　駄洒落の読みが存在
			if debug {
				fmt.Println(rKana, len(hitStr), sen.kana, len(hitKana), fixSentence(sen.kana), len(hitKana2))
			}
			if len(hitStr) > 0 && len(hitStr) < int(math.Max(float64(len(hitKana)), float64(len(hitKana2)))) {
				return true, w.kana
			}
		}
	}
	return false, ""
}

//置き換え可能な文字を考慮した正規表現を返す。
func fixWord(text string) string {
	text = strings.Replace(text, "ッ", "[ツッ]?", -1)
	text = strings.Replace(text, "ァ", "[アァ]?", -1)
	text = strings.Replace(text, "ィ", "[イィ]?", -1)
	text = strings.Replace(text, "ゥ", "[ウゥ]?", -1)
	text = strings.Replace(text, "ェ", "[エェ]?", -1)
	text = strings.Replace(text, "ォ", "[オォ]?", -1)
	text = strings.Replace(text, "ズ", "[スズヅ]", -1)
	text = strings.Replace(text, "ヅ", "[ツズヅ]", -1)
	text = strings.Replace(text, "ヂ", "[チジヂ]", -1)
	text = strings.Replace(text, "ジ", "[シジヂ]", -1)
	text = strings.Replace(text, "ガ", "[カガ]", -1)
	text = strings.Replace(text, "ギ", "[キギ]", -1)
	text = strings.Replace(text, "グ", "[クグ]", -1)
	text = strings.Replace(text, "ゲ", "[ケゲ]", -1)
	text = strings.Replace(text, "ゴ", "[コゴ]", -1)
	text = strings.Replace(text, "ザ", "[サザ]", -1)
	text = strings.Replace(text, "ゼ", "[セゼ]", -1)
	text = strings.Replace(text, "ゾ", "[ソゾ]", -1)
	text = strings.Replace(text, "ダ", "[タダ]", -1)
	text = strings.Replace(text, "デ", "[テデ]", -1)
	text = strings.Replace(text, "ド", "[トド]", -1)
	re := regexp.MustCompile("[ハバパ]")
	text = re.ReplaceAllString(text, "[ハバパ]")
	re = regexp.MustCompile("[ヒビピ]")
	text = re.ReplaceAllString(text, "[ヒビピ]")
	re = regexp.MustCompile("[フブプ]")
	text = re.ReplaceAllString(text, "[フブプ]")
	re = regexp.MustCompile("[ヘベペ]")
	text = re.ReplaceAllString(text, "[ヘベペ]")
	re = regexp.MustCompile("[ホボポ]")
	text = re.ReplaceAllString(text, "[ホボポ]")
	re = regexp.MustCompile("([アカサタナハマヤラワャ])ー")
	text = re.ReplaceAllString(text, "$1[アァ]?")
	re = regexp.MustCompile("([イキシチニヒミリ])ー")
	text = re.ReplaceAllString(text, "$1[イィ]?")
	re = regexp.MustCompile("([ウクスツヌフムユルュ])ー")
	text = re.ReplaceAllString(text, "$1[ウゥ]?")
	re = regexp.MustCompile("([エケセテネへメレ])ー")
	text = re.ReplaceAllString(text, "$1[イィエェ]?")
	re = regexp.MustCompile("([オコソトノホモヨロヲョ])ー")
	text = re.ReplaceAllString(text, "$1[ウゥオォ]?")
	text = strings.Replace(text, "ャ", "[ヤャ]", -1)
	text = strings.Replace(text, "ュ", "[ユュ]", -1)
	text = strings.Replace(text, "ョ", "[ヨョ]", -1)
	text = strings.Replace(text, "ー", "[ー]?", -1)
	return text
}

//本文から省略可能文字を消したパターンを返す。
func fixSentence(text string) string {
	text = strings.Replace(text, "ッ", "", -1)
	text = strings.Replace(text, "ー", "", -1)
	text = strings.Replace(text, "、", "", -1)
	text = strings.Replace(text, ",", "", -1)
	text = strings.Replace(text, "　", "", -1)
	text = strings.Replace(text, " ", "", -1)
	return text
}

//テキストからsentenceオブジェクトを作る。
func getSentences(text string, mode tokenizer.TokenizeMode) []sentence {
	var sentences []sentence
	t := tokenizer.New()

	// http://www.serendip.ws/archives/6307
	kanaConv := unicode.SpecialCase {
		// ひらがなをカタカナに変換
		unicode.CaseRange{
			0x3041, // Lo: ぁ
			0x3093, // Hi: ん
			[unicode.MaxCase]rune{
				0x30a1 - 0x3041, // UpperCase でカタカナに変換
				0,               // LowerCase では変換しない
				0x30a1 - 0x3041, // TitleCase でカタカナに変換
			},
		},
	}

	text = strings.Replace(text, "。", "\n", -1)
	text = strings.Replace(text, ".", "\n", -1)
	text = strings.Replace(text, "?", "?\n", -1)
	text = strings.Replace(text, "!", "!\n", -1)
	text = strings.Replace(text, "？", "？\n", -1)
	text = strings.Replace(text, "！", "！\n", -1)
	senstr := strings.Split(text, "\n")

	for i := 0; i < len(senstr); i++ {
		tokens := t.Analyze(senstr[i], mode)
		var words []word
		var kana string
		for j := 0; j < len(tokens); j++ {
			tk := tokens[j]
			ft := tk.Features()
			if len(ft) > 7 {
				w := word{str: ft[6],
					kana:  ft[7],
					wtype: ft[0],
				}
				words = append(words, w)
				kana += ft[7]
			} else if len(ft) == 7 {
				lk := strings.ToUpperSpecial(kanaConv, tk.Surface)
				w := word{str: lk,
					kana: lk,
					wtype: ft[0],
				}
				words = append(words, w)
				kana += lk
			}
		}
		sentences = append(sentences,
			sentence{
				str:   senstr[i],
				words: words,
				kana:  kana,
			})
	}
	return sentences
}
