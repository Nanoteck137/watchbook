package colui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/utils"
)

func runGroupEditInfo(group *library.CollectionGroup) error {
	groupName := group.Name
	groupOrder := strconv.Itoa(group.Order)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Group Name").
				Value(&groupName),
			huh.NewInput().
				Title("Group Order (number)").
				Validate(validateNumber).
				Value(&groupOrder),
		),
	)
	err := form.Run()
	if err != nil {
		return err
	}

	order, _ := strconv.ParseInt(groupOrder, 10, 64)
	group.Name = groupName
	group.Order = int(order)

	return nil
}

func runGroupEditEntry(col *collection, group *library.CollectionGroup) error {

	var options []huh.Option[int]
	for i, entry := range group.Entries {
		m := col.Media[entry.Path]
		options = append(options, huh.NewOption(entry.Name+": "+m.General.Title, i))
	}

	options = append(options, huh.NewOption("Back", -1))

	var value int
	form := huh.NewSelect[int]().
		Title("Select entry").
		Options(options...).
		Value(&value)
	err := form.Run()
	if err != nil {
		return err
	}

	if value == -1 {
		return nil
	}

	clearScreen()
	// TODO(patrik): Print entry info
	prettyPrintCollection(col)

	entry := &group.Entries[value]

	{
		name := entry.Name
		searchSlug := entry.SearchSlug
		orderStr := strconv.Itoa(entry.Order)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Name").
					Value(&name),
				huh.NewInput().
					Title("Search Slug").
					Value(&searchSlug),
				huh.NewInput().
					Title("Order").
					Validate(validateNumber).
					Value(&orderStr),
			),
		)
		err := form.Run()
		if err != nil {
			return err
		}

		name = strings.TrimSpace(name)
		searchSlug = utils.Slug(searchSlug)
		order, _ := strconv.ParseInt(orderStr, 10, 64)

		entry.Name = name
		entry.SearchSlug = searchSlug
		entry.Order = int(order)
	}

	reorderGroupEntries(group)

	return nil
}

func runGroupAddEntry(col *collection, group *library.CollectionGroup) error {
	var options []huh.Option[string]

	for _, entry := range col.MissingEntries {
		m := col.Media[entry]
		options = append(options, huh.NewOption(m.General.Title, entry))
	}

	options = append(options, huh.NewOption("Back", "back"))

	var value string
	form := huh.NewSelect[string]().
		Title("Add entry").
		Options(options...).
		Value(&value)
	err := form.Run()
	if err != nil {
		return err
	}

	if value == "back" {
		return nil
	}

	m := col.Media[value]
	name := m.General.Title

	var searchWord string
	orderStr := "0"
	{
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Entry Name").
					Value(&name),
				huh.NewInput().
					Title("Entry Search Word (defaults to name)").
					Value(&searchWord),
				huh.NewInput().
					Title("Entry Order").
					Validate(validateNumber).
					Value(&orderStr),
			),
		)
		err := form.Run()
		if err != nil {
			return err
		}
	}

	if searchWord == "" {
		searchWord = name
	}

	order, _ := strconv.ParseInt(orderStr, 10, 64)

	group.Entries = append(group.Entries, library.CollectionEntry{
		Path:       value,
		Name:       name,
		SearchSlug: utils.Slug(searchWord),
		Order:      int(order),
	})

	err = col.Invalidate()
	if err != nil {
		return fmt.Errorf("failed to invalidate collection: %w", err)
	}

	return nil
}

func runGroupRemoveEntry(group *library.CollectionGroup) error {
	var options []huh.Option[int]

	for i, entry := range group.Entries {
		options = append(options, huh.NewOption(entry.Name, i))
	}

	options = append(options, huh.NewOption("Back", -1))

	var selected int
	form := huh.NewSelect[int]().
		Options(options...)
	err := form.Run()
	if err != nil {
		return err
	}

	if selected == -1 {
		return nil
	}

	group.Entries = slices.Delete(group.Entries, selected, selected+1)
	reorderGroupEntries(group)

	return nil
}

