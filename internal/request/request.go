package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)


const(
	stateInitialized=iota           // The initial state of the request before any data has been read. In this state,
                                    //  the request is waiting for the first line of the HTTP request to be read.
	stateDone                       // The state of the request after the entire HTTP request has been read and processed. In this state,
)

type Request struct {
	RequestLine RequestLine
	stateLine  int               // This tracks: Are we at the start? Are we done?
}

type RequestLine struct {                //the request line of http request contains those three parts 
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RightMethods(method string) bool {       //validate the method of the request line (only uppercase letters)
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}


func (r *Request) parse(data []byte) (int,error) {         //parsing the request line of the http request and validating it
	switch r.stateLine {
	case stateInitialized:                                         // In this state, we expect to read the request line of the HTTP request. We will parse the request line and validate it.
		parsedLine,consumed, err := parseRequestLine(data)        // If there is an error during parsing, we return the error. If the request line is successfully parsed, we update the Request struct with the parsed request line and transition to the stateDone state.
		if err != nil {
			return 0,err
		}
		if consumed>0 {                                            // If the request line is successfully parsed, we update the Request struct with the parsed request line and transition to the stateDone state.
			r.RequestLine = *parsedLine
		    r.stateLine = stateDone
		}
		return consumed,nil
	case stateDone:
		return 0,errors.New("trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
}

func httpValid(version string) bool {                                  // This function checks if the HTTP version string is valid. 
	HttpVersion := strings.Split(version, "/")                         //It checks if the version string is in the format "HTTP/x.x" where x.x is a valid HTTP version (e.g., "1.1").
	if len(HttpVersion) != 2 {                                         // The function splits the version string by the "/" character and checks if it has exactly two parts. If not, it returns false.
		return false
	}

	if HttpVersion[0] != "HTTP" {
		return false
	}

	if HttpVersion[1] != "1.1" {
		return false
	}

	return true
}


func parseRequestLine(b []byte ) (*RequestLine,int, error) {              //same for header parsing but for the request line of the http request
	CRLF:=[]byte("\r\n")
	idx:=bytes.Index(b, CRLF)
	if idx==-1{
		return nil,0,nil
	}

    line :=string(b[:idx])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil,0,errors.New("invalid request line")
	}

	if !RightMethods(parts[0]) {
		return nil,0, errors.New("invalid method")
	}

	if !httpValid(parts[2]) {
		return nil,0,errors.New("invalid http version")
	}

	versionOnly := strings.TrimPrefix(parts[2], "HTTP/")                    // Extract the version number (e.g., "1.1") from the HTTP version string (e.g., "HTTP/1.1") by removing the "HTTP/" prefix.

	rl:=&RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   versionOnly,
	}
	return rl,idx+2,nil
}


																	// This function reads data from the provided io.Reader 
																	// and attempts to parse it as an HTTP request. 
																	// It uses a buffer to read data in chunks and calls the parse method of the Request struct to process the data. 
																	//The function continues reading and parsing until the entire HTTP request has been successfully parsed or an error occurs.
func RequestFromReader(reader io.Reader) (*Request, error) {             
	buffer:=make([]byte, 1024)
	readIdx:=0

	r :=&Request{
		stateLine:stateInitialized,
	}

	for r.stateLine!=stateDone{                          // The function continues reading and parsing until the entire HTTP request has been successfully parsed or an error occurs.
		if readIdx==len(buffer){
			newbuffer:=make([]byte,len(buffer)*2)
			copy(newbuffer,buffer)		
			buffer=newbuffer
		}

		n, err := reader.Read(buffer[readIdx:])      // Read data from the reader into the buffer starting at the current read index. The number of bytes read is stored in n, and any error that occurs during reading is stored in err.
		readIdx+=n

		consumed, parseErr := r.parse(buffer[:readIdx])         // Call the parse method of the Request struct to process the data in the buffer up to the current read index. The parse method returns the number of bytes consumed from the buffer and any error that occurs during parsing.
																//The parse method returns the number of bytes consumed from the buffer and any error that occurs during parsing.
		if parseErr != nil {
			return nil, parseErr
		}

		if consumed>0{                                        		// If the parse method successfully consumed some bytes from the buffer, we need to shift the remaining unprocessed data to the beginning of the buffer and update the read index accordingly. 
																	//This allows us to continue reading new data into the buffer without losing any unprocessed data.
			copy(buffer, buffer[consumed:readIdx])
			readIdx-=consumed
		}
		
		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.stateLine != stateDone {
					return nil, errors.New("unexpected end of input")

				}
				break
				
			}
			return nil, err
		}
	}
	return r,nil
	
}
