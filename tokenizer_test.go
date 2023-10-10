package tokenizer

import "testing"

func TestTokenCounter(t *testing.T) {
	tt, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("small_text", func(t *testing.T) {
		if tokens := tt.Tokens("hello world!"); tokens != 3 {
			t.Errorf("expected 3 got %d", tokens)
		}
	})

	t.Run("small_normalize", func(t *testing.T) {
		if tokens := tt.Tokens("™"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}

		if tokens := tt.Tokens("ϰ"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}
	})

	t.Run("allows_special_tokens", func(t *testing.T) {
		if tokens := tt.Tokens("™"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}

		if tokens := tt.Tokens("<EOT>"); tokens != 1 {
			t.Errorf("expected 1 got %d", tokens)
		}
	})
}
