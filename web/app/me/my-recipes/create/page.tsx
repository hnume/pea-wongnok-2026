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
} from "@/types/FormValues/recipeCreationForm";
import { zodResolver } from "@hookform/resolvers/zod";
import z4 from "zod/v4";

const recipeCreationValidateSchema = z4.object({
  name: z4.string().nonempty({ error: "กรุณากรอกชื่อเมนูอาหาร" }),
  description: z4.string().nonempty({ error: "กรุณากรอกรายละเอียดเมนูอาหาร" }),
  imageUrl: z4.url().optional(),
  level: z4.enum(LEVELS, {
    error: "Please select level of recipe",
  }),
  time: z4.string(),
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
  const methods = useForm<RecipeCreationFormValues>({
    mode: "onTouched",
    defaultValues: {
      level: LEVELS.EASY,
      time: "JUST_MINUTES",
      ingredients: [{ description: "" }],
      instructions: [{ description: "" }],
    },
    resolver: zodResolver(recipeCreationValidateSchema),
  });

  const handleSubmit = (data: unknown) => {
    console.log("SUBMITTED", data);
  };

  return (
    <div className="px-4 py-8">
      <h1 className="wongnok-text-h2">Create Recipe</h1>
      <p>description </p>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(handleSubmit)}>
          <RecipeTheDishForm />
          <RecipeEffortForm />
          <RecipeIngredientForm />
          <RecipeInstructionForm />
          <Button type="submit">Submit</Button>
        </form>
      </FormProvider>
    </div>
  );
};

export default RecipeCreationPage;
