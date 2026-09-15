import { axios } from "@/lib/axios";
import {
  useMutation,
  UseMutationOptions,
  useQuery,
  type UseQueryOptions,
} from "@tanstack/react-query";
import { AxiosError } from "axios";

/** -------------------- GET ----------------------- */
interface IRecipesItem {
  id: number;
  name: string;
  description: string;
  imageUrl: string;
  difficulty: { id: string; name: string };
  duration: { id: string; name: string };
  ingredients: { id: string; description: string }[];
  instructions: { id: string; description: string }[];
  creator: { id: string; name: string };
  isFavorite?: boolean;
  rating: { average: number; total: number };
  createdAt: string | Date;
  updatedAt: string | Date;
}

export interface IRecipes {
  total: number;
  results: IRecipesItem[];
}

export interface IGetRecipesQuery {
  page?: number;
  limit?: number;
}

export type GetRecipesQueryOptions = Omit<
  UseQueryOptions<IRecipes>,
  "queryKey" | "queryFn"
>;

// page/limit are part of the queryKey, so each page is cached separately and
// changing either triggers a fetch automatically.
export const useGetRecipes = (
  params?: IGetRecipesQuery,
  queryOptions?: GetRecipesQueryOptions,
) => {
  const response = useQuery<IRecipes>({
    queryKey: ["recipes", params?.page, params?.limit],
    queryFn: async (): Promise<IRecipes> => {
      const response = await axios.get<IRecipes>("/recipes", { params });
      return response.data;
    },
    ...queryOptions,
  });
  return response;
};

/** -------------------- CREATE ----------------------- */

interface IRecipeCreationPayload {
  name: string;
  description: string;
  imageUrl?: string;
  difficultyId: "easy" | "medium" | "hard";
  durationId: "10m" | "30m" | "60m" | "long";
  ingredients: { description: string }[];
  instructions: { description: string }[];
}

interface IRecipeCreationResponse {
  id: number;
}

export type CreateRecipesMutationOptions = Omit<
  UseMutationOptions<
    IRecipeCreationResponse,
    AxiosError,
    IRecipeCreationPayload
  >,
  "queryKey" | "queryFn"
>;

export const useCreateRecipe = (
  payload: IRecipeCreationPayload,
  mutationOption?: CreateRecipesMutationOptions,
) => {
  const response = useMutation({
    mutationKey: ["create", "recipes"],
    mutationFn: async (): Promise<IRecipeCreationResponse> => {
      const response = await axios.post<IRecipeCreationResponse>("/recipes", {
        payload,
      });
      return response.data;
    },
    ...mutationOption,
  });
  return response;
};

/** -------------------- OTHERS ----------------------- */

// PUT/PATCH -> Update
export const useUpdateRecipe = () => {};

// DELETE -> Delete
export const useDeleteRecipe = () => {};
