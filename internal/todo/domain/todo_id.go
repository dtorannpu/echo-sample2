package domain

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type TodoID uuid.UUID

func NewTodoID() TodoID {
	return TodoID(uuid.New())
}

func (id TodoID) String() string {
	return uuid.UUID(id).String()
}

func (id TodoID) Value() (driver.Value, error) {
	return id.String(), nil
}

func (id *TodoID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("invalid type for TodoID: %T", value)
	}

	u, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*id = TodoID(u)
	return nil
}
