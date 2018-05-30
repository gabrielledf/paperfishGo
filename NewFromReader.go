package paperfishGo

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/luisfurquim/stonelizard"
	"io"
	"net/http"
	"reflect"
	"strings"
)

func NewFromReader(contract io.Reader, client *http.Client) (*WSClientT, error) {
	var err error
	var n int
	var ct stonelizard.SwaggerT
	var ws *WSClientT
	var basepath string
	var pathname string
	var pathdef stonelizard.SwaggerPathT
	var peekChar []byte
	var op *OperationT
	var swaggerParm stonelizard.SwaggerParameterT
	var paperParm *ParameterT
	var k reflect.Kind
	var method string
	var operation *stonelizard.SwaggerOperationT
	var coder interface{}

	ws = &WSClientT{
		GetOperation:     map[string]*OperationT{},
		PostOperation:    map[string]*OperationT{},
		PutOperation:     map[string]*OperationT{},
		DeleteOperation:  map[string]*OperationT{},
		OptionsOperation: map[string]*OperationT{},
		HeadOperation:    map[string]*OperationT{},
		PatchOperation:   map[string]*OperationT{},
	}
	if client == nil {
		ws.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					RootCAs:            x509.NewCertPool(),
					Certificates:       []tls.Certificate{tls.Certificate{}},
				},
				DisableCompression: true,
			},
		}
	} else {
		ws.Client = client
	}

	peekChar = make([]byte, 1)
	n, err = contract.Read(peekChar)
	if err != nil || n != 1 {
		Goose.New.Logf(1, "Error peeking first byte of service contract: %d/%c/%s", n, peekChar[0], err)
		return nil, err
	}

	if peekChar[0] == '<' {
		// This is a WSDL contract
		return nil, errors.New("WSDL is not supported yet")
	} else {
		// This is a Swagger / OpenAPI contract
		err = json.NewDecoder(io.MultiReader(bytes.NewReader(peekChar), contract)).Decode(&ct)
		if err != nil {
			Goose.New.Logf(1, "Error decoding service contract: %s", err)
			return nil, err
		}

		ws.Host = ct.Host
		basepath = ct.BasePath
		if basepath[0] == '/' {
			basepath = basepath[1:]
		}
		if basepath[len(basepath)-1] == '/' {
			basepath = basepath[:len(basepath)-1]
		}
		ws.BasePath = basepath
		ws.Schemes = ct.Schemes

		// consumes
		coder, err = getCoder(ct.Consumes)
		if err != nil {
			Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
			return nil, err
		}

		ws.Encoder = coder.(Encoder)

		// Produces
		coder, err = getCoder(ct.Produces)
		if err != nil {
			Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
			return nil, err
		}

		ws.Decoder = coder.(Decoder)

		for pathname, pathdef = range ct.Paths {
		OperLoop:
			for method, operation = range pathdef {
				op = new(OperationT)

				if pathname[0] == '/' {
					pathname = pathname[1:]
				}
				op.Path = pathname

				// consumes
				coder, err = getCoder(operation.Consumes)
				if err != nil {
					Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
					return nil, err
				}
				op.Encoder = coder.(Encoder)

				// Produces
				coder, err = getCoder(operation.Produces)
				if err != nil {
					Goose.New.Logf(1, "Error parsing 'consumes' global encoding: %s", err)
					return nil, err
				}
				op.Decoder = coder.(Decoder)

				if operation.Schemes == nil {
					op.Schemes = ws.Schemes
				} else {
					op.Schemes = operation.Schemes
				}

				Goose.New.Logf(4, " %s %#v -> %s+%s -> %s %#v\n", operation.Consumes, op.Encoder, method, operation.OperationId, operation.Produces, op.Decoder)

				switch strings.ToLower(method) {
				case "get":
					ws.GetOperation[operation.OperationId] = op
				case "post":
					ws.PostOperation[operation.OperationId] = op
				case "put":
					ws.PutOperation[operation.OperationId] = op
				case "delete":
					ws.DeleteOperation[operation.OperationId] = op
				case "options":
					ws.OptionsOperation[operation.OperationId] = op
				case "head":
					ws.HeadOperation[operation.OperationId] = op
				case "patch":
					ws.PatchOperation[operation.OperationId] = op
				default:
					Goose.New.Logf(1, "Ignoring operation %s.%s: %s", method, operation.OperationId, ErrUnknownMethod)
					continue OperLoop
				}

				for _, swaggerParm = range operation.Parameters {
					if swaggerParm.Type != "" {
						k, err = getKind(swaggerParm.Type)
					} else if (swaggerParm.Schema != nil) && (swaggerParm.Schema.Type != "") {
						k, err = getKind(swaggerParm.Schema.Type)
					} else {
						k = kindString
					}
					if err != nil {
						Goose.New.Logf(1, "Ignoring operation %s.%s.%s: %s", method, operation.OperationId, swaggerParm.Name, err)
					}

					paperParm = &ParameterT{
						Name: swaggerParm.Name,
						Kind: k,
					}

					switch strings.ToLower(swaggerParm.In) {
					case "path":
						op.PathParm = append(op.PathParm, paperParm)
					case "header":
						op.HeaderParm = append(op.HeaderParm, paperParm)
					case "query":
						op.QueryParm = append(op.QueryParm, paperParm)
					case "body":
						op.BodyParm = paperParm
					case "form":
						op.FormParm = append(op.FormParm, paperParm)
					}
				}
			}
		}
	}

	return ws, nil
}
