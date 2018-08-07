package paperfishGo

import (
   "sync"
   "reflect"
)

func (wsc *WSockClientT) SendSync(opName string, params map[string]interface{}, fn interface{}, fnfail func(int)) error {
   var wg sync.WaitGroup
   var err error
   var fnval reflect.Value
//   var fnval2 reflect.Value
//   var in []reflect.Type

   fnval = reflect.ValueOf(fn)

   if fnval.Kind() != reflect.Func {
      Goose.Fetch.Logf(1, "Error %s for %s callback function", ErrWrongParmType, opName)
      return ErrWrongParmType
   }



   success := func(in []reflect.Value) []reflect.Value {
      rv := fnval.Call(in)
      wg.Done()
      return  rv
   }

   fnval2 := reflect.MakeFunc(fnval.Type(), success)

   wg.Add(1)
   err = wsc.Send(
      opName,
      params,
      fnval2.Interface(),
      func(httpStat int) {
         fnfail(httpStat)
         wg.Done()
      })
   wg.Wait()

   return err
}
