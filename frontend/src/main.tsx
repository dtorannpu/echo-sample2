import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@/index.css";
import App from "@/components/App";
import { AuthProvider } from "react-oidc-context";
import { onSigninCallback, queryClient, userManager } from "@/config";
import ProtectedApp from "@/components/ProtectedApp";
import { QueryClientProvider } from "@tanstack/react-query";
import AppWithErrorBoundary from "./components/AppWithErrorBoundary";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AppWithErrorBoundary>
      <AuthProvider
        userManager={userManager}
        onSigninCallback={onSigninCallback}
      >
        <QueryClientProvider client={queryClient}>
          <ProtectedApp>
            <App />
          </ProtectedApp>
          <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
      </AuthProvider>
    </AppWithErrorBoundary>
  </StrictMode>,
);
