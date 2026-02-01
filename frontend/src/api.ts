type Todos = {
  id: string;
  title: string;
  description: string;
};

type CreateTodoRequest = {
  title: string;
  description: string;
};

export const getTodos = async (): Promise<Todos[]> => {
  return Promise.resolve([]);
};

export const postTodo = async (request: CreateTodoRequest): Promise<void> => {};
