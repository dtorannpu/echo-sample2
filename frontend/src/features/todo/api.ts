import client from "@/api/client.ts";
import type {
  CreateTodoRequest,
  TodoList,
  UpdateParam,
} from "@/features/todo/types.ts";
import type { UUID } from "@/types.ts";

export const getTodos = async () =>
  client.get<TodoList>("/todos").then((res) => res.data);
export const createTodo = async (request: CreateTodoRequest) =>
  client.post("/todos", request);
export const deleteTodo = async (id: UUID) => client.delete(`/todos/${id}`);
export const updateTodo = async (param: UpdateParam) =>
  client.put(`/todos/${param.id}`, param.request);
