package jsonnext

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSlice(t *testing.T) {
	fs := &flag.FlagSet{}
	ss := StringSlice(fs, "food", "food to eat")
	err := fs.Parse([]string{"-food", "banana", "-food", "pie"})
	require.NoError(t, err)
	require.Equal(t, &[]string{"banana", "pie"}, ss)

	// Test Get() method
	f := fs.Lookup("food")
	require.NotNil(t, f)
	require.Equal(t, *ss, f.Value.(flag.Getter).Get())
}

func TestVMVarMapGet(t *testing.T) {
	fs := &flag.FlagSet{}
	m := VMVarMap{}
	fs.Var(vmVarMapValue{m, NewExtStr}, "ext-str", "ext var string")
	err := fs.Parse([]string{"-ext-str", "foo=bar"})
	require.NoError(t, err)

	// Test Get() method
	f := fs.Lookup("ext-str")
	require.NotNil(t, f)
	require.Equal(t, m, f.Value.(flag.Getter).Get())
}
