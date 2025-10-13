package service

import "errors"

// ErrNotFound возвращается, если сущность не найдена.
var ErrNotFound = errors.New("entity not found")

// ErrAlreadyExist возвращается, если сущность уже существует,
// при попытке ее добавить.
var ErrAlreadyExist = errors.New("entity already exist")
