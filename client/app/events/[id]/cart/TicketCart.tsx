"use client";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useState, useEffect } from "react";
import { ShoppingCart } from 'lucide-react';
import { Separator } from "@/components/ui/separator";

interface TicketType {
  id: number;
  name: string;
  price: number;
}

interface TicketCartProps {
  ticketSelection: Record<number, number>; // id -> quantity
  onCheckout: () => void;
}

export default function TicketCart({ ticketSelection, onCheckout }: TicketCartProps) {
  const [ticketTypes, setTicketTypes] = useState<Record<number, TicketType>>({});

  // fetch ticket types
  useEffect(() => {
    const ticketIds = Object.keys(ticketSelection).map(Number);

    if (ticketIds.length === 0) return; // no items, no fetch

    fetch("http://localhost:8080/api/ticketTypes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ticketIds })
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch");
        return res.json();
      })
      .then((data: Record<number, TicketType>) => setTicketTypes(data));

  }, [ticketSelection]); // must refetch when selection changes

  const getTotalPrice = () => {
    let total = 0;

    for (const ticketIdStr of Object.keys(ticketSelection)) {
      const ticketId = Number(ticketIdStr);
      const quantity = ticketSelection[ticketId];
      const ticket = ticketTypes[ticketId];

      if (!ticket) continue; // still loading

      total += ticket.price * quantity;
    }

    return total;
  };

  return (
    <Card className="w-96 p-4 rounded-none">
      <h2 className="flex items-center text-lg font-semibold mb-4 space-x-2">
        <ShoppingCart className="w-5 h-5" />
        <span>Cart Summary</span>
      </h2>

      <CardContent className="space-y-2">
        {Object.keys(ticketSelection).map((idStr) => {
          const ticket = ticketTypes[Number(idStr)];
          if (!ticket) return null; // still loading

          const ticketName = ticket.name;
          const price = ticket.price;
          const quantity = ticketSelection[Number(idStr)]

          return (
            <div className="grid grid-cols-5"key={Number(idStr)}>
              <div className="col-span-4">{quantity} × {ticketName}</div> <div className="justify-end">${(price * quantity).toFixed(2)} </div>
            </div>
          );
        })}
        <br />
        <div className="grid grid-cols-5 font-semibold">
              <div className="col-span-4">Total:</div> <div className="justify-end">${getTotalPrice().toFixed(2)}</div>
            </div>

        <Separator />
      </CardContent>

      <CardFooter className="flex justify-between items-center mt-4">
        <span className="">
          
        </span>
        <Button onClick={onCheckout}>Checkout with Stripe</Button>
      </CardFooter>
    </Card>
  );
}
