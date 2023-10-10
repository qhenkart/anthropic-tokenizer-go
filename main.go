package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkoukk/tiktoken-go"
	"golang.org/x/text/unicode/norm"
)

type Config struct {
	PatternStr    string         `json:"pat_str"`
	SpecialTokens map[string]int `json:"special_tokens"`
	//BPERanks       map[string]int `json:"bpe_ranks"`
	//BPERanks       string `json:"bpe_ranks"`
	ExplicitNVocab int `json:"explicit_n_vocab"`
}

func main() {
	fileContent, err := os.ReadFile("./claude.json")
	if err != nil {
		panic(err)
	}

	var conf Config
	err = json.Unmarshal(fileContent, &conf)
	if err != nil {
		panic(err)
	}

	//bpe, err := tiktoken.NewCoreBPE(conf.BPERanks, conf.SpecialTokens, conf.PatternStr)
	//if err != nil {
	//err = fmt.Errorf("getEncoding: %v", err)
	//return
	//}

	//fmt.Println(bpe)

	//// encode
	//token := tke.Encode(text, nil, nil)

	////tokens
	//fmt.Println((token))
	//// num_tokens
	//fmt.Println(len(token))
	//updateJson2(conf)

}

func getTokenizer() *tiktoken.Tiktoken {
	fileContent, err := os.ReadFile("claude2.json")
	if err != nil {
		panic(err)
	}

	var conf Config
	err = json.Unmarshal(fileContent, &conf)
	if err != nil {
		panic(err)
	}

	encoder := tiktoken.NewDefaultBpeLoader()
	ranks, err := encoder.LoadTiktokenBpe("./output.txt")
	if err != nil {
		panic(err)
	}

	bpe, err := tiktoken.NewCoreBPE(ranks, conf.SpecialTokens, conf.PatternStr)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		panic(err)
	}

	enc := &tiktoken.Encoding{
		Name:           "CLAUDE-2",
		PatStr:         conf.PatternStr,
		MergeableRanks: ranks,
		SpecialTokens:  conf.SpecialTokens,
		ExplicitNVocab: conf.ExplicitNVocab,
	}

	specialTokensSet := map[string]any{}
	for k := range enc.SpecialTokens {
		specialTokensSet[k] = true
	}

	tt := tiktoken.NewTiktoken(bpe, enc, specialTokensSet)

	fmt.Printf("\n\n%+v\n", tt)
	return tt

	// const tokenizer = getTokenizer();
	// const encoded = tokenizer.encode(text.normalize('NFKC'), 'all');
	// tokenizer.free();
	// return encoded.length;
	// }
	//
	//export function countTokens(text: string): number {
}

func countTokens(text string) int {
	tt := getTokenizer()

	normal := norm.NFKC.String(text)
	return len(tt.Encode(normal, []string{"all"}, nil))
}

//func updateJson2(conf Config) {
//// Process the string
//bytePairs := strings.Split(conf.BPERanks, " ")

//file, err := os.Create("./output.txt")
//if err != nil {
//panic(err)
//}
//defer file.Close()

//for i, pair := range bytePairs {
//file.WriteString(fmt.Sprintf("%s %d\n", pair, i))
//}

//}

//func updateJson(conf Config) {
////	 Process the string
//bytePairs := strings.Split(conf.BPERanks, " ")
//rank := 0
//bpeRanksMap := make(map[string]int)

//for _, pair := range bytePairs {
//bpeRanksMap[pair] = rank
//rank++
//}

//out := map[string]interface{}{"pat_str": conf.PatternStr, "special_tokens": conf.SpecialTokens, "bpe_ranks": bpeRanksMap}
//// Save the result back to the same file
//jsonData, err := json.MarshalIndent(out, "", "  ")
//if err != nil {
//fmt.Printf("Error encoding JSON: %s", err)
//return
//}

//err = os.WriteFile("claude2.json", jsonData, 0644)
//if err != nil {
//fmt.Printf("Error writing to file: %s", err)
//return
//}

//}
