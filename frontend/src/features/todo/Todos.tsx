import {
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from "@tanstack/react-query";
import { deleteTodo, getTodos } from "@/features/todo/api";
import type { Todo } from "@/features/todo/types.ts";

type Props = {
  onUpdate: (data: Todo) => void;
  onDelete: () => void;
};

const Todos = ({ onUpdate, onDelete }: Props) => {
  const queryClient = useQueryClient();

  const deleteMutation = useMutation({
    mutationFn: deleteTodo,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todo"] });
      onDelete();
    },
  });

  const { data: todoList, isFetching } = useSuspenseQuery({
    queryKey: ["todo"],
    queryFn: getTodos,
  });

  if (deleteMutation.isPending) return <div>Deleting...</div>;
  if (isFetching) return <div>Updating...</div>;

  return (
    <div>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>タイトル</th>
            <th>内容</th>
            <th>アクション</th>
          </tr>
        </thead>
        <tbody>
          {todoList.todos?.map((todo) => (
            <tr key={todo.id}>
              <td>{todo.id}</td>
              <td>{todo.title}</td>
              <td>{todo.description}</td>
              <td>
                <button onClick={() => onUpdate(todo)}>更新</button>
                <button onClick={() => deleteMutation.mutate(todo.id)}>
                  削除
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Todos;
