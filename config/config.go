package config

type Services struct{
	// Server Server 
	Database Database
}

type Database struct {
	Username         string 
	Password         string 
	Host             string 
	Port             string 
	Name             string 
}

