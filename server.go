package main

import (
  "github.com/mrdulin/go-rpc-cnode/modules/user"
  api "github.com/mrdulin/go-rpc-cnode/utils/http"
  "log"
  "net"
  "net/rpc"
  "net/rpc/jsonrpc"
)

const (
  port string = "3000"
  baseurl string = "https://cnodejs.org/api/v1"
)

//type HttpConn struct {
//  in  io.Reader
//  out io.Writer
//}
//
//func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
//func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
//func (c *HttpConn) Close() error                      { return nil }

func main() {
  httpClient := api.NewClient()
  userService := user.Service{httpClient, baseurl}

  l, e := net.Listen("tcp", ":" + port)
  if e != nil {
    log.Fatal("listen error:", e)
  }
  defer l.Close()
  
  rpcserver := rpc.NewServer()
  e = rpcserver.RegisterName("UserService", &userService)
  if e != nil {
    log.Fatal("RegisterName error:", e) 
  }
  
  // https://stackoverflow.com/questions/36610140/call-golang-call-jsonrpc-with-curl
  //http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  //if r.URL.Path == "/rpc" {
  //  serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
  //  w.Header().Set("Content-type", "application/json")
  //  w.WriteHeader(200)
  //  err := rpcserver.ServeRequest(serverCodec)
  //  if err != nil {
  //    log.Printf("Error while serving JSON request: %v", err)
  //    http.Error(w, "Error while serving JSON request, details have been logged.", 500)
  //    return
  //  }
  //}
  //}))
  
  for {
   conn, err := l.Accept()
   if err != nil {
     log.Fatal("accept error:", err)
   }
   go rpcserver.ServeCodec(jsonrpc.NewServerCodec(conn))
  }
}

