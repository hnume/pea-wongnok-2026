"use client";

import { Button } from "@/components/bases";

import { products } from "./_shared/products";

const ProductList = () => {
  const handleAddToCart = (productId: string) => {
    // TODO: dispatch addItem
    void productId;
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
          <Button size="small" onClick={() => handleAddToCart(product.id)}>
            + Add to cart
          </Button>
        </div>
      ))}
    </div>
  );
};

export default ProductList;
