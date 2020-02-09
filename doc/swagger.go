// Package mserv Mserv API.
//
// Provides access to operations over an Mserv service.
//
//     Schemes: http, https
//     BasePath: /
//     Host: localhost:8989
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key
//
//    SecurityDefinitions:
//    api_key:
//      type: apiKey
//      name: X-Api-Key
//      in: header
//
// swagger:meta
package doc

import (
	coprocess "github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/health"
	"github.com/TykTechnologies/mserv/models"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/go-openapi/runtime"
)

// swagger:response genericErrorResponse
type GenericErrorResponse struct {
	// in: body
	Body models.Payload
}

// swagger:response healthResponse
type HealthResponse struct {
	// in: body
	Body struct {
		models.BasePayload
		Payload health.HReport
	}
}

// swagger:response invocationResponse
type InvocationResponse struct {
	// in: body
	Body struct {
		models.BasePayload
		Payload coprocess.Object
	}
}

// swagger:response mwIDResponse
type MiddlewareIDResponse struct {
	// in: body
	Body struct {
		models.BasePayload
		Payload struct {
			BundleID string
		}
	}
}

// swagger:response mwResponse
type MiddlewareResponse struct {
	// in: body
	Body struct {
		models.BasePayload
		Payload storage.MW
	}
}

// swagger:response mwListResponse
type MiddlewareListResponse struct {
	// in: body
	Body struct {
		models.BasePayload
		Payload []storage.MW
	}
}

// swagger:response mwBundleResponse
type MiddlewareBundleResponse struct {
	// in: body
	File runtime.File
}

// swagger:parameters mwDelete mwFetch mwFetchBundle
type GenericMiddlewareParameters struct {
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters invoke
type InvocationParameters struct {
	// in: path
	// required: true
	Name string `json:"name"`
	// in: body
	// required: true
	Body coprocess.Object
}

// swagger:parameters mwAdd
type AddMiddlewareParameters struct {
	// in: formData
	// required: true
	// swagger:file
	UploadFile runtime.File `json:"uploadfile"`
	// in: query
	StoreOnly bool `json:"store_only"`
	// in: query
	ApiID string `json:"api_id"`
}

// swagger:parameters mwUpdate
type UpdateMiddlewareParameters struct {
	GenericMiddlewareParameters
	// in: formData
	// required: true
	// swagger:file
	UploadFile runtime.File `json:"uploadfile"`
}
