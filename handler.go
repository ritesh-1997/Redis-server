package main
import(
	"sync"
)
// var Handlers = map[string]func([]Value) Value
// Handlers for RESP commands

func ping(args []Value) Value {
	if len(args) == 0{
		return Value{typ: "string", str: "PONGing"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET": set,
	"GET": get,
	"HSET": hset,
	"HGET": hget,
	// Add more handlers for other commands here
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk
	// store the key-value pair in the SETs map
	// using a mutex to ensure thread safety
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()
	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk
	// retrieve the value from the SETs map
	// using a mutex to ensure thread safety
	SETsMu.RLock()
	value, exists := SETs[key]
	SETsMu.RUnlock()

	if !exists {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) < 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, exists := HSETs[hash]; !exists {
		HSETs[hash] = make(map[string]string)
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"} 
}

func hget(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, exists := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !exists {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}