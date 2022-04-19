package main

import (
"net/http"
"errors"
"mime"
"github.com/gorilla/mux"
)

type postServer struct {
	Data map[string]*Config 
}
