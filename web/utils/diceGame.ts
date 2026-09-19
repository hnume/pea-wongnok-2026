import { rollDice } from "./dice";

export const diceGame = () => {
  const result = rollDice();

  if (result === 1) return "You Win!";
  else return `You Lose, dice is ${result}`;
};
