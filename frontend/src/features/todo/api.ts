import client from "@/api/client.ts";
import type { CreateTodoRequest, Todo } from "@/features/todo/types.ts";
import type { UUID } from "@/types.ts";

export const getTodos = async () =>
  client.get<Todo[]>("/todos").then((res) => res.data);
export const createTodo = async (request: CreateTodoRequest) =>
  client.post("/todos", request);
export const deleteTodo = async (id: UUID) => client.delete(`/todos/${id}`);
