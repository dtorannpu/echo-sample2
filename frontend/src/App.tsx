import "@/App.css";
import Todos from "@/Todos.tsx";
import {
  QueryClient,
  QueryClientProvider,
  useQueryErrorResetBoundary,
} from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { Suspense } from "react";
import { ErrorBoundary } from "react-error-boundary";

const queryClient = new QueryClient();

const App = () => {
  const { reset } = useQueryErrorResetBoundary();

  return (
    <ErrorBoundary
      onReset={reset}
      fallbackRender={({ resetErrorBoundary }) => (
        <div>
          There was an error!
          <button onClick={() => resetErrorBoundary()}>Try again</button>
        </div>
      )}
    >
      <QueryClientProvider client={queryClient}>
        <Suspense fallback="Loading...">
          <Todos />
        </Suspense>

        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </ErrorBoundary>
  );
};

export default App;
