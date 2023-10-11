package tokenizer

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	_ "embed"

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

//go:embed claude.json
var claudeJSON []byte

//go:embed bpe_ranks.txt
var bpeRanks string

// New creates a new Tokenizer. This should be initialized on server start up. Running this for every request could cause memory failures
func New() (*Tokenizer, error) {
	var conf config
	if err := json.Unmarshal(claudeJSON, &conf); err != nil {
		return nil, err
	}

	ranks, err := loadBPE(bpeRanks)
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

func loadBPE(contents string) (map[string]int, error) {
	bpeRanks := make(map[string]int)
	for _, line := range strings.Split(contents, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		token, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, err
		}
		rank, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		bpeRanks[string(token)] = rank
	}
	return bpeRanks, nil
}
