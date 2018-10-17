package misc

import "fmt"

func mapKvMustNotBeEmpty(m map[string]string) {
	for k, v := range m {
		if k != "" && v != "" {
			continue
		}
		panic(fmt.Sprintf("empty key-value: `%s = %s`", k, v))
	}
}
