package jsondiff

import "io"
import "encoding/json"

func Decode(r io.Reader) (map[string]interface{}, error) {
	d := json.NewDecoder(r)
	m := make(map[string]interface{})

	err := d.Decode(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func AlignKeys(leftKeys, rightKeys []string, ch chan string) {
	allKeysCount := len(leftKeys) + len(rightKeys)
	combinedKeys := make([]string, 0, allKeysCount)

	combinedKeys = append(combinedKeys, leftKeys...)
	combinedKeys = append(combinedKeys, rightKeys...)

	seen := make(map[string]struct{})

	for _, key := range combinedKeys {
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			ch <- key
		}
	}

	close(ch)
}
