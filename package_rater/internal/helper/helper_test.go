package helper

import (
	"testing"
)

// Test suite on helper package

func TestBase64Decode(t *testing.T) { // test Base63Decode function in the helper package
	str := Base64Decode("test")     // call the function
	val := interface{}(str)         // cast to interface for type comparison
	if _, ok := val.(string); !ok { // compare type
		t.Error("Error decoding base 64")
	}
}

func TestGetLastWeek(t *testing.T) { // test GetLastWeek function in the helper package
	date := GetLastWeek()           // call the function
	val := interface{}(date)        // cast to interface for type comparison
	if _, ok := val.(string); !ok { // compare type
		t.Error("Error getting last week's date") // error message
	}
}

func TestGetPackageName(t *testing.T) {
	url := "https://npmjs.com/package/node"
	name := GetPackageName(url)
	if name != "node" {
		t.Error("Error getting NPM package name from url")
	}
}

func TestGetOwnerAndName(t *testing.T) {
	url := "https://github.com/nodejs/node"
	owner, name := GetOwnerAndName(url)
	if name != "node" || owner != "nodejs" {
		t.Error("Error getting Github package name and owner from url")
	}
}
