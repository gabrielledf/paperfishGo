# paperfishGo
RESTfull client for the Go language.

The swagger client is almost complete, but not extensively tested. The API is almost stabilized and changes,
now, is not likely to occur, except if we face some unpredicted big problem and, even in this case, we'll do our
best efforts to make any change as small as possible. HTTP methods implemented are GET and POST only, others will
come soon. So, consider it beta code in respect to GET and PUT methods.

The SOAP client is very alpha code and not suitable even for tests. So, consider this part of the software for study
purposes only.


## Example

```Go
   transp =

   httpclient = &http.Client{
      Transport:  &http.Transport{
         TLSClientConfig:     <some TLS configuratio>,
         // Any other configuration that fits to your needs
      },
   }

   ws, err = paperfishGo.NewFromURL("https://your.host/swagger.json",httpclient)
   if err != nil {
      fmt.Printf("Error fetching swagger: %s\n",err)
      return err
   }

   // Example below assumes that your.host offers an operation identified by "someOperationId"
   // which was defined using HTTP POST method
   httpStatus, err = ws.Post("someOperationId", map[string]interface{}{
      // Any parameters needed by the "someOperationId" operation
      .
      .
      .
   }, &resp)
   if err != nil {
      fmt.Printf("Error: %s\n",err)
      return err
   }

   fmt.Printf("Response: %d [%s]\n", httpStatus, resp)

   // Example below assumes that your.host offers an operation identified by "anotherOperationId"
   // which was defined using HTTP GET method
   httpStatus, err = ws.Get("anotherOperationId", map[string]interface{}{
      // Any parameters needed by the "anotherOperationId" operation
      .
      .
      .
   }, &id)
   if err != nil {
      fmt.Printf("Erro fazendo upload do pdf: %s\n",err)
      return
   }

   fmt.Printf("Response: %d [%d]\n", httpStatus, id)

   // Example below assumes that your.host offers an operation identified by "yetAnotherOperationId"
   // which was defined using HTTP GET method and starts a websocket connection which defines an
   // event called websocketEventId
   httpStatus, err = ws.Get("yetAnotherOperationId", map[string]interface{}{
      // Any parameters needed by the "yetAnotherOperationId" operation
      .
      .
      .
   }, &wsock)
   if err != nil {
      fmt.Printf("Error connecting to websocket yetAnotherOperationId: %s\n",err)
      return
   }

   fmt.Printf("Response: %d [%d]\n", httpStatus, *wsock)

   wg.Add(1)

   wsock.On("websocketEventId",func(<parameters are the data sent by the server event>){
      // Code to execute whenever the event fires
      .
      .
      .
   })

   wg.Wait()
```

