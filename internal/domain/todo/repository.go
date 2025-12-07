package todo

type TodoRepository interface {
	FindAll() ([]Todo, error)
	FindByID(id uint) (*Todo, error)
}
