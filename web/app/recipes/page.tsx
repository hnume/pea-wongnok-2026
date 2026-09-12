"use client";

import { RecipeCard, RecipeLevel } from "@/components/RecipeMenu";
import recipes from "./_data.json";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/bases";
import { IRecipes, useGetRecipes } from "@/services/recipes";
import { useQuery } from "@tanstack/react-query";
import { axios } from "@/lib/axios";

const ITEMS_PER_PAGE = 12;

// Show every page up to 5; past that, keep first/last plus current ±1
// e.g. 1 2 3 … 7 | 1 … 5 6 7 … 12 | 1 … 5 6 7
// const getPageRange = (
//   currentPage: number,
//   totalPages: number,
// ): (number | "ellipsis")[] => {
//   if (totalPages <= 5) {
//     return Array.from({ length: totalPages }, (_, i) => i + 1);
//   }

//   const start =
//     currentPage >= totalPages - 1
//       ? totalPages - 2
//       : Math.max(2, currentPage - 1);
//   const end = currentPage <= 2 ? 3 : Math.min(totalPages - 1, currentPage + 1);

//   const range: (number | "ellipsis")[] = [1];
//   if (start > 2) range.push("ellipsis");
//   for (let page = start; page <= end; page++) range.push(page);
//   if (end < totalPages - 1) range.push("ellipsis");
//   range.push(totalPages);

//   return range;
// };

const pageHref = (page: number) => `/recipes?page=${page}`;

const RecipesPage = ({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}) => {
  const { data: recipesData, isLoading, error } = useGetRecipes();

  if (isLoading) return <p>Loading ...</p>;

  // const { page } = searchParams;
  // const totalPages = Math.max(1, Math.ceil(recipes.length / ITEMS_PER_PAGE));
  // const requestedPage = Number(page) || 1;
  // const currentPage = Math.min(
  //   Math.max(1, Math.floor(requestedPage)),
  //   totalPages,
  // );

  // const recipesByPage = recipes.slice(
  //   (currentPage - 1) * ITEMS_PER_PAGE,
  //   currentPage * ITEMS_PER_PAGE,
  // );

  // const isFirstPage = currentPage === 1;
  // const isLastPage = currentPage === totalPages;

  const { total, results: recipes = [] } = recipesData ?? ({} as IRecipes);

  return (
    <div className="px-6 pt-8 pb-10">
      <h1 className="wongnok-text-h2">All Recipes</h1>
      {/* <p className="wongnok-text-base">
        {recipes.length} recipes from home cooks · page {currentPage} of{" "}
        {totalPages}
      </p> */}

      <div className="grid grid-col-1 sm:grid-cols-2 md:grid-cols-4 gap-6 mt-6">
        {recipes.map((recipe, index) => (
          <RecipeCard
            key={index}
            name={recipe.name}
            imageUrl={recipe.imageUrl}
            level={recipe.difficulty.id.toUpperCase() as RecipeLevel}
            owner={{ name: recipe.creator.name }}
          />
        ))}
      </div>
      {/* <Pagination className="mt-8">
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious
              href={pageHref(isFirstPage ? currentPage : currentPage - 1)}
              aria-disabled={isFirstPage}
            />
          </PaginationItem>
          {getPageRange(currentPage, totalPages).map((item, index) => (
            <PaginationItem
              key={item === "ellipsis" ? `ellipsis-${index}` : item}
            >
              {item === "ellipsis" ? (
                <PaginationEllipsis />
              ) : (
                <PaginationLink
                  href={pageHref(item)}
                  isActive={item === currentPage}
                >
                  {item}
                </PaginationLink>
              )}
            </PaginationItem>
          ))}
          <PaginationItem>
            <PaginationNext
              href={pageHref(isLastPage ? currentPage : currentPage + 1)}
              aria-disabled={isLastPage}
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination> */}
    </div>
  );
};

export default RecipesPage;
