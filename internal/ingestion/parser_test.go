package ingestion

import "testing"

func TestParseMovieRejectsMalformedRow(t *testing.T) {
	row := []string{"False", "", "0", "[]", "", "862", "tt0114709", "en", "Toy Story", "ok", "21.9", "/poster.jpg", "[]", "[]", "1995-10-30", "0", "81", "[]", "Released", "", "Toy Story", "False", "7.7"}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ParseMovie panicked on a malformed row: %v", r)
		}
	}()

	_, err := ParseMovie(row)
	if err == nil {
		t.Fatal("expected malformed row to return an error")
	}
}
