package helpers

import (
	"fmt"
	"strings"
)

type ChainErr map[string]error

func NewChainErr() ChainErr {
	return make(map[string]error)
}

func (e ChainErr) Error() string {
	var result []string

	for helper, err := range e {
		result = append(result, fmt.Sprintf("\n\t%s: %s", helper, err))
	}

	if len(result) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s]", strings.Join(result, ",")+"\n")
}

func (e ChainErr) Add(helper string, err error) {
	e[helper] = err
}
