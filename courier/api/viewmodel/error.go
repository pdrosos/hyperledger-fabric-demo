package viewmodel

type Error struct {
	Status int
	Code   string
	Title  string
	Errors []ErrorDetails
}

type ErrorDetails struct {
	Title string
	Field string
}
