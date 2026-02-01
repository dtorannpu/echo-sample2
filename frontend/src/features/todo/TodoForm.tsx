"use no memo";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { createTodoSchema, type FormValues } from "@/features/todo/schema.ts";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createTodo, updateTodo } from "@/features/todo/api";
import type { Todo } from "@/features/todo/types.ts";
import { useEffect } from "react";

type Props = {
  todo?: Todo;
  onComplete: () => void;
};

const TodoForm = ({ todo, onComplete }: Props) => {
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(createTodoSchema),
  });

  useEffect(() => {
    if (todo) {
      reset({
        title: todo.title,
        description: todo.description,
      });
    } else {
      reset({ title: "", description: "" });
    }
  }, [reset, todo]);

  const createMutation = useMutation({
    mutationFn: createTodo,
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ["todo"] });
      onComplete();
      reset({ title: "", description: "" });
    },
  });

  const updateMutation = useMutation({
    mutationFn: updateTodo,
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ["todo"] });
      onComplete();
      reset({ title: "", description: "" });
    },
  });

  const onSubmit = (values: FormValues) => {
    if (todo) {
      updateMutation.mutate({ id: todo.id, request: values });
    } else {
      createMutation.mutate(values);
    }
  };

  if (createMutation.isPending) return <div>登録中...</div>;
  if (createMutation.isError) {
    return (
      <div>
        <p>登録に失敗しました</p>
        <button onClick={() => createMutation.reset()}>もう一度実行</button>
      </div>
    );
  }

  if (updateMutation.isPending) return <div>更新...</div>;
  if (updateMutation.isError) {
    return (
      <div>
        <p>更新に失敗しました</p>
        <button onClick={() => updateMutation.reset()}>もう一度実行</button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {todo && (
        <div>
          <label>ID</label>
          <p>{todo.id}</p>
        </div>
      )}
      <div>
        <label>タイトル</label>
        <input {...register("title")} />
        <p>{errors.title?.message}</p>
      </div>
      <div>
        <label>内容</label>
        <textarea {...register("description")} />
        <p>{errors.description?.message}</p>
      </div>
      <input type="submit" value={todo ? "更新" : "登録"} />
    </form>
  );
};

export default TodoForm;
