package scraper

import (
	"github.com/buke/quickjs-go"
	"search/internal/manga"
	"strconv"
	"strings"
)

func ValidateUrl(url string) bool {
	return false
}

func IsEvalFunction(script string) bool {
	script = strings.TrimSpace(script)
	return len(script) > 4 && script[:4] == "eval"
}

func DropLastPageIfBroken(pages *[]manga.Page) *[]manga.Page {
	lastTwo := (*pages)[len(*pages)-2:]
	var pageNums []int
	for _, p := range lastTwo {
		if p.Url != "" {
			v := strings.Split(p.Url, "/")
			s := v[len(v)-1]
			pn := strings.Split(s, ".")[0]
			pageNumber, err := strconv.Atoi(pn[len(pn)-2:])
			if err != nil {
				*pages = (*pages)[:len(*pages)-1]
				return pages
			}
			pageNums = append(pageNums, pageNumber)
		} else {
			// remove last page
			*pages = (*pages)[:len(*pages)-1]
			return pages
		}
	}

	if pageNums[0] == 0 && pageNums[1] == 1 {
		return pages
	} else if pageNums[1]-pageNums[0] == 1 {
		return pages
	}
	*pages = (*pages)[:len(*pages)-1]
	return pages
}

func ExtractSecretKey(doc string) (string, error) {
	runtime := quickjs.NewRuntime()
	defer runtime.Close()
	context := runtime.NewContext()
	defer context.Close()

	scriptStart := strings.Index(doc, "eval(function(p,a,c,k,e,d)")
	scriptEnd := strings.Index(doc[scriptStart:], "</script>")
	script := doc[scriptStart : scriptStart+scriptEnd]
	script = strings.TrimPrefix(script, "eval")

	s, err := context.Eval(script)
	if err != nil {
		return "", err
	}
	deObfuscatedScript := s.String()

	keyStart := strings.Index(deObfuscatedScript, "'")
	keyEnd := strings.Index(deObfuscatedScript, ";")
	keyString := deObfuscatedScript[keyStart:keyEnd]

	key, err := context.Eval(keyString)
	if err != nil {
		return "", nil
	}

	return key.String(), nil
}
