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

// Delete the characters defined in chr inside str
func stripchars(str, chr string) string {
    return strings.Map(func(r rune) rune {
        if strings.IndexRune(chr, r) < 0 {
            return r
        }
        return -1
    }, str)
}

// Sort the map by its values (int)
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

// Defines the structure Pair which is a part of a more global hash table
type Pair struct {
  Key string
  Value float64
}

// Defines the struture of a pseudo hash table
type PairList []Pair

// Few useful functions for PairList
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

// Prints the full PairList structure
func (p PairList) PrettyPrint(n int) {
    for pair := 0; pair < n; pair++ {
        fmt.Printf("Key:\t%v\t\t\tvalue:\t%.4f\n",p[pair].Key,  p[pair].Value)
    }
}

// Joins two maps together
func join(m map[string]int, toJoin map[string]int) map[string]int {
    for k := range toJoin {
        m[k] += toJoin[k]
    }
    return m
}

// Calculate the TFIDF of a pairlist containing all the words and a
// map array containing maps of each text words
func (p PairList) TFIDF(m []map[string]int) PairList {
    ref := p[0].Value
    for pair := range p {
        count := 0
        for text := range m {
            if m[text][p[pair].Key] != 0 { // if the word exists in another map
                count += 1
            }
        }
        p[pair].Value = 100 * (math.Log10(float64(len(m) / count))) * float64(p[pair].Value) / float64(ref)
    }
    return p
}

func main() {

    n := flag.Int("n", 20, "The value of n defines the length of the output")
    flag.Parse()
    texts := flag.Args()
    if len(texts) == 0 { // error, the user did not enter any input
        fmt.Println("No input provided. Job done")
        return
    }
    TextsCounts := make([]map[string]int, len(texts)) // structure to store the results
    for text := range texts {
        fmt.Printf("=================== PROCESSING %s ===================\n", texts[text])
        file, err := os.Open(texts[text]) // file opening
        if err != nil {
            log.Fatal(err)
        }
        wordMap := make(map[string]int)
        var result string
        for {
            buf := make([]byte, 100) // buffer
            count, err := file.Read(buf)
            if count == 0 { // done reading
                break
            }
            if err != nil { // Reading error
                log.Fatal(err)
            }
            content := string(buf[:count]) // converts to string
            content = strings.ToLower(stripchars(content, ",?;.:/!\"'\\+-*#`")) // delete all the characters corresponding in the string
            result += content // concatenate with result
        }
        contentFields := strings.Fields(result) // extract fields (space separated words)
        for field := range contentFields {
            wordMap[contentFields[field]] += 1 // add one to the hash map for each field to count appearing of the field
        }
        TextsCounts[text] = wordMap // Store in the map collection
        fmt.Printf("======================= DONE ========================\n")
    }
    union := make(map[string]int) // Defines the variable that holds the union of all the hash maps
    union = TextsCounts[0]
    for i := 1; i < len(texts); i++ {
        union = join(union, TextsCounts[i])
    }
    unionPairs := rankByWordCount(union) // Sort by count
    unionPairs = unionPairs.TFIDF(TextsCounts) // Calculate the TFIDF and store it in place of counts
    sort.Sort(unionPairs) // Sort again
    unionPairs.PrettyPrint(int(math.Min(float64(len(unionPairs)), float64(*n))))
}
