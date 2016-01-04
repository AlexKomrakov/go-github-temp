package models

import (
	"testing"
)

func TestStore(t *testing.T) {
	token := Token{"test_user", "test_token"}
	token.Store()
	token.Delete()
}