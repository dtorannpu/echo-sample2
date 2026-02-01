import TodoForm from "@/features/todo/TodoForm.tsx";
import { useState } from "react";
import Todos from "@/features/todo/Todos.tsx";
import type { Todo } from "@/features/todo/types.ts";

const TodoPage = () => {
  const [todo, setTodo] = useState<Todo | undefined>(undefined);
  const onComplete = () => setTodo(undefined);

  return (
    <div>
      <div>
        <TodoForm todo={todo} onComplete={onComplete} />
      </div>
      <div>
        <Todos onUpdate={setTodo} onDelete={() => setTodo(undefined)} />
      </div>
    </div>
  );
};

export default TodoPage;
