package main

var (
	programName    = "gpdb"
	programVersion = "3.3.2"
)

func main() {
	// Execute the cobra CLI & run the program
	rootCmd.Execute()
}
