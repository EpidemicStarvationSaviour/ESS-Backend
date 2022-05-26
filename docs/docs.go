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
        "/admin/users": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "get all users conditionally",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "page_num",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "type",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/admin.AllUserResp"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "delete certain user",
                "parameters": [
                    {
                        "description": "User Id",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/admin.AdminDeleteUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/group/create": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "create a new group",
                "parameters": [
                    {
                        "description": "New Group Info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/group.GroupCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/group.GroupCreateResp"
                        }
                    }
                }
            }
        },
        "/group/edit/{id}": {
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Agent Edit Group",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "edit group id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Agent Group Edit",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/group.GroupEditReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/group/join": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Join a group (create new order)",
                "parameters": [
                    {
                        "description": "Join Group Info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/group.GroupJoinReq"
                        }
                    }
                ],
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
        "/group/list": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "get groups I joined conditional",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "page_num",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "type",
                        "in": "query"
                    }
                ],
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
        "/group/own": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "Agent Get Own Groups",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "page_num",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "type",
                        "in": "query"
                    }
                ],
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
        "/group/search": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "group"
                ],
                "summary": "search exist groups",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "group_type",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page_num",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "search_type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "value",
                        "in": "query"
                    }
                ],
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
                    "type": "string",
                    "example": "西湖区"
                },
                "city": {
                    "type": "string",
                    "example": "杭州市"
                },
                "detail": {
                    "type": "string",
                    "example": "浙江大学紫金港校区"
                },
                "is_default": {
                    "type": "boolean",
                    "example": true
                },
                "province": {
                    "type": "string",
                    "example": "浙江省"
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
        "admin.AdminDeleteUser": {
            "type": "object",
            "properties": {
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "admin.AllUserResp": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/admin.UserData"
                    }
                }
            }
        },
        "admin.UserAddress": {
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
        "admin.UserData": {
            "type": "object",
            "properties": {
                "user_address": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/admin.UserAddress"
                    }
                },
                "user_created_time": {
                    "type": "integer"
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
                },
                "user_updated_time": {
                    "type": "integer"
                }
            }
        },
        "group.GroupCreateReq": {
            "type": "object",
            "properties": {
                "address_id": {
                    "type": "integer"
                },
                "commodities": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "remark": {
                    "type": "string"
                },
                "user_group_id": {
                    "type": "integer"
                }
            }
        },
        "group.GroupCreateResp": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "group.GroupEditReq": {
            "type": "object",
            "properties": {
                "address_id": {
                    "type": "integer"
                },
                "commodities": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "deleted_users": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "remark": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "group.GroupInfoAddress": {
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
        "group.GroupInfoCommodity": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "number": {
                    "type": "number"
                },
                "price": {
                    "type": "number"
                },
                "total_number": {
                    "type": "number"
                },
                "type_id": {
                    "type": "integer"
                }
            }
        },
        "group.GroupInfoData": {
            "type": "object",
            "properties": {
                "commodity_detail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/group.GroupInfoCommodity"
                    }
                },
                "created_time": {
                    "type": "integer"
                },
                "creator_address": {
                    "$ref": "#/definitions/group.GroupInfoAddress"
                },
                "creator_id": {
                    "type": "integer"
                },
                "creator_name": {
                    "type": "string"
                },
                "creator_phone": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "remark": {
                    "type": "string"
                },
                "total_my_price": {
                    "type": "number"
                },
                "total_price": {
                    "type": "number"
                },
                "type": {
                    "type": "integer"
                },
                "user_number": {
                    "type": "integer"
                }
            }
        },
        "group.GroupInfoResp": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/group.GroupInfoData"
                    }
                }
            }
        },
        "group.GroupJoinData": {
            "type": "object",
            "properties": {
                "commodity_id": {
                    "type": "integer"
                },
                "number": {
                    "type": "number"
                }
            }
        },
        "group.GroupJoinReq": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/group.GroupJoinData"
                    }
                },
                "id": {
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
                    "example": 3
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
                    "type": "string",
                    "example": "西湖区"
                },
                "city": {
                    "type": "string",
                    "example": "杭州市"
                },
                "detail": {
                    "type": "string",
                    "example": "浙江大学玉泉校区"
                },
                "province": {
                    "type": "string",
                    "example": "浙江省"
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
