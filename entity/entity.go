package entity


type Entity interface {
	DBName()string
	ProtectedFields()[]string
}