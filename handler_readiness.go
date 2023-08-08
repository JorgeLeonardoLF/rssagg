package main

import "net/http"

/*
	This function signature is what you need to use if you want to define an http handler in the way the go standard library expects
	- a ResponseWriter as the first param and a pointer to an http.Request as the second

	- In the body of the function we can just call our respondWithJSON()
*/
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	/* We are passing to our responseWithJSON function that is defined in the json.go file
	- w - our http.ResponseWrite object
	- a status code of 200
	- struct {}{} and empty struct for the time being*/

	respondWithJSON(w, 200, struct{}{})
}
