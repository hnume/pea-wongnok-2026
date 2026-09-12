"use client";

import {
  QueryClient,
  QueryClientProvider,
  environmentManager,
} from "@tanstack/react-query";
import * as React from "react";

const makeQueryClient = () => {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 1000 * 60 * 3,
        gcTime: 1000 * 60 * 10,
        retry: 0,
        refetchOnWindowFocus: false,
      },
    },
  });
};

let browserQueryClient: QueryClient | undefined = undefined;

const getQueryClient = () => {
  if (environmentManager.isServer()) {
    return makeQueryClient();
  } else {
    if (!browserQueryClient) browserQueryClient = makeQueryClient();
    return browserQueryClient;
  }
};

const TanstackQueryProvider = (props: { children: React.ReactNode }) => {
  const queryClient = getQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      {props.children}
    </QueryClientProvider>
  );
};

export default TanstackQueryProvider;
