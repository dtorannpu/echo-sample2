import { useQueryErrorResetBoundary } from "@tanstack/react-query";
import { type ReactNode } from "react";
import { ErrorBoundary } from "react-error-boundary";

const AppWithErrorBoundary = ({ children }: { children: ReactNode }) => {
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
      {children}
    </ErrorBoundary>
  );
};

export default AppWithErrorBoundary;
