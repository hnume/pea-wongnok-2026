export type Product = {
  id: string;
  name: string;
  price: number;
};

export const products: Product[] = [
  { id: "p1", name: "Wireless Mouse", price: 490 },
  { id: "p2", name: "Mechanical Keyboard", price: 2490 },
  { id: "p3", name: "USB-C Hub", price: 890 },
  { id: "p4", name: "Desk Lamp", price: 690 },
];
