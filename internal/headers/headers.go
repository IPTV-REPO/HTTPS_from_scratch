package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {      // NewHeaders creates and returns a new Headers map. 
	return make(Headers)
}



func isValidHeaderChar(b byte) bool {                      // valid header characters 
	switch {
	case b >= 'a' && b <= 'z':
		return true
	case b >= 'A' && b <= 'Z':
		return true
	case b >= '0' && b <= '9':
		return true
	case bytes.ContainsAny([]byte{b}, "!#$%&'*+-.^_`|~"):
		return true
	default:
		return false
	}
}



func (h Headers) Parse(data []byte) (n int, done bool, err error){        

	CRLF:=[]byte("\r\n")                                                 	                                          
	idx:=bytes.Index(data,CRLF)                                            
	if idx==0 {
		return 2,true,nil
		
	}
	if idx==-1{
		return 0,false,nil
	}

	CLN:=[]byte(":")
    idxColon:=bytes.Index(data,CLN)
	if idxColon==-1 {
		return 0,false,errors.New("invalid header format: missing colon")
	}
	if idxColon == 0 || data[idxColon-1] == ' '{
		return 0,false,errors.New("invalid header format")
	}

	HostPart := data[:idxColon]                                                          // Extract the header key part (before the colon)

	//  Character Validation Loop
	for _, b := range HostPart {
		if !isValidHeaderChar(b) {
			return 0, false, errors.New("invalid character in header key")
		}
	}

	//  Lowercase and Trim the Key
	host := strings.ToLower(strings.TrimSpace(string(HostPart)))
	if len(host) == 0 {
		return 0, false, errors.New("empty header key")
	}

	valuePart:=string(data[idxColon+1 : idx])                                    // Extract the header value part (after the colon and before CRLF)
	value:=strings.TrimSpace(valuePart)

	h[host]=value                                                             // Store the header in the map with the lowercase key and trimmed value

	return idx+2,false,nil
	

}


