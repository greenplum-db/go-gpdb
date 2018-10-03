package core

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Ask user what is the choice from the list provided.
func Prompt_choice(TotalOptions int) int {

	var choice_entered int
	fmt.Print("\nEnter your choice from the above list (eg.s 1 or 2 etc): ")

	// Start the new scanner to get the user input
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {

		// The choice entered
		choice_entered, err := strconv.Atoi(input.Text())

		// If user enters a string instead of a integer then ask to re-enter
		if err != nil {
			fmt.Println("Incorrect value: Please choose a integer (eg.s 1 or 2 etc) from the above list")
			return Prompt_choice(TotalOptions)
		}

		// If its a valid value move on
		if choice_entered > 0 && choice_entered <= TotalOptions {
			return choice_entered
		} else { // Else ask for re-entering the selection
			fmt.Println("Invalid Choice: The choice you entered is not on the list above, try again.")
			return Prompt_choice(TotalOptions)
		}
	}

	return choice_entered
}


// Prompt for confirmation
func YesOrNoConfirmation() string {

	// Start the new scanner to get the user input
	fmt.Print("You can use \"gpdb env -v <version>\" to set the env, do you wish to continue (Yy/Nn)?: ")
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {

		// The choice entered
		choice_entered := input.Text()

		// If its a valid value move on
		if YesOrNo[strings.ToLower(choice_entered)] == "y" {  // Is it Yes
			return choice_entered
		} else if YesOrNo[strings.ToLower(choice_entered)] == "n" { // Is it No
			return choice_entered
		} else { // Invalid choice, ask to re-enter
			fmt.Println("Invalid Choice: Please enter Yy/Nn, try again.")
			return YesOrNoConfirmation()
		}
	}

	return ""
}