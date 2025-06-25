package object

import (
	"github.com/dadiYazZ/xin-da-libs/xin-da-fmt"
	"testing"
)

func Test_QuickRandom(t *testing.T) {

	for i := 1; i < 5; i++ {
		response := QuickRandom(4)
		xin_da_fmt.Dump(response)
	}
}
