"use client";

import { useEffect, useState, type MouseEvent } from "react";

import { RecipeCard, RecipeLevel } from "@/components/RecipeMenu";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/bases";
import { useGetRecipes } from "@/services/recipes";
import { useSearchParams } from "next/navigation";

const ITEMS_PER_PAGE = 12;

// Show every page up to 5; past that, keep first/last plus current ±1
// e.g. 1 2 3 … 7 | 1 … 5 6 7 … 12 | 1 … 5 6 7
const getPageRange = (
  currentPage: number,
  totalPages: number,
): (number | "ellipsis")[] => {
  if (totalPages <= 5) {
    return Array.from({ length: totalPages }, (_, i) => i + 1);
  }

  const start =
    currentPage >= totalPages - 1
      ? totalPages - 2
      : Math.max(2, currentPage - 1);
  const end = currentPage <= 2 ? 3 : Math.min(totalPages - 1, currentPage + 1);

  const range: (number | "ellipsis")[] = [1];
  if (start > 2) range.push("ellipsis");
  for (let page = start; page <= end; page++) range.push(page);
  if (end < totalPages - 1) range.push("ellipsis");
  range.push(totalPages);

  return range;
};

const pageHref = (page: number) => `/recipes?page=${page}`;

const RecipeList = () => {
  const search = useSearchParams();
  const pageSearchParam = search.get("page");

  const [currentPage, setCurrentPage] = useState(Number(pageSearchParam));

  const { data, isPending, isFetching, error, refetch } = useGetRecipes({
    page: currentPage,
    limit: ITEMS_PER_PAGE,
  });

  // refetch() is the only fetch trigger. It has to run from an effect rather
  // than straight from the click handler: useQuery pushes the new queryFn to
  // its observer in an effect of its own, so a refetch() called during the
  // handler would still request the previous page.
  useEffect(() => {
    refetch();
  }, [currentPage, refetch]);

  const handlePageChange =
    (page: number) => (event: MouseEvent<HTMLAnchorElement>) => {
      event.preventDefault();
      if (page === currentPage) return;

      setCurrentPage(page);
      // Keeps the URL shareable without an RSC round-trip; the Next router
      // picks up native history updates.
      window.history.replaceState(null, "", pageHref(page));
    };

  if (isPending) return <p>Loading ...</p>;
  if (error) return <p>Failed to load recipes.</p>;

  const { total = 0, results: recipes = [] } = data ?? {};
  const totalPages = Math.max(1, Math.ceil(total / ITEMS_PER_PAGE));
  const isFirstPage = currentPage === 1;
  const isLastPage = currentPage === totalPages;

  return (
    <div className="px-6 pt-8 pb-10">
      <h1 className="wongnok-text-h2">All Recipes</h1>
      <p className="wongnok-text-base">
        {total} recipes from home cooks · page {currentPage} of {totalPages}
      </p>

      {recipes.length === 0 ? (
        <p className="mt-6 text-muted-foreground">No recipes found.</p>
      ) : (
        <div
          className={`grid grid-col-1 sm:grid-cols-2 md:grid-cols-4 gap-6 mt-6 transition-opacity ${
            isFetching ? "opacity-60" : ""
          }`}
        >
          {recipes.map((recipe) => (
            <RecipeCard
              key={recipe.id}
              name={recipe.name}
              imageUrl={recipe.imageUrl}
              level={recipe.difficulty.id.toUpperCase() as RecipeLevel}
              owner={{ name: recipe.creator.name }}
            />
          ))}
        </div>
      )}

      <Pagination className="mt-8">
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious
              href={pageHref(isFirstPage ? currentPage : currentPage - 1)}
              aria-disabled={isFirstPage}
              onClick={handlePageChange(
                isFirstPage ? currentPage : currentPage - 1,
              )}
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
                  onClick={handlePageChange(item)}
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
              onClick={handlePageChange(
                isLastPage ? currentPage : currentPage + 1,
              )}
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  );
};

export default RecipeList;
