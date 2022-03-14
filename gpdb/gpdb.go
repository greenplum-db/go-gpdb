package main

var (
	programName    = "gpdb"
	programVersion = "3.5.1"
)

func main() {
	// Execute the cobra CLI & run the program
	rootCmd.Execute()
}
