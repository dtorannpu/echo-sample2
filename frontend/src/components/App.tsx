import "@/components/App.css";
import { Suspense } from "react";
import TodoPage from "@/features/todo/TodoPage.tsx";

const App = () => {
  return (
    <Suspense fallback="Loading...">
      <TodoPage />
    </Suspense>
  );
};

export default App;
