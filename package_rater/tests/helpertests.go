package tests

// import (
// 	"fmt"
// 	"group-24-ECE461/internal/helper"
// )

// // Function running the Test suite on helper package
// func RunHelperTests() {
// 	passed = 0    // initialize passed tests
// 	testCount = 0 // initialize all tests
// 	fmt.Println("Running Tests on Helpers...")
// 	TestBase64Decode()                                                                                      // run test on Base64Decode function
// 	TestGetLastWeek()                                                                                       // run test on GetLastWeek function
// 	fmt.Printf("Passed %d of %d tests. Completion: %%/%d\n\n\n", passed, testCount, (passed/testCount)*100) // Print test results
// }

// func TestBase64Decode() { // test Base63Decode function in the helper package
// 	testCount++                        // increment overall test count
// 	str := helper.Base64Decode("test") // call the function
// 	val := interface{}(str)            // cast to interface for type comparison
// 	if _, ok := val.(string); !ok {    // compare type
// 		fmt.Println("Error decoding base 64") // error message
// 	} else {
// 		fmt.Println("Success decoding base 64") // success message
// 		passed++                                // increment passed test count
// 	}
// }

// func TestGetFiveContributions() { // test GetFiveContributors function in the helper package
// }

// func TestCountCommits() { // test CountCommits function in the helper package

// }

// func TestGetLastWeek() { // test GetLastWeek function in the helper package
// 	testCount++                     // increment overall test count
// 	date := helper.GetLastWeek()    // call the function
// 	val := interface{}(date)        // cast to interface for type comparison
// 	if _, ok := val.(string); !ok { // compare type
// 		fmt.Println("Error getting last week's date") // error message
// 	} else {
// 		fmt.Println("Success getting last week's date") // success message
// 		passed++                                        // increment passed test count
// 	}
// }

// func TestPrintRepo() { // test PrintRepo function in the helper package

// }
