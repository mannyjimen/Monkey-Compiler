package object

import "testing"

func TestStringHashKeys(t *testing.T) {
	key1 := &String{Value: "foo"}
	key2 := &String{Value: "foo"}

	alt1 := &String{Value: "bar"}
	alt2 := &String{Value: "bar"}

	if key1.HashKey() != key2.HashKey() {
		t.Errorf("Keys with same value result in different hash key")
	}

	if alt1.HashKey() != alt2.HashKey() {
		t.Errorf("Keys with same value result in different hash key")
	}

	if key1.HashKey() == alt1.HashKey() {
		t.Errorf("Keys with different values result in same hash key.")
	}
}

func TestAlternateTypeHashKeys(t *testing.T) {
	key1 := &Integer{Value: 5}
	key2 := &Integer{Value: 5}

	alt1 := &Boolean{Value: true}
	alt2 := &Boolean{Value: true}

	if key1.HashKey() != key2.HashKey() {
		t.Errorf("Keys with same value result in different hash key")
	}

	if alt1.HashKey() != alt2.HashKey() {
		t.Errorf("Keys with same value result in different hash key")
	}

	if key1.HashKey() == alt1.HashKey() {
		t.Errorf("Keys with different values result in same hash key.")
	}
}
