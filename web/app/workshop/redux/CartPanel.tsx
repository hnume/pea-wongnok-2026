"use client";

import { Button } from "@/components/bases";
import { useAppSelector } from "@/lib/redux/hooks";

const CartPanel = () => {
  const shoppingList = useAppSelector(state => state.shoppingSlice)

  const items = shoppingList.items;
  const total = items.reduce(
    (sum, item) => sum + item.price * item.quantity,
    0,
  );

  const handleDecrease = (id: string) => {
    // TODO: dispatch decreaseQuantity
    void id;
  };
  const handleIncrease = (id: string) => {
    // TODO: dispatch increaseQuantity
    void id;
  };
  const handleRemove = (id: string) => {
    // TODO: dispatch removeItem
    void id;
  };
  const handleClear = () => {
    // TODO: dispatch clearCart
  };

  // Swap to this when the cart has no items (items.length === 0):
  // if (items.length === 0) {
  //   return (
  //     <div className="rounded-lg border border-border p-4">
  //       <h2 className="wongnok-text-h3">Your Cart</h2>
  //       <p className="mt-4 text-muted-foreground">Cart is empty.</p>
  //     </div>
  //   );
  // }

  return (
    <div className="flex flex-col gap-4 rounded-lg border border-border p-4">
      <h2 className="wongnok-text-h3">Your Cart</h2>

      <ul className="flex flex-col gap-4">
        {items.map((item) => (
          <li
            key={item.id}
            className="flex flex-col gap-2 border-b border-border pb-4 last:border-b-0 last:pb-0"
          >
            <div>
              <p className="wongnok-text-body">{item.name}</p>
              <p className="wongnok-text-sm text-muted-foreground">
                ฿{item.price.toLocaleString()} × {item.quantity}
              </p>
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outlined"
                size="small"
                onClick={() => handleDecrease(item.id)}
              >
                -
              </Button>
              <span className="wongnok-text-sm">{item.quantity}</span>
              <Button
                variant="outlined"
                size="small"
                onClick={() => handleIncrease(item.id)}
              >
                +
              </Button>
              <Button
                variant="text"
                color="error"
                size="small"
                className="ml-auto"
                onClick={() => handleRemove(item.id)}
              >
                Remove
              </Button>
            </div>
          </li>
        ))}
      </ul>

      <div className="flex items-center justify-between border-t border-border pt-4">
        <span className="wongnok-text-body">Total</span>
        <span className="wongnok-text-h3">฿{total.toLocaleString()}</span>
      </div>

      <Button variant="outlined" color="error" onClick={handleClear}>
        Clear cart
      </Button>
    </div>
  );
};

export default CartPanel;
