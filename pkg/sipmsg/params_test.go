package sipmsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamsAdd(t *testing.T) {
	tests := map[string]struct {
		initial  Params
		key, val string
		want     Params
	}{
		`key with value to an empty initial value`: {
			"", "user", "phone", "user=phone"},
		`key without a value to an empty initial value`: {
			"", "lr", "", "lr"},
		`key with value to an existing params`: {
			"transport=UDP", "user", "phone", "transport=UDP;user=phone"},
		`key without value to an existing params`: {
			"transport=UDP", "lr", "", "transport=UDP;lr"},
		`unchanged when key is empty`: {
			"transport=UDP", "", "foo", "transport=UDP"},
		`unchanged when key and name are empty`: {
			"transport=UDP", "", "", "transport=UDP"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			prms := tt.initial
			assert.Equal(t, tt.want, prms.Add(tt.key, tt.val))
		})
	}
}

func TestParamsSet(t *testing.T) {
	tests := map[string]struct {
		initial  Params
		key, val string
		want     Params
	}{
		`key with value when empty initial value`: {
			"", "user", "phone", ""},
		`key with value when not match`: {
			"trans=tcp;user=ip;ttl=64", "maddr", "phone", "trans=tcp;user=ip;ttl=64"},
		`key with value when match`: {
			"trans=tcp;user=ip;ttl=64", "user", "phone", "trans=tcp;user=phone;ttl=64"},
		`match the first param`: {
			"trans=tcp;user=ip;ttl=64", "trans", "sctp", "trans=sctp;user=ip;ttl=64"},
		`match the last param`: {
			"trans=tcp;user=ip;ttl=64", "ttl", "255", "trans=tcp;user=ip;ttl=255"},
		`match when single param`: {
			"transport=tcp", "transport", "udp", "transport=udp"},
		`unchanged when key is empty`: {
			"transport=tcp;ttl=123", "", "udp", "transport=tcp;ttl=123"},
		`match without a value`: {
			"transport=tcp;group=1;ttl=123", "group", "", "transport=tcp;group;ttl=123"},
		`match single name and add a  value`: {
			"transport=tcp;group;ttl=123", "group", "12", "transport=tcp;group=12;ttl=123"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			prms := tt.initial
			assert.Equal(t, tt.want, prms.Set(tt.key, tt.val))
		})
	}
}

func TestParamsDel(t *testing.T) {
	tests := map[string]struct {
		initial Params
		key     string
		want    Params
	}{
		`empty params`: {
			"", "user", ""},
		`not exists in params`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "transport",
			"ttl=255;user=ip;lr;maddr=10.0.0.1",
		},
		`param in the middle`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "user",
			"ttl=255;lr;maddr=10.0.0.1",
		},
		`param in the beginning`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "ttl",
			"user=ip;lr;maddr=10.0.0.1",
		},
		`param in the end`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "maddr",
			"ttl=255;user=ip;lr",
		},
		`param without value`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "lr",
			"ttl=255;user=ip;maddr=10.0.0.1",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			prms := tt.initial
			assert.Equal(t, tt.want, prms.Del(tt.key))
		})
	}
}

func TestParamsGet(t *testing.T) {
	tests := map[string]struct {
		initial Params
		key     string
		wantVal string
		wantOk  bool
	}{
		`empty params`: {
			"", "user", "", false},
		`not exists in params`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "transport", "", false},
		`found in the middle`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "user", "ip", true},
		`found in the beginning`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "ttl", "255", true},
		`param without value`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "lr", "", true},
		`param without value case insensetive`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "Lr", "", true},
		`case insensitive with value`: {
			"ttl=255;user=ip;lr;maddr=10.0.0.1", "TTL", "255", true},
		`match param with empty value`: {
			"ttl=255;user=;lr;maddr=10.0.0.1", "User", "", true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			prms := tt.initial
			val, ok := prms.Get(tt.key)
			assert.Equal(t, tt.wantVal, val)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestParamsSetup(t *testing.T) {
	tests := map[string]struct {
		input, want string
	}{
		`empty string`:          {"", ""},
		`spaces only`:           {"   ", ""},
		`no left trim`:          {"user=ip;ttl=225", "user=ip;ttl=225"},
		`trim left space`:       {" user=ip;ttl=225", "user=ip;ttl=225"},
		`trim left spaces`:      {"     user=ip;ttl=225", "user=ip;ttl=225"},
		`trim left semi`:        {";user=ip;ttl=225", "user=ip;ttl=225"},
		`trim left semis`:       {";;;;user=ip;ttl=225", "user=ip;ttl=225"},
		`trim space and semi`:   {" ;user=ip;ttl=225", "user=ip;ttl=225"},
		`trim semi and space`:   {"; user=ip;ttl=225", "user=ip;ttl=225"},
		`trim semis and spaces`: {"   ;   user=ip;ttl=225", "user=ip;ttl=225"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, Params(tt.want), Params(tt.input).setup())
		})
	}
}

func BenchmarkParamsSet(b *testing.B) {
	prms := Params("transport=tcp;group=12;ttl=123")
	for i := 0; i < b.N; i++ {
		_ = prms.Set("group", "main")
	}
}
