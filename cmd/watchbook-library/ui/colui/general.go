package colui

import "github.com/charmbracelet/huh"

func runCollectionGeneralForm(col *collection) error {
	clearScreen()
	prettyPrintCollection(col)

	form := huh.NewForm(
		huh.NewGroup(
			newCollectionTypePicker().
				Title("Collection Type").
				Value(&col.Data.Type),
			huh.NewInput().
				Title("Collection Name").
				Value(&col.Data.General.Name),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	return nil
}
