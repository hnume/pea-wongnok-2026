export enum LEVELS {
  EASY = "EASY",
  MEDIUM = "MEDIUM",
  HARD = "HARD",
}

export enum TIMES {
  JUST_MINUTES = "JUST_MINUTES",
  HALF_HOUR = "HALF_HOUR",
  ABOUT_HOUR = "ABOUT_HOUR",
  MORE_THAN_HOUR = "MORE_THAN_HOUR",
}

export type RecipeCreationFormValues = {
  name: string;
  description: string;
  imageUrl?: string;
  level: LEVELS;
  time: TIMES;
  ingredients: { description: string }[];
  instructions: { description: string }[];
};
