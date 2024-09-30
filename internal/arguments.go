package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"reflect"
)

// SolidityArgument represents a single constructor argument with type and value.
type SolidityArgument struct {
	Name       string             `json:"name"`                 // Argument name in the constructor
	Type       string             `json:"type"`                 // Argument type in Solidity (e.g., uint256, string)
	Value      interface{}        `json:"value"`                // Value of the argument
	Components []SolidityArgument `json:"components,omitempty"` // Components for struct/tuple types
}

// ConstructorArguments holds a list of SolidityArgument, representing all constructor arguments.
type ConstructorArguments struct {
	Inputs []SolidityArgument `json:"constructorArguments"`
}

// parseSolidityArgument parses a SolidityArgument based on its type and prints it.
func parseSolidityArgument(arg SolidityArgument) (interface{}, error) {
	switch arg.Type {
	case "uint256", "int256", "uint", "int", "uint8", "uint16", "uint32", "uint64":
		return parseNumericArgument(arg)
	case "string":
		return fmt.Sprintf("%s", arg.Value), nil
	case "bool":
		return fmt.Sprintf("%v", arg.Value), nil
	case "address":
		return fmt.Sprintf("%s", arg.Value), nil
	case "bytes":
		return fmt.Sprintf("%s", arg.Value), nil // bytes are stored as a string in the JSON format (hex)
	case "uint256[]", "string[]", "bool[]", "address[]":
		return parseArrayArgument(arg), nil
	case "tuple":
		return parseTupleArgument(arg.Components), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", arg.Type)
	}
}

func parseNumericArgument(arg SolidityArgument) (*big.Int, error) {
	return big.NewInt(arg.Value.(int64)), nil // Signed types
}

// parseArrayArgument handles array types like `uint256[]` or `string[]`.
func parseArrayArgument(arg SolidityArgument) []interface{} {
	array := []interface{}{}
	switch reflect.TypeOf(arg.Value).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(arg.Value)
		for i := 0; i < s.Len(); i++ {
			array = append(array, s.Index(i).Interface())
		}
	}
	return array
}

// parseTupleArgument handles struct/tuple arguments with nested components.
func parseTupleArgument(components []SolidityArgument) map[string]interface{} {
	tuple := map[string]interface{}{}
	for _, component := range components {
		val, err := parseSolidityArgument(component)
		if err != nil {
			log.Fatalf("Failed to parse tuple argument: %v", err)
		}
		tuple[component.Name] = val
	}
	return tuple
}

// parseConstructorArguments takes the entire JSON and deserializes the constructor arguments.
func parseConstructorArguments(jsonData []byte) (ConstructorArguments, error) {
	var args ConstructorArguments
	err := json.Unmarshal(jsonData, &args)
	if err != nil {
		return args, err
	}
	return args, nil
}

// printConstructorArguments prints all constructor arguments in a readable format.
func printConstructorArguments(args ConstructorArguments) {
	fmt.Println("Parsed Constructor Arguments:")
	for _, arg := range args.Inputs {
		vv, err := parseSolidityArgument(arg)
		if err != nil {
			log.Fatalf("Failed to parse argument: %v", err)
		}
		fmt.Printf("Name: %s, Type: %s, Value: %v\n", arg.Name, arg.Type, vv)
	}
}

//func main() {
//	// Example JSON configuration for the Solidity constructor arguments
//	jsonConfig := `{
//		"constructorArguments": [
//			{
//				"name": "_num",
//				"type": "uint256",
//				"value": 42
//			},
//			{
//				"name": "_owner",
//				"type": "address",
//				"value": "0x1234567890abcdef1234567890abcdef12345678"
//			},
//			{
//				"name": "_isActive",
//				"type": "bool",
//				"value": true
//			},
//			{
//				"name": "_greeting",
//				"type": "string",
//				"value": "Hello, Ethereum!"
//			},
//			{
//				"name": "_data",
//				"type": "bytes",
//				"value": "0x68656c6c6f"
//			},
//			{
//				"name": "_scores",
//				"type": "uint256[]",
//				"value": [10, 20, 30, 40]
//			},
//			{
//				"name": "_names",
//				"type": "string[]",
//				"value": ["Alice", "Bob", "Charlie"]
//			},
//			{
//				"name": "_person",
//				"type": "tuple",
//				"components": [
//					{
//						"name": "name",
//						"type": "string",
//						"value": "Alice"
//					},
//					{
//						"name": "age",
//						"type": "uint256",
//						"value": 30
//					},
//					{
//						"name": "wallet",
//						"type": "address",
//						"value": "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef"
//					}
//				]
//			}
//		]
//	}`
//
//	// Parse the JSON configuration
//	args, err := parseConstructorArguments([]byte(jsonConfig))
//	if err != nil {
//		log.Fatalf("Failed to parse constructor arguments: %v", err)
//	}
//
//	// Print the parsed arguments
//	printConstructorArguments(args)
//}
