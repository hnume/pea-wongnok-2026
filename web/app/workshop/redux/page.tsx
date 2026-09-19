"use client";
import CartPanel from "./CartPanel";
import CartSummary from "./CartSummary";
import ProductList from "./ProductList";

const ReduxWorkshopPage = () => {
  return (
    <div className="px-6 pt-8 pb-10">
      <header className="flex items-center justify-between gap-4">
        <h1 className="wongnok-text-h2">Shopping Cart</h1>
        <CartSummary />
      </header>

      <div className="grid grid-cols-1 md:grid-cols-[2fr_1fr] gap-6 mt-6">
        <section>
          <ProductList />
        </section>
        <aside>
          <CartPanel />
        </aside>
      </div>
    </div>
  );
};

export default ReduxWorkshopPage;
