package bb

import "fmt"

func Wrap(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}

	return nil
}
