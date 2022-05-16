// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/group/list": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "get groups I joined",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/group.GroupInfoResp"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "test connection",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api"
                ],
                "summary": "ping",
                "responses": {
                    "200": {
                        "description": "'pong'",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/token/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token"
                ],
                "summary": "login",
                "parameters": [
                    {
                        "description": "login information",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.AuthReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.AuthResp"
                        }
                    }
                }
            }
        },
        "/token/logout": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token"
                ],
                "summary": "logout",
                "responses": {
                    "200": {
                        "description": "'logout'",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/token/refresh": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token"
                ],
                "summary": "refresh token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.AuthResp"
                        }
                    }
                }
            }
        },
        "/user/address": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "add address",
                "parameters": [
                    {
                        "description": "new address information",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/address.AddressCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/address.AddressCreateResp"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "delete address",
                "parameters": [
                    {
                        "description": "address information",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/address.AddressDeleteReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "'success'",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/me": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "get user info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.UserInfoResp"
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "modify user info",
                "parameters": [
                    {
                        "description": "user's new information",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.UserModifyReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "'success'",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "register",
                "parameters": [
                    {
                        "description": "register information",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.UserCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.UserCreateResp"
                        }
                    }
                }
            }
        },
        "/user/workinfo": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "dashboard",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.UserDashboardResp"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api"
                ],
                "summary": "get API version",
                "responses": {
                    "200": {
                        "description": "version",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "address.AddressCreateReq": {
            "type": "object",
            "required": [
                "area",
                "city",
                "detail",
                "is_default",
                "province"
            ],
            "properties": {
                "area": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "detail": {
                    "type": "string"
                },
                "is_default": {
                    "type": "boolean"
                },
                "province": {
                    "type": "string"
                }
            }
        },
        "address.AddressCreateResp": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "address.AddressDeleteReq": {
            "type": "object",
            "required": [
                "address_id"
            ],
            "properties": {
                "address_id": {
                    "type": "integer"
                }
            }
        },
        "group.GroupInfoResp": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                }
            }
        },
        "user.AuthReq": {
            "type": "object",
            "required": [
                "account",
                "password",
                "type"
            ],
            "properties": {
                "account": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "type": {
                    "type": "string",
                    "enum": [
                        "name",
                        "phone"
                    ],
                    "example": "name"
                }
            }
        },
        "user.AuthResp": {
            "type": "object",
            "properties": {
                "login_type": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                },
                "user_phone": {
                    "type": "string"
                },
                "user_role": {
                    "type": "integer"
                },
                "user_token": {
                    "type": "string"
                }
            }
        },
        "user.UserCreateReq": {
            "type": "object",
            "required": [
                "user_address",
                "user_name",
                "user_phone",
                "user_role",
                "user_secret"
            ],
            "properties": {
                "user_address": {
                    "$ref": "#/definitions/user.UserCreateReqAddress"
                },
                "user_name": {
                    "type": "string",
                    "maxLength": 30
                },
                "user_phone": {
                    "type": "string",
                    "maxLength": 20,
                    "example": "13800138000"
                },
                "user_role": {
                    "type": "integer",
                    "example": 1
                },
                "user_secret": {
                    "type": "string",
                    "maxLength": 20
                }
            }
        },
        "user.UserCreateReqAddress": {
            "type": "object",
            "required": [
                "area",
                "city",
                "detail",
                "province"
            ],
            "properties": {
                "area": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "detail": {
                    "type": "string"
                },
                "province": {
                    "type": "string"
                }
            }
        },
        "user.UserCreateResp": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "user.UserDashboardResp": {
            "type": "object",
            "properties": {
                "finished_groups": {
                    "type": "integer"
                },
                "total_commodities": {
                    "type": "integer"
                },
                "total_groups": {
                    "type": "integer"
                },
                "total_users": {
                    "type": "integer"
                }
            }
        },
        "user.UserInfoResp": {
            "type": "object",
            "properties": {
                "user_address": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.UserInfoRespAddress"
                    }
                },
                "user_id": {
                    "type": "integer"
                },
                "user_name": {
                    "type": "string"
                },
                "user_phone": {
                    "type": "string"
                },
                "user_role": {
                    "type": "integer"
                }
            }
        },
        "user.UserInfoRespAddress": {
            "type": "object",
            "properties": {
                "area": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "detail": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_default": {
                    "type": "boolean"
                },
                "lat": {
                    "type": "number"
                },
                "lng": {
                    "type": "number"
                },
                "province": {
                    "type": "string"
                }
            }
        },
        "user.UserModifyReq": {
            "type": "object",
            "properties": {
                "user_default_address_id": {
                    "type": "integer"
                },
                "user_name": {
                    "type": "string"
                },
                "user_phone": {
                    "type": "string"
                },
                "user_role": {
                    "type": "integer"
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
