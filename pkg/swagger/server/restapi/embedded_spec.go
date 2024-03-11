// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "Donders Repository data stager APIs",
    "title": "dr-data-stager",
    "version": "0.1.0"
  },
  "basePath": "/v1",
  "paths": {
    "/dir": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get entities within a filesystem path",
        "parameters": [
          {
            "description": "the directory",
            "name": "dir",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dirPath"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseDirEntries"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/job": {
      "post": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "create a new stager job",
        "parameters": [
          {
            "description": "stager job data",
            "name": "data",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/jobData"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/jobInfo"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/job/{id}": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get stager job information",
        "parameters": [
          {
            "type": "string",
            "description": "job identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/jobInfo"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "job not found",
            "schema": {
              "type": "string",
              "enum": [
                "job not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "delete": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "cancel a stager job",
        "parameters": [
          {
            "type": "string",
            "description": "job identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/jobInfo"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "job not found",
            "schema": {
              "type": "string",
              "enum": [
                "job not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/jobs": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get all jobs of a user",
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "post": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "create multiple new stager jobs",
        "parameters": [
          {
            "description": "stager job data",
            "name": "data",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyJobs"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "207": {
            "description": "multi-status",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/jobs/{status}": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get list of jobs at given status",
        "parameters": [
          {
            "enum": [
              "waiting",
              "processing",
              "failed",
              "succeeded",
              "canceled"
            ],
            "type": "string",
            "description": "job status",
            "name": "status",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/ping": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "check API server health",
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "type": "string",
              "enum": [
                "pong"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "dirEntry": {
      "description": "directory entry",
      "required": [
        "name",
        "type",
        "size"
      ],
      "properties": {
        "name": {
          "description": "name of the entry",
          "type": "string"
        },
        "size": {
          "description": "size of the entry in bytes",
          "type": "integer"
        },
        "type": {
          "description": "type of the entry",
          "type": "string",
          "enum": [
            "regular",
            "dir",
            "symlink",
            "unknown"
          ]
        }
      }
    },
    "dirPath": {
      "description": "directory path data",
      "required": [
        "path"
      ],
      "properties": {
        "path": {
          "description": "path of the directory",
          "type": "string"
        }
      }
    },
    "jobData": {
      "description": "job data",
      "required": [
        "title",
        "stagerUser",
        "drUser",
        "srcURL",
        "dstURL"
      ],
      "properties": {
        "drPass": {
          "description": "password of the DR data-access account",
          "type": "string"
        },
        "drUser": {
          "description": "username of the DR data-access account",
          "type": "string"
        },
        "dstURL": {
          "description": "path or DR namespace (prefixed with irods:) of the destination endpoint",
          "type": "string"
        },
        "srcURL": {
          "description": "path or DR namespace (prefixed with irods:) of the source endpoint",
          "type": "string"
        },
        "stagerUser": {
          "description": "username of stager's local account",
          "type": "string"
        },
        "timeout": {
          "description": "allowed duration in seconds for entire transfer job (0 for no timeout)",
          "type": "integer"
        },
        "timeout_noprogress": {
          "description": "allowed duration in seconds for no further transfer progress (0 for no timeout)",
          "type": "integer"
        },
        "title": {
          "description": "short description about the job",
          "type": "string"
        }
      }
    },
    "jobID": {
      "description": "identifier for scheduled background tasks.",
      "type": "string"
    },
    "jobInfo": {
      "description": "JSON object containing scheduled job information.",
      "required": [
        "id",
        "data",
        "status"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/jobData"
        },
        "id": {
          "$ref": "#/definitions/jobID"
        },
        "status": {
          "$ref": "#/definitions/jobStatus"
        }
      }
    },
    "jobProgress": {
      "description": "job progress information",
      "required": [
        "total",
        "processed",
        "failed"
      ],
      "properties": {
        "failed": {
          "description": "number of failed files",
          "type": "integer"
        },
        "processed": {
          "description": "number of processed files",
          "type": "integer"
        },
        "total": {
          "description": "number of total files to be processed",
          "type": "integer"
        }
      }
    },
    "jobStatus": {
      "description": "status of the background task.",
      "required": [
        "status",
        "error",
        "progress"
      ],
      "properties": {
        "error": {
          "description": "job error message from the last execution.",
          "type": "string"
        },
        "progress": {
          "description": "job progress info from the last execution.",
          "$ref": "#/definitions/jobProgress"
        },
        "status": {
          "description": "job status from the last execution.",
          "type": "string",
          "enum": [
            "waiting",
            "processing",
            "failed",
            "succeeded",
            "canceled"
          ]
        }
      }
    },
    "principal": {
      "description": "authenticated client identifier",
      "type": "string"
    },
    "requestBodyJobs": {
      "description": "JSON object containing a list of job data.",
      "properties": {
        "jobs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/jobData"
          }
        }
      }
    },
    "responseBody400": {
      "description": "JSON object containing error message concerning bad client request.",
      "properties": {
        "errorMessage": {
          "description": "error message specifying the bad request.",
          "type": "string"
        }
      }
    },
    "responseBody500": {
      "description": "JSON object containing server side error.",
      "properties": {
        "errorMessage": {
          "description": "server-side error message.",
          "type": "string"
        },
        "exitCode": {
          "description": "server-side exit code.",
          "type": "integer"
        }
      }
    },
    "responseBodyJobs": {
      "description": "JSON object containing a list of job information.",
      "properties": {
        "jobs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/jobInfo"
          }
        }
      }
    },
    "responseDirEntries": {
      "description": "JSON object containing dir entries.",
      "properties": {
        "entries": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dirEntry"
          }
        }
      }
    }
  },
  "securityDefinitions": {
    "basicAuth": {
      "type": "basic"
    },
    "oauth2": {
      "type": "oauth2",
      "flow": "application",
      "tokenUrl": "https://login.dccn.nl/connect/token",
      "scopes": {
        "urn:dccn:data-stager-api:*": "general access scope for data-stager API server"
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "Donders Repository data stager APIs",
    "title": "dr-data-stager",
    "version": "0.1.0"
  },
  "basePath": "/v1",
  "paths": {
    "/dir": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get entities within a filesystem path",
        "parameters": [
          {
            "description": "the directory",
            "name": "dir",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dirPath"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseDirEntries"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/job": {
      "post": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "create a new stager job",
        "parameters": [
          {
            "description": "stager job data",
            "name": "data",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/jobData"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/jobInfo"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/job/{id}": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get stager job information",
        "parameters": [
          {
            "type": "string",
            "description": "job identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/jobInfo"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "job not found",
            "schema": {
              "type": "string",
              "enum": [
                "job not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "delete": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "cancel a stager job",
        "parameters": [
          {
            "type": "string",
            "description": "job identifier",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/jobInfo"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "404": {
            "description": "job not found",
            "schema": {
              "type": "string",
              "enum": [
                "job not found"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/jobs": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get all jobs of a user",
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      },
      "post": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "create multiple new stager jobs",
        "parameters": [
          {
            "description": "stager job data",
            "name": "data",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/requestBodyJobs"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "207": {
            "description": "multi-status",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/jobs/{status}": {
      "get": {
        "security": [
          {
            "oauth2": [
              "urn:dccn:data-stager-api:*"
            ]
          },
          {
            "basicAuth": []
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get list of jobs at given status",
        "parameters": [
          {
            "enum": [
              "waiting",
              "processing",
              "failed",
              "succeeded",
              "canceled"
            ],
            "type": "string",
            "description": "job status",
            "name": "status",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/responseBodyJobs"
            }
          },
          "400": {
            "description": "bad request",
            "schema": {
              "$ref": "#/definitions/responseBody400"
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    },
    "/ping": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "check API server health",
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "type": "string",
              "enum": [
                "pong"
              ]
            }
          },
          "500": {
            "description": "failure",
            "schema": {
              "$ref": "#/definitions/responseBody500"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "dirEntry": {
      "description": "directory entry",
      "required": [
        "name",
        "type",
        "size"
      ],
      "properties": {
        "name": {
          "description": "name of the entry",
          "type": "string"
        },
        "size": {
          "description": "size of the entry in bytes",
          "type": "integer"
        },
        "type": {
          "description": "type of the entry",
          "type": "string",
          "enum": [
            "regular",
            "dir",
            "symlink",
            "unknown"
          ]
        }
      }
    },
    "dirPath": {
      "description": "directory path data",
      "required": [
        "path"
      ],
      "properties": {
        "path": {
          "description": "path of the directory",
          "type": "string"
        }
      }
    },
    "jobData": {
      "description": "job data",
      "required": [
        "title",
        "stagerUser",
        "drUser",
        "srcURL",
        "dstURL"
      ],
      "properties": {
        "drPass": {
          "description": "password of the DR data-access account",
          "type": "string"
        },
        "drUser": {
          "description": "username of the DR data-access account",
          "type": "string"
        },
        "dstURL": {
          "description": "path or DR namespace (prefixed with irods:) of the destination endpoint",
          "type": "string"
        },
        "srcURL": {
          "description": "path or DR namespace (prefixed with irods:) of the source endpoint",
          "type": "string"
        },
        "stagerUser": {
          "description": "username of stager's local account",
          "type": "string"
        },
        "timeout": {
          "description": "allowed duration in seconds for entire transfer job (0 for no timeout)",
          "type": "integer"
        },
        "timeout_noprogress": {
          "description": "allowed duration in seconds for no further transfer progress (0 for no timeout)",
          "type": "integer"
        },
        "title": {
          "description": "short description about the job",
          "type": "string"
        }
      }
    },
    "jobID": {
      "description": "identifier for scheduled background tasks.",
      "type": "string"
    },
    "jobInfo": {
      "description": "JSON object containing scheduled job information.",
      "required": [
        "id",
        "data",
        "status"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/jobData"
        },
        "id": {
          "$ref": "#/definitions/jobID"
        },
        "status": {
          "$ref": "#/definitions/jobStatus"
        }
      }
    },
    "jobProgress": {
      "description": "job progress information",
      "required": [
        "total",
        "processed",
        "failed"
      ],
      "properties": {
        "failed": {
          "description": "number of failed files",
          "type": "integer"
        },
        "processed": {
          "description": "number of processed files",
          "type": "integer"
        },
        "total": {
          "description": "number of total files to be processed",
          "type": "integer"
        }
      }
    },
    "jobStatus": {
      "description": "status of the background task.",
      "required": [
        "status",
        "error",
        "progress"
      ],
      "properties": {
        "error": {
          "description": "job error message from the last execution.",
          "type": "string"
        },
        "progress": {
          "description": "job progress info from the last execution.",
          "$ref": "#/definitions/jobProgress"
        },
        "status": {
          "description": "job status from the last execution.",
          "type": "string",
          "enum": [
            "waiting",
            "processing",
            "failed",
            "succeeded",
            "canceled"
          ]
        }
      }
    },
    "principal": {
      "description": "authenticated client identifier",
      "type": "string"
    },
    "requestBodyJobs": {
      "description": "JSON object containing a list of job data.",
      "properties": {
        "jobs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/jobData"
          }
        }
      }
    },
    "responseBody400": {
      "description": "JSON object containing error message concerning bad client request.",
      "properties": {
        "errorMessage": {
          "description": "error message specifying the bad request.",
          "type": "string"
        }
      }
    },
    "responseBody500": {
      "description": "JSON object containing server side error.",
      "properties": {
        "errorMessage": {
          "description": "server-side error message.",
          "type": "string"
        },
        "exitCode": {
          "description": "server-side exit code.",
          "type": "integer"
        }
      }
    },
    "responseBodyJobs": {
      "description": "JSON object containing a list of job information.",
      "properties": {
        "jobs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/jobInfo"
          }
        }
      }
    },
    "responseDirEntries": {
      "description": "JSON object containing dir entries.",
      "properties": {
        "entries": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dirEntry"
          }
        }
      }
    }
  },
  "securityDefinitions": {
    "basicAuth": {
      "type": "basic"
    },
    "oauth2": {
      "type": "oauth2",
      "flow": "application",
      "tokenUrl": "https://login.dccn.nl/connect/token",
      "scopes": {
        "urn:dccn:data-stager-api:*": "general access scope for data-stager API server"
      }
    }
  }
}`))
}
