package colui

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type collection struct {
	Dir  string
	Data library.Collection

	Media          map[string]library.Media
	MissingEntries []string
}

func (c *collection) Save() error {
	return library.WriteCollection(c.Dir, c.Data)
}

func (c *collection) ReorderGroups() {
	groups := c.Data.Groups

	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].Order < groups[j].Order
	})
}

func (c *collection) Invalidate() error {
	entries, err := os.ReadDir(c.Dir)
	if err != nil {
		return err
	}

	mediaEntries := []string{}
	c.Media = make(map[string]library.Media)

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		p := path.Join(c.Dir, name)

		if entry.IsDir() {
			m, err := library.ReadMedia(p)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				return err
			}

			c.Media[name] = m
			mediaEntries = append(mediaEntries, name)
		}
	}

	var collectionEntries []string
	for _, group := range c.Data.Groups {
		for _, entry := range group.Entries {
			collectionEntries = append(collectionEntries, entry.Path)
		}
	}

	c.MissingEntries = utils.SliceDifference(collectionEntries, mediaEntries)

	return nil
}

func prettyPrintCollection(col *collection) {
	fmt.Println("--------------------")
	fmt.Printf("Type:   %s\n", col.Data.Type)
	fmt.Printf("Id:     %s\n", col.Data.Id)

	fmt.Printf("Name:   %s\n", col.Data.General.Name)
	fmt.Println()

	fmt.Printf("Images:\n")
	fmt.Printf(" Cover:  %s\n", col.Data.Images.Cover)
	fmt.Printf(" Logo:   %s\n", col.Data.Images.Logo)
	fmt.Printf(" Banner: %s\n", col.Data.Images.Banner)
	fmt.Println()

	fmt.Printf("Groups:\n")

	for _, group := range col.Data.Groups {
		fmt.Printf(" Name:  %s\n", group.Name)
		fmt.Printf(" Order: %d\n", group.Order)
		fmt.Printf(" Entries:\n")
		for _, entry := range group.Entries {
			fmt.Printf("  Path:        %s\n", entry.Path)
			fmt.Printf("  Name:        %s\n", entry.Name)
			fmt.Printf("  Search Slug: %s\n", entry.SearchSlug)
			fmt.Printf("  Order:       %d\n", entry.Order)
			fmt.Println("  ---")
		}
		fmt.Println(" ---")
	}

	fmt.Println("--------------------")

	fmt.Printf("Missing Entries (%d):\n", len(col.MissingEntries))
	for _, entry := range col.MissingEntries {
		m := col.Media[entry]
		fmt.Printf(" %s: %s\n", entry, m.General.Title)
	}

	fmt.Println("--------------------")
}

func validateNumber(val string) error {
	_, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return errors.New("must be valid number")
	}

	return nil
}

// TODO(patrik): Better implementation?
func clearScreen() {
	fmt.Println("\033[2J")
}

func newCollectionTypePicker() *huh.Select[types.CollectionType] {
	return huh.NewSelect[types.CollectionType]().
		Options(
			huh.NewOption("Anime", types.CollectionTypeAnime),
			huh.NewOption("Series", types.CollectionTypeSeries),
		)
}

type NewCollectionForm struct {
	Type types.CollectionType
	Name string
}

func FormNewCollection(initialName string) (NewCollectionForm, error) {
	res := NewCollectionForm{
		Name: initialName,
	}

	form := huh.NewForm(
		huh.NewGroup(
			newCollectionTypePicker().
				Title("Collection Type").
				Value(&res.Type),
			huh.NewInput().
				Title("Collection Name").
				Value(&res.Name),
		),
	)

	err := form.Run()
	if err != nil {
		return NewCollectionForm{}, err
	}

	res.Name = strings.TrimSpace(res.Name)

	return res, nil
}

func EditMainMenu(data library.Collection, dir string) error {
	col := &collection{
		Dir:  dir,
		Data: data,
	}

	err := col.Invalidate()
	if err != nil {
		return fmt.Errorf("failed to invalidate collection: %w", err)
	}

	for {
		clearScreen()
		prettyPrintCollection(col)

		option, err := runCollectionMainMenu()
		if err != nil {
			return fmt.Errorf("failed to run collection main menu")
		}

		quit := false

		switch option {
		case collectionMainMenuGeneral:
			err := runCollectionGeneralForm(col)
			if err != nil {
				return fmt.Errorf("failed to run collection general form: %w", err)
			}
		case collectionMainMenuImages:
			err := runCollectionImagesForm(col)
			if err != nil {
				return fmt.Errorf("failed to run collection images form: %w", err)
			}
		case collectionMainMenuGroups:
			err := runCollectionGroupsForm(col)
			if err != nil {
				return fmt.Errorf("failed to run collection groups form: %w", err)
			}
		case collectionMainMenuQuit:
			quit = true
		default:
			return fmt.Errorf("unsupported menu option: %w", err)
		}

		if quit {
			break
		}

		err = library.WriteCollection(dir, data)
		if err != nil {
			return fmt.Errorf("failed to write collection to disk: %w", err)
		}
	}

	return nil
}

type collectionMainMenu string

const (
	collectionMainMenuUnknown collectionMainMenu = "unknown"
	collectionMainMenuGeneral collectionMainMenu = "general"
	collectionMainMenuImages  collectionMainMenu = "images"
	collectionMainMenuGroups  collectionMainMenu = "groups"
	collectionMainMenuQuit    collectionMainMenu = "quit"
)

func runCollectionMainMenu() (collectionMainMenu, error) {
	var selected collectionMainMenu
	form := huh.NewSelect[collectionMainMenu]().
		Title("Main Menu").
		Options(
			huh.NewOption("Edit General Info", collectionMainMenuGeneral),
			huh.NewOption("Edit Images", collectionMainMenuImages),
			huh.NewOption("Edit Groups", collectionMainMenuGroups),
			huh.NewOption("Quit", collectionMainMenuQuit),
		).
		Value(&selected)
	err := form.Run()
	if err != nil {
		return collectionMainMenuUnknown, err
	}

	return selected, nil
}

func reorderGroupEntries(group *library.CollectionGroup) {
	sort.SliceStable(group.Entries, func(i, j int) bool {
		return group.Entries[i].Order < group.Entries[j].Order
	})
}
