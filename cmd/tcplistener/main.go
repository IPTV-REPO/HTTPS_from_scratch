package main

import (
    "github.com/IPTV-REPO/HTTPS_from_scratch.git/internal/request"
	"fmt"
	"log"
	"net"
	
)





func main(){
   const port=":42069" ;
    ln,err:=net.Listen("tcp",port)                                         //listen on the specefic port 
    
    if err!=nil{
        log.Fatal(err)
    }

    for  {
        conn,err:=ln.Accept()                                                  //accept incoming conn

        fmt.Printf("\t There is a New HTTP Request for you  : \r\n")
	    fmt.Printf("==============================================\r\n")
	

        if err!=nil {
			log.Fatal("ERROR ","ERROR",err)
		}

        req,err:=request.RequestFromReader(conn)                                            //parse the request (from the parse request func) from the conn reader 
        if err!=nil{
			log.Fatal("ERROR ","ERROR",err)
		}

        fmt.Println("Request line:")
        fmt.Printf("- Method: %s\n", req.RequestLine.Method)
        fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
        fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
        fmt.Println("Headers:")
        for k, v := range req.Headers {
            fmt.Printf("- %s: %s\n", k, v)
        }
        fmt.Printf("Body: \n")
        fmt.Printf("%s\n",req.Body)
        
        fmt.Printf("==============================================\r\n")
       
    }
  
    

    
    
    
}