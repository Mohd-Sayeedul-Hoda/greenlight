package data

import(
	"fmt"
	"strconv"
	"errors"
	"strings"
)

var ErrInvalidRuntimeFormate = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error){
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte)error{
	unquoteJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil{
		return ErrInvalidRuntimeFormate
	}
	
	parts := strings.Split(unquoteJSONValue, " ")
	// we are checking that runtime coming is runtime + string if not return error
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormate
	}

	//otherwise parse the string containing the number into int32
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil{
		return ErrInvalidRuntimeFormate
	}

	*r = Runtime(i)

	return nil
}
