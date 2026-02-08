import { getAuthHealth } from "@/api/api";
import { logout } from "@/config";
import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { hasAuthParams, useAuth } from "react-oidc-context";

interface ProtectedAppProps {
  children: React.ReactNode;
}

const ProtectedApp = ({ children }: ProtectedAppProps) => {
  const { isPending: getAuthHealthIsPending, error: getAuthHealthError } =
    useQuery({
      queryKey: ["getAuthHealth"],
      queryFn: getAuthHealth,
      retry: false,
    });

  const auth = useAuth();

  useEffect(() => {
    if (getAuthHealthIsPending || getAuthHealthError) {
      return;
    }
    if (
      !(
        hasAuthParams() ||
        auth.isAuthenticated ||
        auth.activeNavigator ||
        auth.isLoading
      )
    ) {
      void auth.signinRedirect();
    }
  }, [auth, getAuthHealthIsPending, getAuthHealthError]);

  const anyLoading = getAuthHealthIsPending || auth.isLoading;
  const anyErrorMessage = getAuthHealthError?.message || auth.error?.message;

  if (anyLoading) {
    return <h1>Loading...</h1>;
  }
  if (anyErrorMessage) {
    return (
      <>
        <h1>We've hit a snag</h1>
        <div>{anyErrorMessage}</div>
      </>
    );
  }
  if (!auth.isAuthenticated) {
    return (
      <>
        <h1>We've hit a snag</h1>
        <div>Unable to sign in</div>
      </>
    );
  }
  return (
    <>
      <div>
        <button onClick={logout}>ログアウト</button>
      </div>
      {children}
    </>
  );
};

export default ProtectedApp;
