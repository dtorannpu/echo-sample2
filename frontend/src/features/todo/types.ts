import type { UUID } from "@/types.ts";

export type Todo = {
  id: UUID;
  title: string;
  description: string;
};

export type CreateTodoRequest = {
  title: string;
  description: string;
};
