package dbuser

import "testing"

func TestReadByID(t *testing.T) {
	obj, err := ReadByID(1)
	if err != nil {
		t.Error("couldn't get dbuser", err)
	}

	if obj.Username.String != "stealth" {
		t.Error("Username isn't correct")
	}
}

