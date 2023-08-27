// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/login": {
            "post": {
                "tags": [
                    "login"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Login request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.LoginUserDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/users": {
            "get": {
                "tags": [
                    "user"
                ],
                "summary": "Get users",
                "responses": {
                    "200": {
                        "description": "Created user",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/user.UserDto"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            },
            "post": {
                "tags": [
                    "user"
                ],
                "summary": "Create user",
                "parameters": [
                    {
                        "description": "Create request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.CreateUserDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created user",
                        "schema": {
                            "$ref": "#/definitions/user.UserDto"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            },
            "delete": {
                "tags": [
                    "user"
                ],
                "summary": "Delete user",
                "parameters": [
                    {
                        "description": "Delete request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.DeleteUserDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "tags": [
                    "user"
                ],
                "summary": "Get current user info",
                "responses": {
                    "200": {
                        "description": "Created user",
                        "schema": {
                            "$ref": "#/definitions/user.UserDto"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/users/recovery": {
            "post": {
                "tags": [
                    "user"
                ],
                "summary": "Password recovery",
                "parameters": [
                    {
                        "description": "Delete request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/passwordRecovery.PasswordRecoveryRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/users/recovery/email": {
            "post": {
                "tags": [
                    "user"
                ],
                "summary": "Password recovery email",
                "parameters": [
                    {
                        "description": "Send password recovery email request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.SendPasswordRecoveryDto"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        }
    },
    "definitions": {
        "passwordRecovery.PasswordRecoveryRequest": {
            "type": "object",
            "properties": {
                "attemptId": {
                    "type": "string"
                },
                "newPassword": {
                    "type": "string"
                }
            }
        },
        "user.CreateUserDto": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "user.DeleteUserDto": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "user.LoginUserDto": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "user.SendPasswordRecoveryDto": {
            "type": "object",
            "properties": {
                "username": {
                    "type": "string"
                }
            }
        },
        "user.UserDto": {
            "type": "object",
            "properties": {
                "createdOn": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "isDeleted": {
                    "type": "boolean"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "username": {
                    "type": "string"
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