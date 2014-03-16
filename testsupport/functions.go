package testsupport

import (
	"fmt"
)

func Compare(exp, act interface{}) error {
	if exp != act {
		return fmt.Errorf("Expected %v; Got %v", exp, act)
	}
	return nil
}
