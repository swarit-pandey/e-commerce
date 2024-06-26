// Package httpapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package httpapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new user
	// (POST /users)
	PostUsers(c *gin.Context)
	// User login
	// (POST /users/login)
	PostUsersLogin(c *gin.Context)
	// Initiate password reset
	// (POST /users/password-reset)
	PostUsersPasswordReset(c *gin.Context)
	// Update password
	// (PUT /users/password-reset)
	PutUsersPasswordReset(c *gin.Context)
	// Get user profile
	// (GET /users/{userId})
	GetUsersUserId(c *gin.Context, userId int)
	// Add user address
	// (POST /users/{userId})
	PostUsersUserId(c *gin.Context, userId int)
	// Delete user address
	// (DELETE /users/{userId}/addresses/{addressId})
	DeleteUsersUserIdAddressesAddressId(c *gin.Context, userId int, addressId int)
	// Update user address
	// (PUT /users/{userId}/addresses/{addressId})
	PutUsersUserIdAddressesAddressId(c *gin.Context, userId int, addressId int)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// PostUsers operation middleware
func (siw *ServerInterfaceWrapper) PostUsers(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUsers(c)
}

// PostUsersLogin operation middleware
func (siw *ServerInterfaceWrapper) PostUsersLogin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUsersLogin(c)
}

// PostUsersPasswordReset operation middleware
func (siw *ServerInterfaceWrapper) PostUsersPasswordReset(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUsersPasswordReset(c)
}

// PutUsersPasswordReset operation middleware
func (siw *ServerInterfaceWrapper) PutUsersPasswordReset(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutUsersPasswordReset(c)
}

// GetUsersUserId operation middleware
func (siw *ServerInterfaceWrapper) GetUsersUserId(c *gin.Context) {

	var err error

	// ------------- Path parameter "userId" -------------
	var userId int

	err = runtime.BindStyledParameterWithOptions("simple", "userId", c.Param("userId"), &userId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter userId: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUsersUserId(c, userId)
}

// PostUsersUserId operation middleware
func (siw *ServerInterfaceWrapper) PostUsersUserId(c *gin.Context) {

	var err error

	// ------------- Path parameter "userId" -------------
	var userId int

	err = runtime.BindStyledParameterWithOptions("simple", "userId", c.Param("userId"), &userId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter userId: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUsersUserId(c, userId)
}

// DeleteUsersUserIdAddressesAddressId operation middleware
func (siw *ServerInterfaceWrapper) DeleteUsersUserIdAddressesAddressId(c *gin.Context) {

	var err error

	// ------------- Path parameter "userId" -------------
	var userId int

	err = runtime.BindStyledParameterWithOptions("simple", "userId", c.Param("userId"), &userId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter userId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "addressId" -------------
	var addressId int

	err = runtime.BindStyledParameterWithOptions("simple", "addressId", c.Param("addressId"), &addressId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter addressId: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteUsersUserIdAddressesAddressId(c, userId, addressId)
}

// PutUsersUserIdAddressesAddressId operation middleware
func (siw *ServerInterfaceWrapper) PutUsersUserIdAddressesAddressId(c *gin.Context) {

	var err error

	// ------------- Path parameter "userId" -------------
	var userId int

	err = runtime.BindStyledParameterWithOptions("simple", "userId", c.Param("userId"), &userId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter userId: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "addressId" -------------
	var addressId int

	err = runtime.BindStyledParameterWithOptions("simple", "addressId", c.Param("addressId"), &addressId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter addressId: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutUsersUserIdAddressesAddressId(c, userId, addressId)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/users", wrapper.PostUsers)
	router.POST(options.BaseURL+"/users/login", wrapper.PostUsersLogin)
	router.POST(options.BaseURL+"/users/password-reset", wrapper.PostUsersPasswordReset)
	router.PUT(options.BaseURL+"/users/password-reset", wrapper.PutUsersPasswordReset)
	router.GET(options.BaseURL+"/users/:userId", wrapper.GetUsersUserId)
	router.POST(options.BaseURL+"/users/:userId", wrapper.PostUsersUserId)
	router.DELETE(options.BaseURL+"/users/:userId/addresses/:addressId", wrapper.DeleteUsersUserIdAddressesAddressId)
	router.PUT(options.BaseURL+"/users/:userId/addresses/:addressId", wrapper.PutUsersUserIdAddressesAddressId)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xXQW/zNgz9K4K2o/c56ToU8y3ttiJYsQXdAhQoctAsJlFnS65EN+gC//dBkp3ajZw4",
	"W9PDd4ojUST13iNpb2mq8kJJkGhosqUmXUPO3OOEcw3GPRZaFaBRgPvH/MadkDC2//G1AJpQg1rIFa2i",
	"tsFF0CAV+BreUBx6NkqJOnxI8NaykAgr0HbdIMOQtypqVtRfT5CitZ0xYzZK83swgPfwXILB/ZtDzkQW",
	"9qjhuRQaOE0ea7PFgTDzgjOE3jgSNo1l8Mao/gZ5PA9vFtGi8RVKaW5A36mVkL3ZFIdSKQ1oyXI4ns3O",
	"cnhCplDSQECBaQrG/NmPQtDpTKulyGA4rRHtudmxawfD38NKGNQMhZInK+xAKh9Ej4/cctc6vBh0pT62",
	"+u/UV7knoWuXhFwqa8zBpFoUNh+auBTJZDY1JAM0xHol2qUMOiKZ1Zj7USUSJjl5EbAhuAahSVFrJaIo",
	"0GrGO+MqZ0I6nzSiL6CNjzT+MrJ5qwIkKwRN6Pdf7JJFE9cOhNhG9+WkPPEWIofclNOEzpTBuTPxtIDB",
	"a8Vdw0uVRJDuDCuKTKTuVPxklHxr2fbpWw1LmtBv4reeHtcNPe5TYFV5IXjyXIIXo/EZw9YqcWG7dN1o",
	"YAjcAnk5Gu3Tec042aVtbX7ct7lRcpkJr9AfHh72Df4A/QKa/Ky10k5NpsxzZidLHZ8wImHjxOL2PXGx",
	"U8sA+lznOiOHnVYdJG90jnj9rP3+qydjHKg/yUpcKy3+8ayeSogrOQ98i4mmQX2n7aweQElntp+JmuD7",
	"Q5Cey30IflPkpk5ikPY7IE2lQGF12+BCPC52MpQhWMrPR6X7unMuWAapsCswl9cOubbKtvZnyivrbwUB",
	"HG/B4zh3dq7Va5YDui7/uKW2W7j2T5vZTcvG9G3uoi4haqH5fhZWizPXd/NSdKCyg9Qg+UWV8j+V9S2g",
	"n8XNkLVSPVzCnw5yQEcnj6fxx46nCecetvrzKqTWuN4DE2/rx1rDHDLwX0NdiH9y6y2QJ42HSXP+bLBH",
	"QU+sFfd/cTigpRyWdgd+D9Q7Bo702K8A0Y8fCpOdfPteXsJtaEjJncBn3fvfVZQ1cSXoaSl1RhOaqZRl",
	"a2UwuRpdXdJqUf0bAAD//3eLatI3EQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
