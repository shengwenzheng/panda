package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"panda.com/api/pagination"
	"panda.com/api/service"
)

var (
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	contextType      = reflect.TypeOf((*gin.Context)(nil))
	paginationType   = reflect.TypeOf((*pagination.Query)(nil))
	paginationResult = reflect.TypeOf((*pagination.Result)(nil))
)

type handleFunc any

type Response struct {
	Code int 	   `json:"code"`
	Msg  string     `json:"msg"`
	Data any	   `json:"data,omitempty"` // omitempty表示如果Data为空则不返回该字段
}

func (s *Server) handle(fn handleFunc) gin.HandlerFunc {
	if err := validateFunc(fn); err != nil {
		slog.Error("invalid handler func", "error", err)
		panic(err)
	}

	return func(c *gin.Context) {
		ft := reflect.TypeOf(fn)
		args, err := buildParams(ft, c)
		if err != nil {
			c.Error(err)
			return
		}

		result := callHandleFunc(fn, args...)
		if err := result[len(result)-1]; err != nil {
			c.Error(err.(error))
			return
		}

		if ft.In(ft.NumIn()-1) == paginationType {
			r := result[0].(*pagination.Result)
			query := args[len(args)-1].(*pagination.Query)
			c.AbortWithStatusJSON(http.StatusOK, &pagination.Response{
				Code:   http.StatusOK,
				Result: r,
				Links:  pagination.GetLinks(c, r.Total, query),
			},
			)

			return
		}
		c.AbortWithStatusJSON(http.StatusOK, &Response{
			Code: http.StatusOK,
			Data: result[0],
		})
	}
}

func handleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last();
		if err != nil  {
			return 
		}

		slog.Error("handle error", "error", err.Err, "url", c.Request.URL)
		msg := err.Err.Error()
		code := getErrorCode(err.Err, service.ErrorCode);
		if code == -1 {
			msg = "request failed"
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, &Response{
			Code: code,
			Msg:  msg,
		})
	}
}

func getErrorCode(err error, errorCode map[error]int) int {
	if ok := isComparable(reflect.TypeOf(err)); ok {
		if errCode, ok := errorCode[err]; ok {
			return errCode
		}
	}
	return -1
}

func isComparable(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Map, reflect.Func,
		reflect.Slice:
		return false
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if !isComparable(t.Field(i).Type) {
				return false
			}
		}
	}
	return true
}

func buildParams(ft reflect.Type,ctx *gin.Context) ([]any, error) {
	args := []any{ctx}
	if ft.NumIn() == 1 {
		return args, nil
	}

	if ft.In(1) != paginationType {
		reqArg := reflect.New(ft.In(1)).Interface()
		if err := ctx.ShouldBindJSON(reqArg); err != nil {
			return nil, err
		}

		if err := validator.New().Struct(reqArg); err != nil {
			return nil, err
		}

		args = append(args, reqArg)
	}

	if ft.In(ft.NumIn()-1) == paginationType {
		query, err := pagination.Parse(ctx)
		if err != nil {
			return nil, err
		}

		args = append(args, query)
	}

	return args, nil
}

func callHandleFunc(fn handleFunc, args ...any) []any {
	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}

	rs := reflect.ValueOf(fn).Call(params)
	result := make([]any, len(rs))
	for i, r := range rs {
		result[i] = r.Interface()
	}
	return result
}

func validateFunc(fn handleFunc) error {
	ft := reflect.TypeOf(fn)
	if ft.Kind() != reflect.Func || ft.IsVariadic() {
		return fmt.Errorf("need non variadic func in %s", ft.String())
	}

	if ft.NumIn() < 1 || ft.NumIn() > 3 {
		return fmt.Errorf("the size of input parameters is not correct in %s", ft.String())
	}

	if ft.In(0) != contextType {
		return fmt.Errorf("the first parameter must point of context "+
			"in %s", ft.String())
	}

	if ft.NumIn() == 2 && ft.In(1).Kind() != reflect.Ptr {
		return fmt.Errorf("the second parameter must be a "+
			"pointer type in %s", ft.String())
	}

	if ft.NumOut() < 1 || ft.NumOut() > 2 {
		return fmt.Errorf("the number of return values must be "+
			"one or two in %s", ft.String())
	}

	if ft.In(ft.NumIn()-1) == paginationType && ft.Out(0) != paginationResult {
		return fmt.Errorf("the last of input parameter is "+
			"pagination Query type, the first return value must be "+
			"a paginationResult type in %s", ft.String())
	}

	if !ft.Out(ft.NumOut() - 1).Implements(errorType) {
		return fmt.Errorf("the last return value must be an "+
			"error type in %s", ft.String())
	}

	return nil
}