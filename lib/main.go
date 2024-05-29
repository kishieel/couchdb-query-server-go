package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	server := NewQueryServer()

	dispatch := map[string]func(args ...any){
		"ddoc": func(args ...any) {
			// shows - deprecated
			// lists - deprecated
			// updates
			// filters
			// views
			// validate_doc_update
			// rewrites

			// ["ddoc","new","_design/go-2",{"_id":"_design/go-2","_rev":"3-2f8240d2cdd7985df6381eca0c10d62b","filters":{"filter":"func Filter(doc couchgo.Document) { return false }"},"language":"go","views":{"view-number-one":{"map":"func Map(doc couchgo.Document) {  }"},"view-number-two":{"map":"func Map(doc couchgo.Document) { \n    if doc[\"type\"] == \"post\" {\n        couchgo.Emit(doc[\"_id\"], 1)\n    }\n}","reduce":"func Reduce(keys []any, values []any, rereduce bool) any {\n\tout := 0.0\n\n\tfor _, value := range values {\n\t\tout += value.(float64)\n\t}\n\n\treturn out\n}"}}}]
		},
		"reset": func(args ...any) {
			server.Reset()
		},
		"add_fun": func(args ...any) {
			source := args[0].(string)
			server.AddFun(source)
		},
		"add_lib": func(args ...any) {
			server.AddLib()
		},
		"map_doc": func(args ...any) {
			doc := args[0].(Document)
			server.MapDoc(doc)
		},
		"reduce": func(args ...any) {
			sources := make([]string, len(args[0].([]any)))
			for i, source := range args[0].([]any) {
				sources[i] = source.(string)
			}

			keyValues := args[1].([]any)
			keys := make([]any, len(keyValues))
			values := make([]any, len(keyValues))

			for i, kv := range keyValues {
				keys[i] = kv.([]any)[0]
				values[i] = kv.([]any)[1]
			}

			server.Reduce(sources, keys, values, false)
		},
		"rereduce": func(args ...any) {
			sources := make([]string, len(args[0].([]any)))
			for i, source := range args[0].([]any) {
				sources[i] = source.(string)
			}

			values := args[1].([]any)
			server.Reduce(sources, nil, values, true)
		},
	}

	for scanner.Scan() {
		var message []any

		if err := json.Unmarshal(scanner.Bytes(), &message); err != nil {
			Respond([]string{"error", "unnamed_error", err.Error()})
		}

		Log(message)

		command := message[0].(string)
		arguments := message[1:]

		if dispatch[command] != nil {
			dispatch[command](arguments...)
		} else {
			Respond([]string{"error", "unknown_command", fmt.Sprintf("unknown command '%s'", command)})
		}
	}
}
