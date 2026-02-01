import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getTodos, createTodo } from "@/api";

const Todos = () => {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ["todos"],
    queryFn: getTodos,
  });

  const mutation = useMutation({
    mutationFn: createTodo,
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ["todos"] });
    },
  });

  return (
    <div>
      <ul>
        {query.data?.map((todo) => (
          <li key={todo.id}>{todo.title}</li>
        ))}
      </ul>

      <button
        onClick={() => {
          mutation.mutate({
            title: "Do Laundry",
            description: "test",
          });
        }}
      >
        Add Todo
      </button>
    </div>
  );
};

export default Todos;
