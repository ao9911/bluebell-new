package snowflake

import "testing"

func TestGenID(t *testing.T) {
	id1 := GenID()
	if id1 <= 0 {
		t.Fatalf("GenID() = %d, want positive", id1)
	}
	id2 := GenID()
	if id2 <= id1 {
		t.Fatalf("second GenID() = %d, want greater than first %d", id2, id1)
	}
}
