package CommandingDiscord_test

import (
	"CommandingDiscord"
	"fmt"
	"testing"
)

var (
	TestHandler *CommandingDiscord.CommandHandler
	BotID       = "TestID"
	UID         = "TestID2"
)

func init() {
	fmt.Println("Init has been called.")

	_, TestHandler = CommandingDiscord.NewCommandHandler()
}

func TestNewCommandHandler(t *testing.T) {
	err, _ := CommandingDiscord.NewCommandHandler()
	if err != nil {
		t.Errorf("NewCommandHandler() returned error %v", err)
	}
}

func TestCommandHandler_Register(t *testing.T) {
	//err,TestHandler := CommandingDiscord.NewCommandHandler()
	err := TestHandler.Register("testcommand", testcommand, 10)
	if err != nil {
		t.Errorf("Register() returned error %v", err)
	}
}

func testcommand(context CommandingDiscord.Context) {
	fmt.Println("Testing~")
}
