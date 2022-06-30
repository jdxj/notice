package github

import (
	"testing"
)

func TestGithub_T(t *testing.T) {
	g := NewGithub()
	g.getRepos()
	g.run()
}
