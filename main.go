package main

import "os"

func main() {
	a := App{}
	a.Initialise(
		os.Getenv("FINK_DB_USERNAME"),
		os.Getenv("FINK_DB_PASSWORD"),
		os.Getenv("FINK_DB_NAME"))

	a.Run(":8080")
}
