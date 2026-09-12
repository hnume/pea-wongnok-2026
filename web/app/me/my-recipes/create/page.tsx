"use client";

import { Button } from "@/components/bases";
import RecipeTheDishForm from "./_containers/RecipeTheDishForm";
import RecipeEffortForm from "./_containers/RecipeEffortForm";
import RecipeIngredientForm from "./_containers/RecipeIngredientForm";
import RecipeInstructionForm from "./_containers/RecipeInstructionForm";
import { FormProvider, useForm } from "react-hook-form";
import { RecipeCreationFormValues } from "@/types/FormValues/recipeCreationForm";

const RecipeCreationPage = () => {
  const methods = useForm<RecipeCreationFormValues>({
    defaultValues: {
      level: "EASY",
      time: "JUST_MINUTES",
      ingredients: [{ description: "" }],
      instructions: [{ description: "" }],
    },
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
