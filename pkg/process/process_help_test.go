package process

import (
	"fmt"
	"testing"
)

func TestHelp_Help(t *testing.T) {
	h := &Help{}
	fmt.Println(h.Help())
}
