package club

import "testing"

var table = map[string]bool{
	"hjk321_jbnkl": true,
	"client1":      true,
	"client2":      true,
	"client3":      true,
	"client4":      true,
	"client 5":     false,
	"орлрофыв":     false,
	"@#$%^":        false,
	"CLIENT":       false,
}

func TestValidateClientName(t *testing.T) {
	for k, v := range table {
		if ValidateClientName(k) != v {
			t.Errorf("ValidateClientName(%s) failed", k)
		}
	}
}
