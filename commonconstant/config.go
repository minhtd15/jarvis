package commonconstant

import "github.com/pkg/errors"

const (
	XApiKeyHeader string = "X-Api-Key"
	Authorization string = "Authorization"
)

var ErrCourseNotExist = errors.New("Course you request does not exist")
var ErrUserNotExist = errors.New("User not exist")
