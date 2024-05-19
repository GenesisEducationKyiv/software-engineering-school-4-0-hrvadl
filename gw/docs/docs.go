// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/rate": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rate"
                ],
                "summary": "Get USD -\u003e UAH exchange rate",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_hrvadl_converter_gw_internal_transport_http_handlers.Response-float32"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_hrvadl_converter_gw_internal_transport_http_handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/subscribe": {
            "post": {
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rate"
                ],
                "summary": "Subscribe to email rate exchange notification",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email to subscribe",
                        "name": "body",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_hrvadl_converter_gw_internal_transport_http_handlers.EmptyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_hrvadl_converter_gw_internal_transport_http_handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_hrvadl_converter_gw_internal_transport_http_handlers.EmptyResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "github_com_hrvadl_converter_gw_internal_transport_http_handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "github_com_hrvadl_converter_gw_internal_transport_http_handlers.Response-float32": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "number"
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}