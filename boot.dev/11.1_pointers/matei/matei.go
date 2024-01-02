package main

import "fmt"

type Message struct {
	Recipient string
	Text      string
}

// Don't touch above this line

// Fix the bug in the `sendMessage` function. It's *supposed* to print a nicely formatted message to the console containing an SMS's recipient and message body. However, it's not working as expected. Run the code and see what happens, then fix the bug.

func sendMessage(m Message) {
	fmt.Printf("To: %v\n", m.Recipient)
	fmt.Printf("Message: %v\n", m.Text)
}

// Don't touch below this line

func test(recipient string, text string) {
	m := Message{Recipient: recipient, Text: text}
	sendMessage(m)
	fmt.Println("=====================================")
}

func main() {
	test("Lane", "Textio is getting better everyday!")
	test("Allan", "This pointer stuff is weird...")
	test("Tiffany", "What time will you be home for dinner?")
}
