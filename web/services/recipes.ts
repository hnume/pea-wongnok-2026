import { axios } from "@/lib/axios";
import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

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
