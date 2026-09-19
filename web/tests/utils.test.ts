const sum = (a: number, b: number) => {
  return a + b;
};

// Arrange
// Act
// Assert

// story
describe("Validate sum function", () => {
  // use case
  it("Should return 3 when 1 + 2", () => {
    // Act
    const result = sum(1, 2);
    // Assert
    expect(result).toEqual(3);
  });
});
