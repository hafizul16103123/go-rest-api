package student

import (
	"encoding/json"
	"net/http"
	"errors"
	"io"
	"fmt"
	"log/slog"
	"github.com/hafizul16103123/student-api/internal/types"
	"github.com/hafizul16103123/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")
		var student types.Student

		// Decode means converting external data back into Go objects. JSON Data -> Go Struct
		// Encode: convert to External type
		// Decode: Deserialize (convert) back to base form
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err,io.EOF) {

			// custom error message:
			response.WriteJson(w, http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
		
			return
		}
		if err!=nil{
			response.WriteJson(w, http.StatusBadRequest,response.GeneralError(err))
			return
		}
		//req validate
		if err:=validator.New().Struct(student); err!=nil{
			validateErrors:=err.(validator.ValidationErrors)// type casting to validate error
			response.WriteJson(w, http.StatusBadRequest,response.ValidationErrors(validateErrors))
			return
		}

		response.WriteJson(w,http.StatusCreated,map[string]string{"success":"OK"})
	}
}