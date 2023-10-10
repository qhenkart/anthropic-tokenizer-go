package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkoukk/tiktoken-go"
	"golang.org/x/text/unicode/norm"
)

type config struct {
	PatternStr     string         `json:"pat_str"`
	SpecialTokens  map[string]int `json:"special_tokens"`
	ExplicitNVocab int            `json:"explicit_n_vocab"`
}

func New() *tiktoken.Tiktoken {
	fileContent, err := os.ReadFile("claude.json")
	if err != nil {
		panic(err)
	}

	var conf config
	err = json.Unmarshal(fileContent, &conf)
	if err != nil {
		panic(err)
	}

	encoder := tiktoken.NewDefaultBpeLoader()
	ranks, err := encoder.LoadTiktokenBpe("./bpe_ranks.txt")
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

	return tt
}

func countTokens(text string) int {
	tt := New()

	normal := norm.NFKC.String(text)
	return len(tt.Encode(normal, []string{"all"}, nil))
}
