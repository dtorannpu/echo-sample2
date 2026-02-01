import client from "@/api/client";
import type { CreateTodoRequest, Todo } from "@/types";

export const getTodos = () =>
  client.get<Todo[]>("/todos").then((res) => res.data);
export const createTodo = (request: CreateTodoRequest) =>
  client.post("/todos", request);
