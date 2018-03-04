package database

import (
	"encoding/json"
	"testing"
)

func TestIntBool_Scan(t *testing.T) {
	var tests = []struct {
		in    interface{}
		out   bool
		valid bool
	}{
		{int64(1), true, true},
		{int64(0), false, true},
		{nil, false, false},
	}

	for _, test := range tests {
		ib := IntBool{}
		ib.Scan(test.in)
		if test.out != ib.Bool {
			t.Errorf("scan failed expected %t got %t", test.out, ib.Bool)
		}

		if test.valid != ib.Valid {
			t.Errorf("valid failed expected %t got %t", test.valid, ib.Valid)
		}
	}

	ib := IntBool{}
	err := ib.Scan("fail")
	if err == nil {
		t.Error("scan accepted invalid value")
	}
}

func TestIntBool_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		IntBool
	}{
		{IntBool{
			Valid: true,
			Bool:  true,
		}},
		{
			IntBool{
				Valid: false,
				Bool:  false,
			},
		},
	}

	for _, ib := range tests {

		b, err := json.Marshal(ib)
		if err != nil {
			t.Error(err)
		}

		var ibUnmarshal IntBool
		err = json.Unmarshal(b, &ibUnmarshal)
		if err != nil {
			t.Error(err)
		}

		if ib.Bool != ibUnmarshal.Bool {
			t.Errorf("unmarshal failed bools aren't same expected %v got %v", ib.Bool, ibUnmarshal.Bool)
		}

		if ib.Valid != ibUnmarshal.Valid {
			t.Errorf("unmarshal failed valids aren't same expected %v got %v", ib.Valid, ibUnmarshal.Valid)
		}
	}

	var ib IntBool
	err := json.Unmarshal([]byte("1"), &ib)
	if err != nil {
		t.Error(err)
	}

	if ib.Bool != true {
		t.Error("bool not true on 1")
	}
}
