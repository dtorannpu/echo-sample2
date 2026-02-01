import TodoForm from "@/features/todo/TodoForm.tsx";
import Todos from "@/features/todo/api.ts";

const TodoPage = () => {
  return (
    <div>
      <div>
        <TodoForm />
      </div>
      <div>
        <Todos />
      </div>
    </div>
  );
};

export default TodoPage;
