import { axios } from "@/lib/axios";
import { useQuery } from "@tanstack/react-query";

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

interface IRecipes {
  total: number;
  results: IRecipesItem[];
}

export const useGetRecipes = () => {
  const response = useQuery<IRecipes>({
    queryKey: ["recipes"],
    queryFn: async (): Promise<IRecipes> => {
      const response = await axios.get<IRecipes>("/recipes");
      return response.data;
    },
  });
  return response;
};
