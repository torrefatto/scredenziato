package helpers

import (
	"fmt"
	"testing"
)

func TestChainErrError(t *testing.T) {
	expected := `[
	miao: someone woff'd,
	bau: someone purrrr'd,
	ciao: someone said 'hi!'
]`

	errs := NewChainErr()

	errs.Add("miao", fmt.Errorf("someone woff'd"))
	errs.Add("bau", fmt.Errorf("someone purrrr'd"))
	errs.Add("ciao", fmt.Errorf("someone said 'hi!'"))

	if res := errs.Error(); expected != res {
		t.Fatalf("unexpected: %s", res)
	}
}
