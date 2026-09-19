"use client";

import { useState } from "react";
import Image from "next/image";
import { CookingPot, Heart, Star } from "lucide-react";

import { Avatar, Badge, type BadgeProps } from "@/components/bases";
import { RecipeLevel } from "@/components/RecipeMenu";
import RecipeRating from "@/components/RecipeMenu/RecipeRating";
import { cn } from "@/lib/utils";
import { useGetRecipe } from "@/services/recipes";

const RECIPE_LEVEL_LABEL: Record<RecipeLevel, string> = {
  EASY: "Easy",
  MEDIUM: "Medium",
  HARD: "Hard",
};

const RECIPE_LEVEL_COLOR: Record<RecipeLevel, NonNullable<BadgeProps["color"]>> = {
  EASY: "accent",
  MEDIUM: "primary",
  HARD: "gray",
};

function RecipeImagePlaceholder() {
  return (
    <div
      aria-hidden
      className="flex size-full items-center justify-center bg-[repeating-linear-gradient(135deg,var(--muted),var(--muted)_10px,var(--background)_10px,var(--background)_20px)]"
    >
      <CookingPot className="size-12 text-muted-foreground/50" strokeWidth={1.5} />
    </div>
  );
}

type RecipeDetailProps = {
  id: string;
};

const RecipeDetail = ({ id }: RecipeDetailProps) => {
  const recipeId = Number(id);
  const { data: recipe, isPending, error } = useGetRecipe(recipeId);

  // UI-only: not persisted, not wired to POST/DELETE /recipes/{id}/favorite yet.
  const [isFavorite, setIsFavorite] = useState(false);
  // UI-only: not persisted, not wired to POST /recipes/{id}/rating yet.
  const [myRating, setMyRating] = useState(0);
  const [hoverRating, setHoverRating] = useState(0);
  const [checked, setChecked] = useState<Record<number, boolean>>({});

  if (isPending) return <p className="px-6 py-10">Loading ...</p>;

  if (error) {
    const status = (error as { response?: { status?: number } })?.response
      ?.status;
    if (status === 404) {
      return <p className="px-6 py-10">Recipe not found.</p>;
    }
    return <p className="px-6 py-10">Failed to load recipe.</p>;
  }

  if (!recipe) return null;

  const level = recipe.difficulty.id.toUpperCase() as RecipeLevel;
  const shownRating = hoverRating || myRating;

  return (
    <div className="max-w-4xl mx-auto px-6 py-8">
      <div className="relative aspect-video rounded-xl overflow-hidden bg-muted">
        {recipe.imageUrl ? (
          <Image
            src={recipe.imageUrl}
            alt={recipe.name}
            fill
            unoptimized
            className="object-cover"
          />
        ) : (
          <RecipeImagePlaceholder />
        )}

        <button
          type="button"
          aria-pressed={isFavorite}
          aria-label={isFavorite ? "Remove from favorites" : "Add to favorites"}
          onClick={() => setIsFavorite((prev) => !prev)}
          className="absolute top-3 right-3 flex size-10 cursor-pointer items-center justify-center rounded-full bg-card/90 text-secondary-foreground shadow-sm backdrop-blur-sm transition-transform duration-200 ease-[cubic-bezier(0.34,1.6,0.64,1)] outline-none hover:bg-card focus-visible:ring-3 focus-visible:ring-ring/50 active:scale-90 aria-pressed:scale-110 aria-pressed:text-primary"
        >
          <Heart aria-hidden className={cn("size-5", isFavorite && "fill-current")} />
        </button>
      </div>

      <div className="mt-6 flex flex-col gap-4">
        <h1 className="wongnok-text-h1">{recipe.name}</h1>

        <div className="flex items-center gap-4 flex-wrap">
          <span className="inline-flex items-center gap-2.5">
            <Avatar aria-hidden size="small" name={recipe.creator.name} />
            <span className="wongnok-text-sm text-muted-foreground">
              {recipe.creator.name}
            </span>
          </span>
          <Badge color={RECIPE_LEVEL_COLOR[level]}>
            {RECIPE_LEVEL_LABEL[level]}
          </Badge>
          <Badge color="gray" variant="outlined">
            {recipe.duration.name}
          </Badge>
          <RecipeRating rating={recipe.rating.average} ratingCount={recipe.rating.total} />
        </div>

        <p className="wongnok-text-base text-muted-foreground">
          {recipe.description}
        </p>
      </div>

      <div className="mt-10 grid md:grid-cols-[minmax(0,320px)_1fr] gap-8 items-start">
        <section className="bg-card border border-border rounded-2xl p-5">
          <h2 className="wongnok-text-h3 mb-4">Ingredients</h2>
          <div className="flex flex-col gap-1">
            {recipe.ingredients.map((ingredient) => {
              const isChecked = !!checked[ingredient.id];
              return (
                <button
                  key={ingredient.id}
                  type="button"
                  onClick={() =>
                    setChecked((prev) => ({
                      ...prev,
                      [ingredient.id]: !prev[ingredient.id],
                    }))
                  }
                  className="flex items-center gap-3 w-full text-left rounded-lg px-2 py-2 hover:bg-muted/50 transition-colors"
                >
                  <span
                    className={cn(
                      "flex-none size-5 rounded-md border flex items-center justify-center text-[0.7rem] font-bold",
                      isChecked
                        ? "bg-primary text-primary-foreground border-primary"
                        : "bg-card border-border",
                    )}
                  >
                    {isChecked ? "✓" : ""}
                  </span>
                  <span
                    className={cn(
                      "wongnok-text-sm",
                      isChecked
                        ? "text-muted-foreground line-through"
                        : "text-card-foreground",
                    )}
                  >
                    {ingredient.description}
                  </span>
                </button>
              );
            })}
          </div>
        </section>

        <section>
          <h2 className="wongnok-text-h3 mb-4">How to Make</h2>
          <div className="flex flex-col gap-3">
            {recipe.instructions.map((instruction, index) => (
              <div
                key={instruction.id}
                className="flex gap-3.5 items-start bg-card border border-border rounded-xl p-4"
              >
                <span className="flex-none size-7 rounded-md bg-primary-subtle text-primary flex items-center justify-center wongnok-text-xs font-bold">
                  {index + 1}
                </span>
                <p className="wongnok-text-sm text-card-foreground">
                  {instruction.description}
                </p>
              </div>
            ))}
          </div>
        </section>
      </div>

      <div className="mt-10 bg-card border border-border rounded-2xl p-6 flex flex-col gap-3">
        <span className="wongnok-text-sm font-semibold">
          {myRating ? "Your rating" : "Rate this recipe"}
        </span>
        <div className="flex items-center gap-4 flex-wrap">
          <div className="flex gap-1" onMouseLeave={() => setHoverRating(0)}>
            {[1, 2, 3, 4, 5].map((n) => (
              <button
                key={n}
                type="button"
                aria-label={`${n} star${n > 1 ? "s" : ""}`}
                onMouseEnter={() => setHoverRating(n)}
                onClick={() => setMyRating(n)}
                className="border-none bg-transparent p-0 cursor-pointer"
              >
                <Star
                  className={cn(
                    "size-7",
                    n <= shownRating
                      ? "fill-accent text-accent"
                      : "fill-transparent text-muted-foreground/40",
                  )}
                />
              </button>
            ))}
          </div>
          <span className="wongnok-text-sm text-muted-foreground">
            {myRating ? `You rated this ${myRating} out of 5` : "Tap a star"}
          </span>
        </div>
      </div>
    </div>
  );
};

export default RecipeDetail;
