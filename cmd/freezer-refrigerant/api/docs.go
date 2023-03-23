// Code generated by swaggo/swag. DO NOT EDIT.

package api

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "ice.io",
            "url": "https://ice.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/tokenomics/{userId}/extra-bonus-claims": {
            "post": {
                "description": "Claims an extra bonus for the user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tokenomics"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the user",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/tokenomics.ExtraBonusSummary"
                        }
                    },
                    "400": {
                        "description": "if validations fail",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "if not authorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "if not allowed",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "if user not found or no extra bonus available",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "if already claimed",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "if syntax fails",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "504": {
                        "description": "if request times out",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tokenomics/{userId}/mining-sessions": {
            "post": {
                "description": "Starts a new mining session for the user, if not already in progress with another one.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tokenomics"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the user",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.StartNewMiningSessionRequestBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/tokenomics.MiningSummary"
                        }
                    },
                    "400": {
                        "description": "if validations fail",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "if not authorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "if not allowed",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "if user not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "if mining is in progress or if a decision about negative mining progress is required",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "if syntax fails",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "504": {
                        "description": "if request times out",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tokenomics/{userId}/pre-staking": {
            "put": {
                "description": "Starts or updates pre-staking for the user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tokenomics"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer \u003cAdd access token here\u003e",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the user",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.StartOrUpdatePreStakingRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tokenomics.PreStakingSummary"
                        }
                    },
                    "400": {
                        "description": "if validations fail",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "if not authorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "if not allowed",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "if syntax fails",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "504": {
                        "description": "if request times out",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.StartNewMiningSessionRequestBody": {
            "type": "object",
            "properties": {
                "resurrect": {
                    "description": "Specify this if you want to resurrect the user.\n` + "`" + `true` + "`" + ` recovers all the lost balance, ` + "`" + `false` + "`" + ` deletes it forever, ` + "`" + `null/undefined` + "`" + ` does nothing. Default is ` + "`" + `null/undefined` + "`" + `.",
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "main.StartOrUpdatePreStakingRequestBody": {
            "type": "object",
            "properties": {
                "allocation": {
                    "type": "integer",
                    "maximum": 100,
                    "example": 100
                },
                "years": {
                    "type": "integer",
                    "maximum": 5,
                    "example": 1
                }
            }
        },
        "server.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "SOMETHING_NOT_FOUND"
                },
                "data": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "error": {
                    "type": "string",
                    "example": "something is missing"
                }
            }
        },
        "tokenomics.ExtraBonusSummary": {
            "type": "object",
            "properties": {
                "availableExtraBonus": {
                    "type": "integer",
                    "example": 2
                }
            }
        },
        "tokenomics.MiningRateBonuses": {
            "type": "object",
            "properties": {
                "extra": {
                    "type": "integer",
                    "example": 300
                },
                "preStaking": {
                    "type": "integer",
                    "example": 300
                },
                "t1": {
                    "type": "integer",
                    "example": 100
                },
                "t2": {
                    "type": "integer",
                    "example": 200
                },
                "total": {
                    "type": "integer",
                    "example": 300
                }
            }
        },
        "tokenomics.MiningRateSummary-coin_ICE": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string",
                    "example": "1,234,232.001"
                },
                "bonuses": {
                    "$ref": "#/definitions/tokenomics.MiningRateBonuses"
                }
            }
        },
        "tokenomics.MiningRateType": {
            "type": "string",
            "enum": [
                "positive",
                "negative",
                "none"
            ],
            "x-enum-varnames": [
                "PositiveMiningRateType",
                "NegativeMiningRateType",
                "NoneMiningRateType"
            ]
        },
        "tokenomics.MiningRates-tokenomics_MiningRateSummary-coin_ICE": {
            "type": "object",
            "properties": {
                "base": {
                    "$ref": "#/definitions/tokenomics.MiningRateSummary-coin_ICE"
                },
                "positiveTotalNoPreStakingBonus": {
                    "$ref": "#/definitions/tokenomics.MiningRateSummary-coin_ICE"
                },
                "preStaking": {
                    "$ref": "#/definitions/tokenomics.MiningRateSummary-coin_ICE"
                },
                "standard": {
                    "$ref": "#/definitions/tokenomics.MiningRateSummary-coin_ICE"
                },
                "total": {
                    "$ref": "#/definitions/tokenomics.MiningRateSummary-coin_ICE"
                },
                "totalNoPreStakingBonus": {
                    "$ref": "#/definitions/tokenomics.MiningRateSummary-coin_ICE"
                },
                "type": {
                    "$ref": "#/definitions/tokenomics.MiningRateType"
                }
            }
        },
        "tokenomics.MiningSession": {
            "type": "object",
            "properties": {
                "endedAt": {
                    "type": "string",
                    "example": "2022-01-03T16:20:52.156534Z"
                },
                "free": {
                    "type": "boolean",
                    "example": true
                },
                "resettableStartingAt": {
                    "type": "string",
                    "example": "2022-01-03T16:20:52.156534Z"
                },
                "startedAt": {
                    "type": "string",
                    "example": "2022-01-03T16:20:52.156534Z"
                },
                "warnAboutExpirationStartingAt": {
                    "type": "string",
                    "example": "2022-01-03T16:20:52.156534Z"
                }
            }
        },
        "tokenomics.MiningSummary": {
            "type": "object",
            "properties": {
                "availableExtraBonus": {
                    "type": "integer",
                    "example": 2
                },
                "miningRates": {
                    "$ref": "#/definitions/tokenomics.MiningRates-tokenomics_MiningRateSummary-coin_ICE"
                },
                "miningSession": {
                    "$ref": "#/definitions/tokenomics.MiningSession"
                },
                "miningStreak": {
                    "type": "integer",
                    "example": 2
                },
                "remainingFreeMiningSessions": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "tokenomics.PreStakingSummary": {
            "type": "object",
            "properties": {
                "allocation": {
                    "type": "integer",
                    "example": 100
                },
                "bonus": {
                    "type": "integer",
                    "example": 100
                },
                "years": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "latest",
	Host:             "",
	BasePath:         "/v1w",
	Schemes:          []string{"https"},
	Title:            "Tokenomics API",
	Description:      "API that handles everything related to write-only operations for user's tokenomics.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
