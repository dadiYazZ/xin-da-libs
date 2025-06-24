package object

import (
	"github.com/dadiYazZ/xin-da-libs/fmt"
	"testing"
)

func Test_QuickRandom(t *testing.T) {

	for i := 1; i < 5; i++ {
		response := QuickRandom(4)
		fmt.Dump(response)
	}
}
