import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { type FormValues, createTodoSchema } from "@/features/todo/schema.ts";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createTodo } from "@/features/todo/api";

const TodoForm = () => {
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(createTodoSchema),
  });

  const mutation = useMutation({
    mutationFn: createTodo,
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ["todo"] });
      reset();
    },
    throwOnError: true,
  });

  const onSubmit = (values: FormValues) => {
    mutation.mutate(values);
  };

  if (mutation.isPending) return <div>登録中...</div>;
  if (mutation.isError) {
    return (
      <div>
        <p>登録に失敗しました</p>
        <button onClick={() => mutation.reset()}>もう一度実行</button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
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
      <input type="submit" value="登録" />
    </form>
  );
};

export default TodoForm;
