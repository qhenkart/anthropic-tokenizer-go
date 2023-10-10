// Run via the Makefile. Automatically pulls the latest json configuration file from Anthropic's tokenizer repo and modifies it into a processable state by the underlying tokenizer logic
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const configFileLocation = "https://raw.githubusercontent.com/anthropics/anthropic-tokenizer-typescript/main/claude.json"

type config struct {
	BPERanks       string         `json:"bpe_ranks"`
	PatternStr     string         `json:"pat_str"`
	SpecialTokens  map[string]int `json:"special_tokens"`
	ExplicitNVocab int            `json:"explicit_n_vocab"`
}

func main() {
	// pull the file from Github
	content, err := fetchSourceConfig()
	if err != nil {
		panic(err)
	}

	var conf config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		panic(err)
	}

	// Process the string
	bytePairs := strings.Split(conf.BPERanks[4:], " ")

	file, err := os.Create("bpe_ranks.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	for i, pair := range bytePairs {
		file.WriteString(fmt.Sprintf("%s %d\n", pair, i))
	}

	conf.BPERanks = ""

	modifiedContent, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("claude.json", modifiedContent, 0644); err != nil {
		panic(err)
	}
}

// retrieves the config file from the main branch of Anthropics tokenizer
func fetchSourceConfig() ([]byte, error) {
	url := configFileLocation
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download file")
	}

	file, err := os.Create("temp.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	defer func() {
		os.Remove("temp.json")
	}()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return nil, err
	}

	return os.ReadFile("./temp.json")
}
