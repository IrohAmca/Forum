package main
func main() {
	createDatabase()
	insertData(1, "Sebo", "Göt")
	Query(1)
	deleteData("Sebo")
	Query(1)
}
