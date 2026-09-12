export enum LEVELS {
  EASY = "EASY",
  MEDIUM = "MEDIUM",
  HARD = "HARD",
}

export type RecipeCreationFormValues = {
  name: string;
  description: string;
  imageUrl?: string;
  level: LEVELS;
  time: string;
  ingredients: { description: string }[];
  instructions: { description: string }[];
};
