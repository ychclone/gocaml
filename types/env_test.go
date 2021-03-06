package types

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestDumpResult(t *testing.T) {
	env := NewEnv()
	env.Table["test_ident"] = IntType
	env.Table["test_ident2"] = BoolType
	env.Table["external_ident"] = UnitType
	env.Table["external_ident2"] = FloatType

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	env.Dump()

	ch := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		ch <- buf.String()
	}()
	w.Close()
	os.Stdout = old

	out := <-ch
	for _, s := range []string{
		"Variables:\n",
		"test_ident: int",
		"test_ident2: bool",
		"External Variables:\n",
		"external_ident: unit",
		"external_ident2: float",
	} {
		if !strings.Contains(out, s) {
			t.Fatalf("Output does not contain '%s': %s", s, out)
		}
	}
}

func TestEnvHasBuiltins(t *testing.T) {
	env := NewEnv()
	if len(env.Externals) == 0 {
		t.Fatal("Env must contain some external symbols by default because of builtin symbols")
	}
	if _, ok := env.Externals["print_int"]; !ok {
		t.Fatal("'print_int' is not found though it is builtin:", env.Externals)
	}
}
