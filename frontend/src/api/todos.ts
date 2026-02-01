import client from "@/api/client";
import type { CreateTodoRequest, Todo } from "@/features/todo/types.ts";

export const getTodos = async () =>
  client.get<Todo[]>("/todos").then((res) => res.data);
export const createTodo = async (request: CreateTodoRequest) =>
  client.post("/todos", request);
