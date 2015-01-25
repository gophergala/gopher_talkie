package main

type Server interface {
	ListenAndServer() error
}
