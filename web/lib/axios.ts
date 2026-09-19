import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";
import { getSession, signOut } from "next-auth/react";

const instance = axios.create({
  baseURL: "http://api.wongnok.localhost/api/v1",
  timeout: 5000,
  headers: { "Content-Type": "application/json" },
});

let refreshInFlight: Promise<{
  accessToken: string;
  refreshToken: string;
}> | null = null;

let refreshedTokens: { accessToken: string; refreshToken: string } | null =
  null;


async function refreshAccessToken(refreshToken: string) {
  if (!refreshInFlight) {
    refreshInFlight = axios
      .post("http://localhost:8080/api/v1/auth/refresh-token", {
        refreshToken,
      })
      .then(({ data }) => {
        refreshedTokens = {
          accessToken: data.accessToken,
          refreshToken: data.refreshToken,
        };
        return refreshedTokens;
      })
      .finally(() => {
        refreshInFlight = null;
      });
  }
  return refreshInFlight;
}

instance.interceptors.request.use(async config => {
  const session = await getSession()

  config.headers.Authorization = `Bearer ${session?.accessToken}`
  return config
})

instance.interceptors.response.use(
  response => response,
  async error => {
    const { config, response } = error;

    // config.__isRetry guards against looping forever if the retried
    // request also comes back 401 (e.g. the refresh token itself is dead).
    if (response?.status !== 401 || config.__isRetry) {
      return Promise.reject(error);
    }
    config.__isRetry = true;

    const session = await getSession();
    const refreshToken = refreshedTokens?.refreshToken ?? session?.refreshToken;
    if (!refreshToken) {
      return Promise.reject(error);
    }

    try {
      const tokens = await refreshAccessToken(refreshToken);
      config.headers.Authorization = `Bearer ${tokens.accessToken}`;
      return instance(config);
    } catch {
      // Refresh token is invalid/expired too — nothing left to do but
      // force the user back through the sign-in flow.
      refreshedTokens = null;
      await signOut();
      return Promise.reject(error);
    }
  }
);

export { instance as axios };
