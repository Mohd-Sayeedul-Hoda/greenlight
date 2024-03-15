package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"fmt"
	"io"
	"strings"

	"greenlight/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error){

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter") 
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error{
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil{
		return err
	}

	js = append(js, '\n')

	for key, value := range headers{
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error{

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)

	// if field is not json but given by client then 
	// then it will return error if we not insalize the methond then
	// it will ignore it not return error and it will stop reading json at 
	// that point
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)

	if err != nil{
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch{
		// this error happen because
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contain badly-formated JSON (at character %d)", syntaxError.Offset)

		// this error happen due to syntax error in json
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formated JSON")

		// this error occur when there is json have wrong type of schema that we have define or 
		// error is related specific field
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != ""{
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// this error happen because request body is empty
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// when client has but extra field that cannot be
		// map by json decoder then it will return error
		// starting with "json: unknown"
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return	fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		// this error occur when does not pass any thing to json decoder
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF{
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

// use to return a string if exisits in url parameter if not return some default value
func (app *application) readStirng(qs url.Values, key string, defaultValue string) string{
	
	// get the value of key from map of parameter
	s := qs.Get(key)

	if s == ""{
		return defaultValue
	}

	return s

}

// function use to read and return value of string in parameter like when searching 
// generes we have more then one to extract that we use it 
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string{
	
	csv := qs.Get(key)

	if csv == ""{
		return defaultValue
	}

	return strings.Split(csv, ",")
}


// so to read int for query parameter if not there return default value if can't convert
// return error
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator)int{

	s := qs.Get(key)

	if s == ""{
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil{
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i

}
