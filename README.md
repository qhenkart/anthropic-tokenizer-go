# Anthropic Golang Tokenizer

This package provides a convenient way to check how many tokens a given piece of text will be.

## Installation

```sh
go get github.com/qhenkart/anthropic-tokenizer-go
```

## Usage

Do not create a new Tokenizer on each request as it's quite expensive as it has somewhere in the range of 150,000 iterations to initialize the underlying byte pairs.
Create one on start up and share it

A typical pattern might look like this

```
type AnthropicClient struct {
 Tokenizer *tokenizer.Tokenizer
}


tokenizer := tokenizer.New()

client := &AnthropicClient{
  Tokenizer: tokenizer
}


func (c *client) generate(prompt string) error
...
...

promptTokens := tokenizer.Tokens(prompt)
completionTokens := tokenizer.Tokens(resp.Completion)


// evaluate total cost of the tokens based on anthropic pricing https://www-files.anthropic.com/production/images/model_pricing_july2023.pdf


```

## Status

This package was reverse engineered from [Anthropic's official tokenizer for typescript](https://github.com/anthropics/anthropic-tokenizer-typescript)

According to their README:

> This package is in beta. Its internals and interfaces are not stable
> and subject to change without a major semver bump;
> please reach out if you rely on any undocumented behavior.
> We are keen for your feedback; please email us at [support@anthropic.com](mailto:support@anthropic.com)

## Maintenence

Running `make config` will update the configuration based on whatever is in the Anthropic official tokenizer repo.
If Anthropic ever modifies their config file, running this command and running the command will update the config
