"use client";

import { Button } from "@/components/bases";

import { Product, products } from "./_shared/products";
import { addItem } from "@/lib/redux/slices/shopingSlice";
import { useAppDispatch } from "@/lib/redux/hooks";

const ProductList = () => {

  const dispatch = useAppDispatch()

  const handleAddToCart = (productId: Product) => {
    dispatch(addItem(productId))
  };

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
      {products.map((product) => (
        <div
          key={product.id}
          className="flex flex-col gap-4 rounded-lg border border-border p-4"
        >
          <div>
            <h3 className="wongnok-text-h3">{product.name}</h3>
            <p className="wongnok-text-base text-muted-foreground">
              ฿{product.price.toLocaleString()}
            </p>
          </div>
          <Button size="small" onClick={() => handleAddToCart(product)}>
            + Add to cart
          </Button>
        </div>
      ))}
    </div>
  );
};

export default ProductList;
