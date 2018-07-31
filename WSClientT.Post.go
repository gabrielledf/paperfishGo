package paperfishGo

import (
	"io"
	"bytes"
	"io/ioutil"
	"strings"
	"net/http"
)

func (ws *WSClientT) Post(opName string, input map[string]interface{}, output interface{}) (int, error) {
	var err error
	var ok bool
	var op *OperationT
	var trackId string
	var targetURI string
	var sch, scheme string
	var lenBodyParm int
	var req *http.Request
	var resp *http.Response
	var postdata []byte
	var postparm []*ParameterT
	var buf []byte
	var soapbuf []string

	if op, ok = ws.PostOperation[opName]; !ok {
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

	scheme = "http"

	for _, sch = range op.Schemes {
		if strings.ToLower(sch) == "https" {
			scheme += "s"
		}
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

	if op.BodyParm != nil {
		postparm = []*ParameterT{op.BodyParm}
	} else {
		postparm = op.FormParm
	}

	postdata, err = ws.prepareBody(op, opName, postparm, input)
	if err != nil {
		Goose.Fetch.Logf(1, "Error preparing body for %s: %s", opName, err)
		return 0, err
	}

	req, err = http.NewRequest("POST", targetURI, bytes.NewReader(postdata))
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

	Goose.Fetch.Logf(3, "TID:[%s] request URL path %s", trackId, req.URL.Path)
	resp, err = ws.Client.Do(req)
	if err != nil {
		Goose.Fetch.Logf(6, "TID:[%s] Error fetching %s:%s", trackId, opName, err)
		return 0, err
	}
	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		Goose.Fetch.Logf(6, "TID:[%s] Error fetching response for %s:%s", trackId, opName, err)
		return 0, err
	}

	Goose.Fetch.Logf(0, "RESP: %s\n----------\n\n", buf)

	if xopxmlEnvelopRE.Match(buf) {
		soapbuf = strings.Split(string(buf),"\n")
		buf = []byte(soapbuf[len(soapbuf)-2])
		Goose.Fetch.Logf(0, "REBUF: %s\n----------\n\n", buf)
	}

	err = ws.PostOperation[opName].Decoder.Decode(bytes.NewReader(buf), output)
	if err != nil && err != io.EOF {
		Goose.Fetch.Logf(6, "TID:[%s] Error decoding response for %s:%s", trackId, opName, err)
		return 0, err
	}

	return resp.StatusCode, nil
}
