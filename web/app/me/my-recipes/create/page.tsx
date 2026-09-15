"use client";

import { Button } from "@/components/bases";
import RecipeTheDishForm from "./_containers/RecipeTheDishForm";
import RecipeEffortForm from "./_containers/RecipeEffortForm";
import RecipeIngredientForm from "./_containers/RecipeIngredientForm";
import RecipeInstructionForm from "./_containers/RecipeInstructionForm";
import { FormProvider, useForm } from "react-hook-form";
import {
  LEVELS,
  RecipeCreationFormValues,
  TIMES,
} from "@/types/FormValues/recipeCreationForm";
import { zodResolver } from "@hookform/resolvers/zod";
import z4 from "zod/v4";
import { IRecipeCreationPayload, useCreateRecipe } from "@/services/recipes";
import { useRouter } from "next/navigation";
import { useState } from "react";

const recipeCreationValidateSchema = z4.object({
  name: z4.string().nonempty({ error: "กรุณากรอกชื่อเมนูอาหาร" }),
  description: z4.string().nonempty({ error: "กรุณากรอกรายละเอียดเมนูอาหาร" }),
  // imageUrl: z4.url().optional(),
  imageUrl: z4
    .url({ error: "กรุณากรอก URL รูปภาพให้ถูกต้อง" })
    .or(z4.literal(""))
    .optional(),
  level: z4.enum(LEVELS, {
    error: "Please select level of recipe",
  }),
  time: z4.enum(TIMES, {
    error: "Please select time of recipe",
  }),
  ingredients: z4
    .array(
      z4.object({
        description: z4.string().nonempty({
          error: "Please input ingredient",
        }),
      }),
    )
    .min(1, { error: "Please add at least one ingredient" }),
  instructions: z4
    .array(
      z4.object({
        description: z4.string().nonempty({
          error: "Please input how to make",
        }),
      }),
    )
    .min(1, { error: "Please add at least one ingredient" }),
});

const RecipeCreationPage = () => {
  const router = useRouter();
  const methods = useForm<RecipeCreationFormValues>({
    mode: "onTouched",
    defaultValues: {
      level: LEVELS.EASY,
      time: TIMES.JUST_MINUTES,
      ingredients: [{ description: "" }],
      instructions: [{ description: "" }],
    },
    resolver: zodResolver(recipeCreationValidateSchema),
  });

  const [isError, setIsError] = useState(false);

  const { mutate: createRecipeMutation, isPending: isCreateRecipePending } =
    useCreateRecipe({
      onError: () => {
        setIsError(true);
      },
      onSuccess: (response) => {
        if (response) {
          router.push("/");
        }
      },
    });

  const getLevelId = (level: LEVELS) => {
    switch (level) {
      case LEVELS.HARD:
        return "hard";
      case LEVELS.MEDIUM:
        return "medium";
      case LEVELS.EASY:
      default:
        return "easy";
    }
  };

  const getTimeId = (time: TIMES) => {
    switch (time) {
      case TIMES.MORE_THAN_HOUR:
        return "long";
      case TIMES.ABOUT_HOUR:
        return "60m";
      case TIMES.HALF_HOUR:
        return "30m";
      case TIMES.JUST_MINUTES:
      default:
        return "10m";
    }
  };

  const mapRecipePayload = (
    data: RecipeCreationFormValues,
  ): IRecipeCreationPayload => {
    const { level, time, imageUrl, ...restData } = data;
    return {
      ...restData,
      imageUrl: imageUrl ? imageUrl : undefined,
      difficultyId: getLevelId(level),
      durationId: getTimeId(time),
    };
  };

  const handleSubmit = (data: RecipeCreationFormValues) => {
    const recipePayload = mapRecipePayload(data);
    createRecipeMutation(recipePayload);
  };

  return (
    <div className="px-4 py-8">
      {isError && (
        <div
          className={
            "bg-destructive/20 text-destructive-strong px-6 py-3 w-full"
          }
        >
          Recipe Creation Error
        </div>
      )}

      <h1 className="wongnok-text-h2">Create Recipe</h1>
      <p>description </p>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(handleSubmit)}>
          <RecipeTheDishForm />
          <RecipeEffortForm />
          <RecipeIngredientForm />
          <RecipeInstructionForm />
          <Button type="submit" disabled={isCreateRecipePending}>
            Submit
          </Button>
        </form>
      </FormProvider>
    </div>
  );
};

export default RecipeCreationPage;
