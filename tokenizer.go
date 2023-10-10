package tokenizer

import (
	"encoding/json"
	"os"

	"github.com/pkoukk/tiktoken-go"
	"golang.org/x/text/unicode/norm"
)

type config struct {
	PatternStr     string         `json:"pat_str"`
	SpecialTokens  map[string]int `json:"special_tokens"`
	ExplicitNVocab int            `json:"explicit_n_vocab"`
}

// Tokenizer wraps the underlying tiktoken library
type Tokenizer struct {
	*tiktoken.Tiktoken
}

// New creates a new Tokenizer. This should be initialized on server start up. Running this for every request could cause memory failures
func New() (*Tokenizer, error) {
	fileContent, err := os.ReadFile("claude.json")
	if err != nil {
		return nil, err
	}

	var conf config
	if err = json.Unmarshal(fileContent, &conf); err != nil {
		return nil, err
	}
	encoder := tiktoken.NewDefaultBpeLoader()
	ranks, err := encoder.LoadTiktokenBpe("./bpe_ranks.txt")
	if err != nil {
		return nil, err
	}

	bpe, err := tiktoken.NewCoreBPE(ranks, conf.SpecialTokens, conf.PatternStr)
	if err != nil {
		return nil, err
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

	return &Tokenizer{tt}, nil
}

// Tokens returns the amount of tokens in the input
func (t *Tokenizer) Tokens(text string) int {
	normal := norm.NFKC.String(text)
	return len(t.Encode(normal, []string{"all"}, nil))
}
