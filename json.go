package main //In the main package

import (
	"encoding/json"
	"log"
	"net/http"
) //Package for http, hover for details

/* {Note Stamp 1}
-TLDR: This function is our helper function to turn a payload (aka some sort of interface) data into a json object/string that is passed back through the w http.ResponseWriter instance.
	- a json object / string is used to write the data back to as an http response
	- We need to do error handling as well as success handling
- w http.ResponseWriter - same http response writer that http handlers in go use, exposed by the net/http package
- code int - the status code we wil responde with
- payload interface{} - an interface, something that we can Marshal (an acutal function) to a json structure
*/

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	/* - Marshal the payload into a json string
	- This function will take what is passed and return it as bytes
		- It is returned as bytes is so that we can write it in a binary format for the http response
	*/
	dat, err := json.Marshal(payload)
	if err != nil {
		/*
			If something goes wrong when trying to marshal the payload will will write to the response a header with a status code
			In this case 500 - Internal server error
			Log the issue
			return
		*/
		w.WriteHeader(500)
		log.Printf("Failed to marshal JSON response: %v", payload) // Note we are not fatal
		return
	}
	/*
		If no error on marshal
		Add a header to the response to say we are responding with json data
		"Content-Type", "application/json" --> Key/ value pair this is the standard for json responses
		TLDR; this lets it know what to expect
	*/
	w.Header().Add("Content-Type", "application/json")

	/*Now we write status code that was passed to us*/
	w.WriteHeader(code)
	/*Write the data it self*/
	w.Write(dat)
}

/*
{Note Stamp 2}
This function similar to respondWithJSON, will take a payload - this case msg- and turn it into json
this is specific for handling generic errors
*/
func respondWithError(w http.ResponseWriter, code int, msg string) {
	/*Errors in the 400 range are client side. We care about server side with are 500 and up*/
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	/*
		We are defining a struct here because we want to to have a specific format
			- the format is special because we are using a json reflect tag (`json:"error"`)
			1: these are not ' single quotes they are ` backtick
			2: `json:"error"` denotes a json reflect tag that instructs how to marshal and unmarshal the struct in the json object/string
				- In this case we are saying we have a field (Error) of type string, we want the Key for this field to be "error"
				-On the unmarshal side it would look something like
					{
						"error": "Something went wrong"
					}
	*/
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}
