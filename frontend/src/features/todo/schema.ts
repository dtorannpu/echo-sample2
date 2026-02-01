import * as z from "zod";
import type { CreateTodoRequest } from "@/features/todo/types.ts";

export const createTodoSchema = z.object({
  title: z
    .string()
    .min(1, "タイトルは必須です")
    .max(100, "タイトルは100文字以内で入力してください"),
  description: z
    .string()
    .min(1, "内容は必須です")
    .max(1000, "内容は1000文字以内で入力してください"),
}) satisfies z.ZodType<CreateTodoRequest>;

export type FormValues = z.infer<typeof createTodoSchema>;
