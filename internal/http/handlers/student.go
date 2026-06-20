package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/hafizul16103123/student-api/internal/storage"
	"github.com/hafizul16103123/student-api/internal/types"
	"github.com/hafizul16103123/student-api/internal/utils/response"
)

func New(storage storage.IStorage) http.HandlerFunc {
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

		// create student
		id,err:=storage.CreateStudent(student.Name,student.Email,student.Age)
		if err!=nil {
			response.WriteJson(w,http.StatusCreated,err)
			return
		}
		slog.Info("student created successfully")
		response.WriteJson(w,http.StatusCreated,map[string]int64{"id":id})
	}
}

func GetById(storage storage.IStorage) http.HandlerFunc {
	return func(w http.ResponseWriter,r *http.Request){
		id:=r.PathValue("id")
		slog.Info("getting a student",slog.String("id",id))

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		student,err:=storage.GetStudentById(intId)
		if err!=nil {
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		response.WriteJson(w,http.StatusOK,student)

	}
}
