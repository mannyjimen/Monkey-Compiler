package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("smybol defined incorrectly, expected %+v, got %+v", a, expected["a"])
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("smybol defined incorrectly, expected %+v, got %+v", b, expected["b"])
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("symbol %+v could not be resolved", sym)
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, instead got %+v",
				sym.Name, sym, result)
		}
	}
}
