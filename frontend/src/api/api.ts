import client from "./client";

export const getAuthHealth = async () => {
  try {
    const response = await client.get("/auth-well-known-config");
    return response.data;
  } catch {
    throw new Error("Please confirm your auth server is up");
  }
};
