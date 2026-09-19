/**
 * @jest-environment jsdom
 */
import { axios } from "@/lib/axios";
import { IRecipes, useGetRecipes } from "@/services/recipes";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { renderHook, waitFor } from "@testing-library/react";

// Arrange
jest.mock("@/lib/axios", () => ({
  axios: { get: jest.fn() },
}));

const mockAxiosGetFn = jest.mocked(axios.get);

const setUpTanStack = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: 0 } },
  });

  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );

  return { queryClient, wrapper };
};

const mockInitialRecipeList: IRecipes = {
  total: 1,
  results: [
    {
      id: 1,
      name: "RECIPE_NAME",
      description: "RECIPE_DESCRIPTION",
      imageUrl: "https://example.image",
      difficulty: { id: "1", name: "EASY" },
      duration: { id: "1", name: "HALF_HOUR" },
      ingredients: [{ id: "1", description: "RECIPE_INGREDIENT_01" }],
      instructions: [{ id: "1", description: "RECIPE_INSTRUCTION_01" }],
      creator: { id: "1", name: "DEV_POOL" },
      isFavorite: false,
      rating: { average: 5, total: 200 },
      createdAt: "2026-08-20T10:00:00Z",
      updatedAt: "2026-08-20T10:00:00Z",
    },
  ],
};

describe("Validate use get recipe from API", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("Should get recipe response success", async () => {
    // Arrange
    mockAxiosGetFn.mockResolvedValue({ data: mockInitialRecipeList });

    // Act
    const { wrapper } = setUpTanStack();
    const { result } = renderHook(() => useGetRecipes(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    // // Assert
    expect(mockAxiosGetFn).toHaveBeenCalledWith("/recipes", {
      params: undefined,
    });
    expect(result.current.data).toEqual(mockInitialRecipeList);
  });

  it("Should get recipe response empty array when page equal 2", async () => {
    // Arrange
    const mockResults = { ...mockInitialRecipeList, result: [] };
    mockAxiosGetFn.mockResolvedValue({ data: mockResults });

    // Act
    const { wrapper } = setUpTanStack();
    const { result } = renderHook(() => useGetRecipes({ page: 2 }), {
      wrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    // Assert
    expect(mockAxiosGetFn).toHaveBeenCalledWith("/recipes", {
      params: { page: 2 },
    });
    expect(result.current.data).toEqual(mockResults);
  });

  it('Should return "NETWORK ERROR" when API fail', async () => {
    // Arrange
    mockAxiosGetFn.mockRejectedValue(new Error("NETWORK ERROR"));

    // Act
    const { wrapper } = setUpTanStack();
    const { result } = renderHook(() => useGetRecipes(), { wrapper });

    // Assert
    await waitFor(() => expect(result.current.isError).toBeTruthy());

    expect(result.current.error?.message).toBe("NETWORK ERROR");
    expect(result.current.data).toBeUndefined();
  });
});