func runEditGroup(col *collection, index int) error {
	group := &col.Data.Groups[index]

	for {
		clearScreen()
		// TODO(patrik): Print group
		prettyPrintCollection(col)

		var selected string
		form := huh.NewSelect[string]().
			Title("Editing group: "+group.Name).
			Options(
				huh.NewOption("Edit info", "edit-info"),
				huh.NewOption("Edit entry ", "edit-entry"),
				huh.NewOption("Add entry", "add-entry"),
				huh.NewOption("Remove entry", "remove-entry"),
				huh.NewOption("Delete group", "delete-group"),
				huh.NewOption("Back", "back"),
			).
			Value(&selected)

		err := form.Run()
		if err != nil {
			return err
		}

		switch selected {
		case "edit-info":
			err := runGroupEditInfo(group)
			if err != nil {
				return fmt.Errorf("failed to run group edit info: %w", err)
			}

			err = col.Save()
			if err != nil {
				return fmt.Errorf("failed to save collection: %w", err)
			}
		case "edit-entry":
			err := runGroupEditEntry(col, group)
			if err != nil {
				return fmt.Errorf("failed to run group edit entry: %w", err)
			}

			err = col.Save()
			if err != nil {
				return fmt.Errorf("failed to save collection: %w", err)
			}
		case "add-entry":
			err := runGroupAddEntry(col, group)
			if err != nil {
				return fmt.Errorf("failed to run group add entry: %w", err)
			}

			err = col.Save()
			if err != nil {
				return fmt.Errorf("failed to save collection: %w", err)
			}
		case "remove-entry":
			err := runGroupRemoveEntry(group)
			if err != nil {
				return fmt.Errorf("failed to run group remove entry: %w", err)
			}

			err = col.Invalidate()
			if err != nil {
				return fmt.Errorf("failed to invalidate collection: %w", err)
			}

			err = col.Save()
			if err != nil {
				return fmt.Errorf("failed to save collection: %w", err)
			}
		case "delete-group":
			col.Data.Groups = slices.Delete(col.Data.Groups, index, index+1)

			err = col.Invalidate()
			if err != nil {
				return fmt.Errorf("failed to invalidate collection: %w", err)
			}

			return nil
		case "back":
			return nil
		}
	}
}

func runCollectionGroupsForm(col *collection) error {
	for {
		clearScreen()
		prettyPrintCollection(col)

		var selected string
		form := huh.NewSelect[string]().
			Title("Edit Collection Groups").
			Options(
				huh.NewOption("Add Group", "add-group"),
				huh.NewOption("Edit Group", "edit-group"),
				huh.NewOption("Back", "back"),
			).
			Value(&selected)

		err := form.Run()
		if err != nil {
			return err
		}

		switch selected {
		case "add-group":
			var groupName string
			var groupOrder string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Group Name").
						Value(&groupName),
					huh.NewInput().
						Title("Group Order (number)").
						Validate(validateNumber).
						Value(&groupOrder),
				),
			)
			err := form.Run()
			if err != nil {
				return err
			}

			order, _ := strconv.ParseInt(groupOrder, 10, 64)
			col.Data.Groups = append(col.Data.Groups, library.CollectionGroup{
				Name:    groupName,
				Order:   int(order),
				Entries: []library.CollectionEntry{},
			})

			col.ReorderGroups()

			err = col.Save()
			if err != nil {
				return fmt.Errorf("failed to save collection: %w", err)
			}

		case "edit-group":
			options := []huh.Option[int]{}

			for i, group := range col.Data.Groups {
				options = append(options, huh.NewOption(group.Name, i))
			}

			var groupIndex int
			form := huh.NewSelect[int]().
				Title("Select group").
				Options(options...).
				Value(&groupIndex)

			err := form.Run()
			if err != nil {
				return err
			}

			err = runEditGroup(col, groupIndex)
			if err != nil {
				return fmt.Errorf("failed to run edit group form: %w", err)
			}

		case "back":
			return nil
		default:
			return fmt.Errorf("unsupported collection group cmd: %s", selected)
		}
	}
}
