package jsondiff

import "testing"
import "os"
import "reflect"
import "sort"
import "strings"

func TestDecode(t *testing.T) {
	r := strings.NewReader(`{ "foo": "bar" }`)
	e := map[string]interface{}{
		"foo": "bar",
	}

	g, err := Decode(r)

	if err != nil {
		t.Error(err)
	}

	eq := reflect.DeepEqual(g, e)

	if !eq {
		t.Errorf("Wrong result:\ngot: %v\nwant: %v\n", g, e)
	}
}

func TestAlignKeys(t *testing.T) {
	lks := []string{"foo", "bar"}
	rks := []string{"bar", "baz"}

	expected := []string{"foo", "bar", "baz"}
	got := AlignKeys(lks, rks)

	if l := len(got); l != 3 {
		t.Errorf("Wrong number of elements:\ngot: %v, want: %v\n", l, 3)
	}

	sort.Strings(got)
	sort.Strings(expected)

	eq := reflect.DeepEqual(got, expected)

	if !eq {
		t.Errorf("Wrong result:\ngot: %v\nwant: %v\n", got, expected)
	}
}

func TestKeysFromMap(t *testing.T) {
	m := map[string]interface{}{
		"foo": struct{}{},
		"bar": struct{}{},
		"baz": struct{}{},
	}

	expected := []string{"foo", "bar", "baz"}

	got := KeysFromMap(m)

	eq := reflect.DeepEqual(got, expected)

	if !eq {
		t.Errorf("Wrong result:\ngot: %v\nwant: %v\n", got, expected)
	}
}

func TestDiff(t *testing.T) {
	diffLeft, err := os.Open("testdata/diffleft.json")

	if err != nil {
		t.Error(err)
	}

	defer diffLeft.Close()

	diffRight, err := os.Open("testdata/diffright.json")

	if err != nil {
		t.Error(err)
	}

	defer diffRight.Close()

	diffResult, err := os.Open("testdata/diffresult.json")

	if err != nil {
		t.Error(err)
	}

	defer diffResult.Close()

	leftMap, err := Decode(diffLeft)

	if err != nil {
		t.Error(err)
	}

	rightMap, err := Decode(diffRight)

	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
	}

	tree := Diff(leftMap, rightMap)

	t.Errorf("\n\n%v\n\n", tree)
}
