import { useSuspenseQuery } from "@tanstack/react-query";
import { getTodos } from "@/api";

const Todos = () => {
  const { data: todos, isFetching } = useSuspenseQuery({
    queryKey: ["todo"],
    queryFn: getTodos,
  });

  return (
    <div>
      {isFetching && <div>Updating...</div>}
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
          {todos?.map((todo) => (
            <tr key={todo.id}>
              <td>{todo.id}</td>
              <td>{todo.title}</td>
              <td>{todo.description}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Todos;
