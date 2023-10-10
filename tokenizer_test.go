package tokenizer

import "testing"

func TestStuff(t *testing.T) {

	t.Run("small_text", func(t *testing.T) {
		if tokens := countTokens("hello world!"); tokens != 3 {
			t.Errorf("expected 3 got %d", tokens)
		}
	})

	t.Run("small_normalize", func(t *testing.T) {
		if tokens := countTokens("™"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}

		if tokens := countTokens("ϰ"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}

	})

	t.Run("allows_special_tokens", func(t *testing.T) {
		if tokens := countTokens("™"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}

		if tokens := countTokens("<EOT>"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}
	})
}
