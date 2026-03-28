package cmd

import "testing"

func TestRepoNameFromURL(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://github.com/user/my-app.git", "my-app"},
		{"https://github.com/user/my-app", "my-app"},
		{"git@github.com:user/my-app.git", "my-app"},
		{"git@github.com:user/my-app", "my-app"},
		{"ssh://git@github.com/user/my-app.git", "my-app"},
	}
	for _, c := range cases {
		got, err := repoNameFromURL(c.input)
		if err != nil {
			t.Errorf("repoNameFromURL(%q): unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("repoNameFromURL(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestRepoNameFromURL_Invalid(t *testing.T) {
	_, err := repoNameFromURL("https://github.com/user/")
	if err == nil {
		t.Fatal("expected error for trailing slash URL, got nil")
	}
}
