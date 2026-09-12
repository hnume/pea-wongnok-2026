import { Radio, RadioGroup, Select } from "@/components/bases";
import { LEVELS } from "@/types/FormValues/recipeCreationForm";
import { Controller, useFormContext } from "react-hook-form";

const recipeLevelOptions = [
  { label: "Easy", value: LEVELS.EASY },
  { label: "Medium", value: LEVELS.MEDIUM },
  { label: "Hard", value: LEVELS.HARD },
];

const RecipeEffortForm = () => {
  const { register, setValue, getValues, control } = useFormContext();

  return (
    <div className="mt-6">
      <div className="bg-white p-7 rounded-3xl">
        <p className="font-bold">Effort</p>
        <p className="wongnok-text-body text-muted-foreground">
          {`Helps cooks pick something that fits their evening.`}
        </p>
        <div className="flex w-full gap-6 mt-6">
          <div className={"flex-1 shrink-0"}>
            <Controller
              name={"level"}
              control={control}
              render={({ field }) => (
                <Select
                  {...field}
                  label={"Menu Level of Recipe"}
                  options={recipeLevelOptions}
                  required
                  placeholder={"Select level of recipe"}
                  onValueChange={(value) => {
                    setValue("level", value);
                  }}
                />
              )}
            />
          </div>
          <div className={"flex-1 shrink-0"}>
            <Controller
              name="time"
              control={control}
              render={({ field }) => (
                <RadioGroup
                  label={"Time to Make"}
                  variant={"outlined"}
                  className={"grid grid-cols-2"}
                  {...field}
                  required
                >
                  <Radio label={"5 - 10 mins"} value={"JUST_MINUTES"} />
                  <Radio label={"10 - 30 mins"} value={"HALF_HOUR"} />
                  <Radio label={"~1 hour"} value={"ABOUT_HOUR"} />
                  <Radio label={"More than 1 hour"} value={"MORE_THAN_HOUR"} />
                </RadioGroup>
              )}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default RecipeEffortForm;
