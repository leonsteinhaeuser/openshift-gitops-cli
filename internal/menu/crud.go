package menu

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	optionCreate = "Create"
	optionUpdate = "Update"
	optionDelete = "Delete"
	optionDone   = "Done"
)

type actionFunc func() error

func crudMenu(actionLabel string, create, update, delete actionFunc) error {
OUTER:
	for {
		prompt := promptui.Select{
			Label: actionLabel,
			Items: []string{optionCreate, optionUpdate, optionDelete, optionDone},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		switch result {
		case optionCreate:
			err := create()
			if err != nil {
				return err
			}
		case optionUpdate:
			err := update()
			if err != nil {
				return err
			}
		case optionDelete:
			err := delete()
			if err != nil {
				return err
			}
		case optionDone:
			break OUTER
		default:
			return fmt.Errorf("invalid option %s", result)
		}
	}
	return nil
}
