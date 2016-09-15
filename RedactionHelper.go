package main

import (
    "fmt"
    "math"
    "flag"
    "os"
    "strings"
    "log"
    "sort"
)

func stripchars(str, chr string) string {
    return strings.Map(func(r rune) rune {
        if strings.IndexRune(chr, r) < 0 {
            return r
        }
        return -1
    }, str)
}

func rankByWordCount(wordFrequencies map[string]int) PairList{
  pl := make(PairList, len(wordFrequencies))
  i := 0
  for k, v := range wordFrequencies {
    pl[i] = Pair{k, float64(v)}
    i++
  }
  sort.Sort(pl)
  return pl
}

type Pair struct {
  Key string
  Value float64
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

func (p PairList) PrettyPrint() {
    for pair := range p {
        fmt.Printf("Key:\t%v\t\t\tvalue:\t%.4f\n",p[pair].Key,  p[pair].Value)
    }
}

func join(m map[string]int, toJoin map[string]int) map[string]int {
    for k := range toJoin {
        m[k] += toJoin[k]
    }
    return m
}

func (p PairList) TFIDF(m []map[string]int) PairList {
    ref := p[0].Value
    for pair := range p {
        count := 0
        for text := range m {
            if m[text][p[pair].Key] != 0 {
                count += 1
            }
        }
        p[pair].Value = 100 * (math.Log10(float64(len(m) / count))) * float64(p[pair].Value) / float64(ref)
    }
    return p
}

func main() {

    flag.Parse()
    texts := flag.Args()
    TextsCounts := make([]map[string]int, len(texts))
    for text := range texts {
        fmt.Printf("=================== PROCESSING %s ===================\n", texts[text])
        file, err := os.Open(texts[text])
        if err != nil {
            log.Fatal(err)
        }
        wordMap := make(map[string]int)
        var result string
        for {
            buf := make([]byte, 100)
            count, err := file.Read(buf)
            if count == 0 {
                break
            }
            if err != nil {
                log.Fatal(err)
            }
            content := string(buf[:count])
            content = strings.ToLower(stripchars(content, ",?;.:/!\"'\\+-*#`"))
            result += content
        }
        contentFields := strings.Fields(result)
        for field := range contentFields {
            wordMap[contentFields[field]] += 1
        }
        TextsCounts[text] = wordMap
        fmt.Printf("======================= DONE ========================\n")
    }
    union := make(map[string]int)
    union = TextsCounts[0]
    for i := 1; i < len(texts); i++ {
        union = join(union, TextsCounts[i])
    }
    unionPairs := rankByWordCount(union)
    unionPairs = unionPairs.TFIDF(TextsCounts)
    sort.Sort(unionPairs)
    unionPairs.PrettyPrint()


}
