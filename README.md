# paperfishGo
RESTful client for the Go language.

The Swagger client is almost complete but has not been extensively tested. The API is nearly stable and
changes are unlikely to occur at this point, unless we encounter some unforeseen major issue — and even
then, we will do our best to keep any changes as small as possible. The HTTP methods currently implemented
are GET and POST only; others will come soon. Consider it beta code with respect to GET and POST methods.

The SOAP client is very early alpha code and is not suitable even for testing. Consider this part of the
software for study purposes only.


## Features

- Automatic client generation from Swagger/OpenAPI and WSDL contracts
- Support for HTTP GET and POST methods
- WebSocket connections with event-based communication
- XOP/MTOM multipart SOAP response handling
- Optional digital signature of HTTP requests using [paracentric](https://github.com/luisfurquim/paracentric)


## Example

```Go
   httpclient = &http.Client{
      Transport:  &http.Transport{
         TLSClientConfig:     <some TLS configuration>,
         // Any other configuration that fits your needs
      },
   }

   ws, err = paperfishGo.NewFromURL("https://your.host/swagger.json", httpclient)
   if err != nil {
      fmt.Printf("Error fetching swagger: %s\n", err)
      return err
   }

   // Example below assumes that your.host offers an operation identified by "someOperationId"
   // which was defined using the HTTP POST method
   httpStatus, err = ws.Post("someOperationId", map[string]interface{}{
      // Any parameters needed by the "someOperationId" operation
      .
      .
      .
   }, &resp)
   if err != nil {
      fmt.Printf("Error: %s\n", err)
      return err
   }

   fmt.Printf("Response: %d [%s]\n", httpStatus, resp)

   // Example below assumes that your.host offers an operation identified by "anotherOperationId"
   // which was defined using the HTTP GET method
   httpStatus, err = ws.Get("anotherOperationId", map[string]interface{}{
      // Any parameters needed by the "anotherOperationId" operation
      .
      .
      .
   }, &id)
   if err != nil {
      fmt.Printf("Error: %s\n", err)
      return
   }

   fmt.Printf("Response: %d [%d]\n", httpStatus, id)

   // Example below assumes that your.host offers an operation identified by "yetAnotherOperationId"
   // which was defined using the HTTP GET method and starts a WebSocket connection which defines an
   // event called websocketEventId
   httpStatus, err = ws.Get("yetAnotherOperationId", map[string]interface{}{
      // Any parameters needed by the "yetAnotherOperationId" operation
      .
      .
      .
   }, &wsock)
   if err != nil {
      fmt.Printf("Error connecting to websocket yetAnotherOperationId: %s\n", err)
      return
   }

   fmt.Printf("Response: %d [%d]\n", httpStatus, *wsock)

   wg.Add(1)

   wsock.On("websocketEventId", func(<parameters are the data sent by the server event>) {
      // Code to execute whenever the event fires
      .
      .
      .
   })

   wg.Wait()
```


## Request Signing

paperfishGo supports optional digital signing of HTTP requests using the
[paracentric](https://github.com/luisfurquim/paracentric) PKI library. When enabled, every request
is signed with RSA-PSS (SHA256), and the signature is sent via HTTP headers. This is compatible with
the signing mechanism implemented in the JavaScript [paperfish](https://github.com/luisfurquim/paperfish) library.

Each request includes two additional headers:
- `X-Request-Signature` — base64-encoded RSA-PSS signature
- `X-Request-Signer` — the signer's key ID (derived from the certificate's SubjectKeyId)

The message that is signed follows the format: `METHOD+urlPath` for requests without a body,
or `METHOD+urlPath\nbody` for requests that include a body (e.g. POST).

### Example

```Go
   pki := paracentric.New()

   err = pki.NewPemKeyFromFile("private.pem", "password")
   if err != nil {
      fmt.Printf("Error loading private key: %s\n", err)
      return
   }

   err = pki.NewPemCertFromFile("cert.pem")
   if err != nil {
      fmt.Printf("Error loading certificate: %s\n", err)
      return
   }

   ws.SetPki(pki)

   // From this point on, all GET and POST requests will be digitally signed.
   httpStatus, err = ws.Post("someOperationId", map[string]interface{}{
      .
      .
      .
   }, &resp)
```
