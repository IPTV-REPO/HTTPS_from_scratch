package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)



func GetLineChannel(f io.ReadCloser)<-chan string{
    
    pipechannel:=make(chan string)
    
    go func(){
        defer close(pipechannel)
        defer f.Close()

        var currentline string
        data:=make([]byte,8)
      
        for{
        
            count,err:=f.Read(data)
            if count>0{

                str:=string(data[:count])
                part:=strings.Split(str, "\n")

                for i := 0; i < len(part)-1; i++ {
                    fullline:=currentline+part[i]
                    fmt.Printf("read: %s\n",fullline) 
                    currentline="" 
                }
                currentline+=part[len(part)-1]
            }

            if err==io.EOF {
            
                if currentline!="" {
                   fmt.Printf("read: %s\n",currentline) 
            
                }
                fmt.Println("----------------------")
                fmt.Println("End of file reached")
                break
            }
        
            if err!=nil{
               log.Fatal(err)
            }
          
        }


    }()

    return  pipechannel
}


func main(){
    port:=":42069"
    ln,err:=net.Listen("tcp",port)    
    
    if err!=nil{
        log.Fatal(err)
    }

    for  {
        conn,err:=ln.Accept()
        if err!=nil{
            log.Fatal(err)
            break
        }
        
        fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())
        line:=GetLineChannel(conn)

    
        for l:=range line{
            fmt.Printf("%s\n", l)
    
        }
        conn.Close()
    }
  
    

    
    
    
}