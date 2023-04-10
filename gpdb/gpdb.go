package main

var (
	programName    = "gpdb"
	programVersion = "for 7-beta.2"
)

func main() {
	// Execute the cobra CLI & run the program
	rootCmd.Execute()
}
