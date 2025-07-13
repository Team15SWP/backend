package constants

const (
	//  Roles
	//  ----------------------------------------------------------------

	Admin = iota + 1
	User

	//  Handler Constants
	//  ----------------------------------------------------------------

	AuthorizationHeader = "Authorization"
	HandlerName         = "HANDLER_NAME"
	ServiceName         = "SERVICE_NAME"

	//  JWT Token Constants
	//  ----------------------------------------------------------------

	Role   = "role"
	UserID = "user"
	Name   = "name"
	Email  = "email"

	//  Entities
	//  ----------------------------------------------------------------

	TaskEntity     = "task"
	SyllabusEntity = "syllabus"

	//  Operations
	//  ----------------------------------------------------------------

	Get    = "get"
	Create = "create"
	Update = "update"
	Delete = "delete"
)
