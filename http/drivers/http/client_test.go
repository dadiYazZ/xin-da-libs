package http

import (
	"github.com/dadiYazZ/xin-da-libs/http/contract"
	"testing"
)

func Test_NewClient(t *testing.T) {
	helper, err := NewHttpClient(&contract.ClientConfig{})
	if err != nil {
		t.Error(err)
	}

	if helper == nil {
		t.Error(err)
	}

}
