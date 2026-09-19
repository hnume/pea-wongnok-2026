import { use } from "react";
import RecipeDetail from "./_containers/RecipeDetail";

export default function Page({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } =  use(params);

  return <RecipeDetail id={id} />;
}
