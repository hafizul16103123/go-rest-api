package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/go-playground/validator/v10"
)

type Response struct{
	Status string `json:"status"`
	Error string `json:"error"`
}
const(
	StatusOK="OK"
	StatusError="Error"
)
func WriteJson(w http.ResponseWriter,status int, data any) error{
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(status)

	// Encode means converting Go data into a format that can be stored or transferred.
	// Go Struct -> json data
	// Encode: convert to External type
	// Decode: Deserialize (convert) back to base form
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response{
	return Response{
		Status:StatusError,
		Error: err.Error(),
	}
}

func ValidationErrors(errs validator.ValidationErrors) Response {

	var errMsgs []string

	for _,err:= range errs{
		switch err.ActualTag(){
		case "required":
			errMsgs=append(errMsgs,fmt.Sprintf("field %s is required", err.Field()))
		
		default:
			errMsgs=append(errMsgs,fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return Response{
		Status:StatusError,
		Error: strings.Join(errMsgs,", "),
	}

}