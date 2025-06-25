package object

import (
	"github.com/dadiYazZ/xin-da-libs/xin-da-fmt"
	"testing"
)

func Test_Attribute_Set_AND_Get(t *testing.T) {

	mapTest := NewAttribute(&HashMap{
		"gun": "model",
	})

	mapTest.SetAttribute("weapon.bullet", 100)
	mapTest.SetAttribute("weapon.shield.strength", "strong")

	bulletCount := mapTest.GetAttribute("weapon.bullet", 0)
	if bulletCount != 100 {
		t.Error("get bullet error")
		xin_da_fmt.Dump(bulletCount)
	}

	shieldStrength := mapTest.Get("weapon.shield.strength", "")
	if shieldStrength != "strong" {
		t.Error("get shield error")
		xin_da_fmt.Dump(shieldStrength)
	}

}
