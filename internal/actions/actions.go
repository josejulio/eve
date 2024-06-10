package actions

import(
	"fmt"
	"errors"

	"github.com/josejulio/eve/internal/session"
)

func ExecuteAction(action string, session session.Session) ([]string, error) {
	if action == "add_numbers" {
		a := session.GetSlot("first_number").(int)
		b := session.GetSlot("second_number").(int)
		return []string{fmt.Sprintf("%d + %d = %d !!!", a, b, a + b)}, nil
	}

	return nil, errors.New("Not implemented")
}