package paperfishGo

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"reflect"
	"strings"
)

func (ws *WSClientT) Get(opName string, input map[string]interface{}, output interface{}) (int, error) {
	var err error
	var ok bool
	var op *OperationT
	var trackId string
	var cfg *websocket.Config
	var targetURI string
	var swaggerURI string
	var wsock *websocket.Conn
	var sch, scheme string
	var useTLS, useWebSock bool
	var lenBodyParm int
	var req *http.Request
	var resp *http.Response
	var WSockClient WSockClientT

	if op, ok = ws.GetOperation[opName]; !ok {
		Goose.Fetch.Logf(1, "%s", ErrUnknownOperation)
		return 0, ErrUnknownOperation
	}

	if op.BodyParm != nil {
		lenBodyParm = 1
	}

	if (len(op.PathParm) + len(op.HeaderParm) + len(op.QueryParm) + lenBodyParm + len(op.FormParm)) != len(input) {
		Goose.Fetch.Logf(1, "%s for %s", ErrWrongParmCount, opName)
		return 0, ErrWrongParmCount
	}

	trackId, err = NewTrackId()
	if err != nil {
		Goose.Fetch.Logf(1, "Error generating new Track ID: %s", err)
		return 0, err
	}

	for _, sch = range op.Schemes {
		switch strings.ToLower(sch) {
		case "https":
			useTLS = true
		case "wss":
			useTLS = true
			useWebSock = true
		case "ws":
			useWebSock = true
		}
	}

	if !useWebSock {
		scheme = "http"
		if useTLS {
			scheme += "s"
		}

		targetURI, err = ws.prepareURL(op, opName, scheme, op.PathParm, input)
		if err != nil {
			Goose.Fetch.Logf(1, "Error preparing URL for %s: %s", opName, err)
			return 0, err
		}

		targetURI, err = ws.prepareQuery(op, opName, scheme, op.QueryParm, input, targetURI)
		if err != nil {
			Goose.Fetch.Logf(1, "Error preparing URL query for %s: %s", opName, err)
			return 0, err
		}

		req, err = http.NewRequest("GET", targetURI, nil)
		if err != nil {
			return 0, err
		}

		Goose.Fetch.Logf(6, "TID:[%s] Request:%#v", trackId, req)

		req.Header.Set("Connection", "keep-alive")

		err = ws.prepareHeaders(op, opName, op.HeaderParm, input, req)
		if err != nil {
			Goose.Fetch.Logf(1, "Error preparing headers for %s: %s", opName, err)
			return 0, err
		}

		Goose.Fetch.Logf(1, "TID:[%s] request URL %s", trackId, req.URL.Path)
		resp, err = ws.Client.Do(req)
		if err != nil {
			Goose.Fetch.Logf(1, "TID:[%s] Error fetching %s:%s", trackId, opName, err)
			return 0, err
		}

		defer resp.Body.Close()
		err = ws.GetOperation[opName].Decoder.Decode(resp.Body, output)
		if err != nil && err != io.EOF {
			Goose.Fetch.Logf(1, "TID:[%s] Error decoding response for %s:%s", trackId, opName, err)
			return 0, err
		}

		return resp.StatusCode, nil

	} else {

		Goose.Fetch.Logf(4, "Suboperations of %s: %#v", opName, op.SubOperations)

		switch output.(type) {
		case **WSockClientT:
		default:
			Goose.Fetch.Logf(6, "TID:[%s] %s for %s", trackId, ErrWrongReturnParmType, opName)
			return 0, ErrWrongReturnParmType
		}

		if useTLS {
			scheme = "s"
		} else {
			scheme = ""
		}

		targetURI, err = ws.prepareURL(op, opName, "ws"+scheme, op.PathParm, input)
		if err != nil {
			Goose.Fetch.Logf(1, "Error preparing URL for %s: %s", opName, err)
			return 0, err
		}

		targetURI, err = ws.prepareQuery(op, opName, "ws"+scheme, op.QueryParm, input, targetURI)
		if err != nil {
			Goose.Fetch.Logf(1, "Error preparing URL query for %s: %s", opName, err)
			return 0, err
		}

		swaggerURI = fmt.Sprintf("ws%s://%s/swagger.json", scheme, ws.Host)

		cfg, err = websocket.NewConfig(targetURI, swaggerURI)
		if err != nil {
			Goose.Fetch.Logf(1, "Error configuring websocket connection for %s: %s", targetURI, err)
			return 0, err
		}

		cfg.TlsConfig = ws.Client.Transport.(*http.Transport).TLSClientConfig

		wsock, err = websocket.DialConfig(cfg)
		if err != nil {
			Goose.Fetch.Logf(1, "Error dialing to websocket at %s: %s", targetURI, err)
			return 0, err
		}

		WSockClient = WSockClientT{
			SubOperations: op.SubOperations,
			receiver:      make(chan []interface{}),
			cli2srvch:     make(chan WSockRequest),
			bindch:        make(chan WSockRequest),
		}

		go func() {
			var err error
			var message []interface{}
			var ok bool
			var id, stat float64

			for {
				message = []interface{}{}
				err = websocket.JSON.Receive(wsock, &message)
				if err != nil {
					Goose.Fetch.Logf(1, "Error receiving from websocket at %s: %s", targetURI, err)
					return
				}

				if id, ok = message[0].(float64); !ok {
					Goose.Fetch.Logf(1, "Error receiving from websocket at %s: first element must be a number", targetURI)
					continue
				}

				message[0] = uint32(id)

				if id == 0 {
					if _, ok = message[1].(string); !ok {
						Goose.Fetch.Logf(1, "Error receiving from websocket at %s: second element must be a string, got [%#v]", targetURI, message[1])
						continue
					}
				} else {
					if stat, ok = message[1].(float64); !ok {
						Goose.Fetch.Logf(1, "Error receiving from websocket at %s: second element must be a number, got [%#v]", targetURI, message[1])
						continue
					}

					message[1] = uint32(stat)

					/*
					   if len(message) >= 3 {
					      if _, ok = message[2].(string) ; !ok {
					         Goose.Fetch.Logf(1,"Error receiving from websocket at %s: third element must be a string, got [%#v]", targetURI, message[2])
					         continue
					      }
					   }
					*/
				}

				WSockClient.receiver <- message
			}
		}()

		go func() {
			var ok bool
			var err error
			var wstrackId uint32
			var req WSockRequest
			var pending map[uint32]CallbackT
			var pendingEvents map[string][]reflect.Value
			var fn CallbackT
			var resp []interface{}
			var evtName string
			var resppar []interface{}
			var parmval []reflect.Value
			var httpStat uint32

			pending = map[uint32]CallbackT{}
			pendingEvents = map[string][]reflect.Value{}

			for {
				Goose.Fetch.Logf(4, "Entering select")
				select {
				case req = <-WSockClient.bindch:
					for {
						wstrackId, err = NewWSTrackId()
						if err != nil {
							Goose.Fetch.Logf(1, "Error generating new Track ID: %s", err)
							return
						}

						if _, ok = pending[wstrackId]; !ok {
							break
						}
					}

					err = websocket.JSON.Send(wsock, &[]interface{}{wstrackId, "bind", req.SubOperation})
					if err != nil {
						Goose.Fetch.Logf(1, "Error sending to websocket at %s: %s", targetURI, err)
						return
					}

					pending[wstrackId] = func(name string, callback reflect.Value) CallbackT {
						Goose.Fetch.Logf(5, "On receiver -> fn=%#v", callback)
						return CallbackT{
							Callback: reflect.ValueOf(func(httpStat uint32) {
								if httpStat != 200 {
									Goose.Fetch.Logf(1, "Error binding event on websocket at %s [%#v]>", targetURI, httpStat)
									return
								}

								Goose.Fetch.Logf(5, "Get receiver -> fn=%#v", callback)

								pendingEvents[name] = append(pendingEvents[name], callback)
							}),
							FailCallback: func(int) {},
						}
					}(req.SubOperation, req.Callback)

				case req = <-WSockClient.cli2srvch:
					Goose.Fetch.Logf(5, "Received message to send to server: %#v", req)
					for {
						Goose.Fetch.Logf(5, "Will generate new Track ID")
						wstrackId, err = NewWSTrackId()
						Goose.Fetch.Logf(5, "Generated new Track ID")
						if err != nil {
							Goose.Fetch.Logf(1, "Error generating new Track ID: %s", err)
							return
						}

						Goose.Fetch.Logf(1, "Generated new Track ID %d without error", wstrackId)
						if _, ok = pending[wstrackId]; !ok {
							break
						}
					}

					Goose.Fetch.Logf(5, "Sending %#v to %s", []interface{}{wstrackId, req.SubOperation, req.Params}, targetURI)
					err = websocket.JSON.Send(wsock, &[]interface{}{wstrackId, req.SubOperation, req.Params})
					if err != nil {
						Goose.Fetch.Logf(1, "Error sending to websocket at %s: %s", targetURI, err)
						return
					}

					pending[wstrackId] = req.CallbackT

				case resp = <-WSockClient.receiver:
					if wstrackId, ok = resp[0].(uint32); !ok {
						Goose.Fetch.Logf(1, "Error interface {} is %T %T, not uint32 %#v", resp[0], resp[0], resp)
						break
					}
					Goose.Fetch.Logf(5, "Callback of %d", wstrackId)
					if wstrackId == 0 {
						evtName = resp[1].(string)

						if len(resp) != 3 {
							Goose.Fetch.Logf(1, "Error wrong parameter count to websocket at event %s", evtName)
							break
						}

						if resppar, ok = resp[2].([]interface{}); !ok {
							Goose.Fetch.Logf(1, "Error wrong 3rd parameter type at websocket event %s", evtName)
							break
						}

						for _, fn.Callback = range pendingEvents[evtName] {
							parmval, err = mkParms(fn.Callback, resppar)
							if err != nil {
								Goose.Fetch.Logf(1, "Error %s at %s", err, targetURI)
								break
							}
							go fn.Callback.Call(parmval)
						}
					} else {
						if fn, ok = pending[wstrackId]; !ok {
							Goose.Fetch.Logf(1, "Error handler not found at websocket %s", targetURI)
							break
						}

						httpStat, ok = resp[1].(uint32)
						if !ok {
							Goose.Fetch.Logf(1, "Error %s (%T) at %s", ErrProtocol, resp[1], targetURI)
							break
						}

						if int(httpStat) >= http.StatusBadRequest { // Check for error function
							Goose.Fetch.Logf(5, "Fail Callback of %d", wstrackId)
							go fn.FailCallback(int(httpStat))
						} else {
							Goose.Fetch.Logf(5, "Callback of %d => %#v", wstrackId, resp[1:])
							parmval, err = mkParms(fn.Callback, resp[1:])
							if err != nil {
								Goose.Fetch.Logf(1, "Error %s at %s", err, targetURI)
								break
							}

							Goose.Fetch.Logf(5, "Callback of %d => parmval=%#v", wstrackId, parmval)
							go fn.Callback.Call(parmval)
							Goose.Fetch.Logf(5, "Callback of %d returned", wstrackId)
						}
					}

				}
			}
		}()

		switch out := output.(type) {
		case **WSockClientT:
			*out = &WSockClient
		}

		return http.StatusOK, nil
	}
}
