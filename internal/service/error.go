package service

import "fmt"

var (
	// ErrNoGroup возникает когда пользователь не состоит в группе
	ErrNoGroup        = fmt.Errorf("user is not a member of any group")
	ErrInvalidGroupID = fmt.Errorf("invalid group id")
	ErrNoOddMonday    = fmt.Errorf("odd monday date is not set")
	ErrTooManyGroups  = fmt.Errorf("user has too many groups")
	ErrNoOwnedGroup   = fmt.Errorf("user does not own this group")
)
