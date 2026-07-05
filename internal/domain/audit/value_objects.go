package audit

// Action enum
type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionLogin  Action = "login"
	ActionLogout Action = "logout"
	ActionAccess Action = "access"
)

func (a Action) IsValid() bool {
	switch a {
	case ActionCreate, ActionUpdate, ActionDelete, ActionLogin, ActionLogout, ActionAccess:
		return true
	}
	return false
}
