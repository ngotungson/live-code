package api

import (
	"context"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
)

func yamlBodyDecoder(r io.Reader, h http.Header, s *openapi3.SchemaRef, fn openapi3filter.EncodingFn) (interface{}, error) {
	return "", nil
}

func OpenapiInputValidator(openapiFile string) gin.HandlerFunc {
	router := openapi3filter.NewRouter().WithSwaggerFromFile(openapiFile)
	openapi3filter.RegisterBodyDecoder("text/yaml", yamlBodyDecoder)

	return func(c *gin.Context) {
		// before request
		httpReq := c.Request
		ctx := context.Background()

		//context.TODO()SO

		var requestValidationInput *openapi3filter.RequestValidationInput
		// Find route
		route, pathParams, _ := router.FindRoute(httpReq.Method, httpReq.URL)

		// Validate request
		requestValidationInput = &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		}

		c.Next()
	}
}
