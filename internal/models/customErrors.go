package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record in database")

var ErrInvalidCredentials = errors.New("models: invalid credentials")

var ErrAlreadyHasSkill = errors.New("models: character already has that custom skill")

var ErrDuplicateFileName = errors.New("models: file of that name already exists")

var ErrNameTaken = errors.New("models: a user with that name already exists")
