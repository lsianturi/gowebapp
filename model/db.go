package main

type db interface {
	SelectPeople() ([]*Person, error)
}