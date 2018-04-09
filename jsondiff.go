package jsondiff

import "io"
import "encoding/json"
import "fmt"
import "reflect"

type Change string

const (
	Added     Change = "Added"
	Updated   Change = "Updated"
	Removed   Change = "Removed"
	Unchanged Change = "Unchanged"
)

func Decode(r io.Reader) (map[string]interface{}, error) {
	d := json.NewDecoder(r)
	m := make(map[string]interface{})

	err := d.Decode(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func KeysFromMap(m map[string]interface{}) []string {
	s := make([]string, 0, len(m))

	for k := range m {
		s = append(s, k)
	}

	return s
}

func minLength(l, r []string) int {
	if ll, rl := len(l), len(r); ll < rl {
		return ll
	} else {
		return rl
	}
}

func AlignKeys(lks, rks []string) []string {
	seen := make(map[string]struct{})
	ml := minLength(lks, rks)

	keys := make([]string, 0, ml)

	for _, key := range lks {
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}

	for _, key := range rks {
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}

	return keys
}

type Pair struct {
	Key    string      `json:"key"`
	Left   interface{} `json:"left"`
	Right  interface{} `json:"right"`
	Change Change      `json:"change,omitempty"`
}

func (p Pair) String() string {
	return fmt.Sprintf("{Key: %v Left: %v Right: %v Change: %v}", p.Key, p.Left, p.Right, p.Change)
}

type Tree struct {
	Root     *Pair   `json:"root"`
	Children []*Tree `json:"children,omitempty"`
}

func (t *Tree) String() string {
	res := fmt.Sprintf("%v\n", t.Root.String())

	for _, p := range t.Children {
		res = fmt.Sprintf("%v  %v", res, p)
	}

	return res
}

func Diff(left, right map[string]interface{}) *Tree {
	tree := new(Tree)

	diff(tree, left, right, "")

	return tree
}

func diff(tree *Tree, left, right map[string]interface{}, rootKey string) *Tree {
	lks := KeysFromMap(left)
	rks := KeysFromMap(right)
	aks := AlignKeys(lks, rks)

	tree.Root = &Pair{
		Key:   rootKey,
		Right: right,
		Left:  left,
	}

	for _, key := range aks {
		rv := right[key]
		lv := left[key]
		lt := reflect.TypeOf(lv)
		rt := reflect.TypeOf(rv)
		if rt == nil && lt == nil {
			tree.Children = append(tree.Children, &Tree{Root: &Pair{
				Key:    key,
				Change: Unchanged,
				Left:   nil,
				Right:  nil,
			}})
		} else if rt == nil {
			tree.Children = append(tree.Children, &Tree{Root: &Pair{
				Key:    key,
				Change: Removed,
				Left:   lv,
				Right:  nil,
			}})
		} else if lt == nil {
			tree.Children = append(tree.Children, &Tree{Root: &Pair{
				Key:    key,
				Change: Added,
				Left:   nil,
				Right:  rv,
			}})
		} else if lt != rt {
			tree.Children = append(tree.Children, &Tree{Root: &Pair{
				Key:    key,
				Change: Updated,
				Left:   lv,
				Right:  rv,
			}})
		} else if lt.Kind() == reflect.Map {
			tree.Children = append(tree.Children, diff(new(Tree), lv.(map[string]interface{}), rv.(map[string]interface{}), key))
		} else if lv != rv {
			tree.Children = append(tree.Children, &Tree{Root: &Pair{
				Key:    key,
				Change: Updated,
				Left:   lv,
				Right:  rv,
			}})
		} else {
			tree.Children = append(tree.Children, &Tree{Root: &Pair{
				Key:    key,
				Change: Unchanged,
				Left:   lv,
				Right:  rv,
			}})
		}

	}

	// tree change overall=

	return tree
}
