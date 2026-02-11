import { privateClient } from "@/api/client";
import type {
  CreateTodoRequest,
  TodoList,
  UpdateParam,
} from "@/features/todo/types.ts";
import type { UUID } from "@/types.ts";

export const getTodos = async () =>
  privateClient.get<TodoList>("/todos").then((res) => res.data);
export const createTodo = async (request: CreateTodoRequest) =>
  privateClient.post("/todos", request);
export const deleteTodo = async (id: UUID) =>
  privateClient.delete(`/todos/${id}`);
export const updateTodo = async (param: UpdateParam) =>
  privateClient.put(`/todos/${param.id}`, param.request);
