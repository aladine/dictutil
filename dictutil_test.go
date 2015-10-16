package dictutiltest

import(
  "testing"
  "github.com/aladine/dictutil"
)

var datatests = []struct {
	in  string
	out string
}{
	{"/start", "[%a]"},
	{"/help", "[%-a]"},
	{"/start" + BOT_NAME, "[%+a]"},
	{"/help" + BOT_NAME, "[%+a]"},
	{"", "[%#a]"},
	{"television", "[% a]"},
	{" television", "[%0a]"},
	{"television ", "[%1.2a]"},
	{"show up", "[%-1.2a]"},
}


func TestGetDefinitionFromDb(t *testing.T) {
	user := telebot.User{
		ID:        123456,
		FirstName: "Dan",
		LastName:  "Tr",
	}

	for i := 0; i < datatests; i++ {
		t.Equal(datatests[i].in, GetDefinitionFromDb(datatests[i].in, "user", user))
	}
	group := telebot.User{
		ID:        -123456,
		FirstName: "WATO",
		LastName:  "Tr",
		Title: "We are the one"
	}

	
	for i := 0; i < datatests; i++ {
		t.Equal(datatests[i].in, GetDefinitionFromDb(datatests[i].in, "group", group))
	}
}

func TestGetDefinition(t *testing.T) {

}

func TestGetDescription(t *testing.T) {

}
