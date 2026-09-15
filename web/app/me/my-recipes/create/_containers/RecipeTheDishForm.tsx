import { Textarea, TextField } from "@/components/bases";
import { RecipeCreationFormValues } from "@/types/FormValues/recipeCreationForm";
import { useFormContext } from "react-hook-form";

const RecipeTheDishForm = () => {
  const {
    register,
    formState: { errors },
  } = useFormContext<RecipeCreationFormValues>();

  return (
    <div className="bg-white p-5 mt-4">
      <p className="wongnok-text-h3">The Dish</p>
      <p className="wongnok-text-sm text-muted-foreground">
        {`Give it a name people will recognise, and a line about why it's good.`}
      </p>
      <div className="flex flex-col gap-4 mt-4">
        <TextField
          {...register("name")}
          label={"Menu Name"}
          placeholder={"e.g. Thai basil chicken with a crisp fried egg"}
          error={!!errors.name?.message}
          errorMessage={errors.name?.message}
          required
        />
        <Textarea
          {...register("description")}
          label={"Menu Description"}
          placeholder={
            "Two or three sentences — what it tastes like, when you cook it, any shortcut you love."
          }
          error={!!errors.description?.message}
          errorMessage={errors.description?.message}
          required
        />
        <TextField
          {...register("imageUrl")}
          label={"Image URL"}
          placeholder={"https://…/my-dish.jpg"}
          helperText={"Paste a link to a photo — landscape works best."}
          error={!!errors.imageUrl?.message}
          errorMessage={errors.imageUrl?.message}
        />
      </div>
    </div>
  );
};

export default RecipeTheDishForm;
