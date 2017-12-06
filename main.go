package main

func main() {
	//yToken := os.Getenv("Y_TOKEN")
	//if yToken == "" {
	//	panic("Y_TOKEN is required variable")
	//}
	yToken := "dict.1.1.20171201T214544Z.c0def3859d70a33d.88c6cf03a4e01eae8d732fce76205af4d35a7956"

	//dbPath := os.Getenv("DB_PATH")
	//if dbPath == "" {
	//	panic("DB_PATH is required variable")
	//}
	dbPath := "database.db"

	yDict := NewYDict(yToken)
	dbConnect := NewDB(dbPath)
	ui := NewUI(yDict, dbConnect)

	ui.Run()
}
