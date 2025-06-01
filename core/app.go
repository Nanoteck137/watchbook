package core

import (
	"github.com/nanoteck137/pyrin/trail"
	"github.com/nanoteck137/watchbook/config"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
)

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	Logger() *trail.Logger

	DB() *database.Database
	Config() *config.Config

	WorkDir() types.WorkDir

	Bootstrap() error
}
