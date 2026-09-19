import { rollDice } from "@/utils/dice";
import { diceGame } from "@/utils/diceGame";

// Arrange
jest.mock("@/utils/dice.ts");
const mockedRollDiceFn = jest.mocked(rollDice);

describe("Validate dice game", () => {
  it('Should return "You Win!" when roll dice equal 1', () => {
    // Arrange
    mockedRollDiceFn.mockReturnValue(1);
    // Act
    const result = diceGame();
    // Assert
    expect(result).toEqual("You Win!");
  });

  it('Should return "You Lose, dice is <dice_face>" when roll dice not equal 1', () => {
    // Arrange
    mockedRollDiceFn.mockReturnValue(3);
    // Act
    const result = diceGame();
    // Assert
    expect(result).toEqual("You Lose, dice is 3");
  });
});
