package main

import "git.apache.org/thrift.git/lib/go/thrift"

func main() {
	server, err := NewServiceLocator().Get("thrift_service")
	if err != nil {
		panic(err)
	}

	server.(*thrift.TSimpleServer).Serve()
}
