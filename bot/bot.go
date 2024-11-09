package bot

import "fmt"

func SendMessage(to, message string) error {
	fmt.Printf("Sending message to %s: %s", to, message)
	return nil
}
