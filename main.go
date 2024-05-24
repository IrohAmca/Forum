package main
func main() {
	createDatabase()
	insertData(1, "Sebo", "Göt")
	Query(1)
	deleteData(1)
	Query(1)
}
