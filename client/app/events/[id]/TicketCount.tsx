'use client';

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface TicketTypeRowProps {
  ticketType: string;
  price: number;
  count: number; // controlled value
  min?: number;
  max?: number;
  onQuantityChange: (quantity: number) => void;
}

export function TicketTypeRow({
  ticketType,
  price,
  count,
  min = 0,
  max = 10,
  onQuantityChange,
}: TicketTypeRowProps) {

  const handleChange = (newCount: number) => {
    if (newCount < min) newCount = min;
    if (newCount > max) newCount = max;
    onQuantityChange(newCount);
  };

  return (
    <div className="flex items-center justify-between p-4 border rounded-md">
      <div>
        <div className="font-semibold">{ticketType}</div>
        <div className="text-sm text-gray-500">${price.toFixed(2)}</div>
      </div>
      <div className="flex items-center space-x-2">
        <Button variant="outline" size="sm" onClick={() => handleChange(count - 1)}>–</Button>
        <Input
          type="number"
          value={count}
          onChange={(e) => handleChange(Number(e.target.value))}
          className="w-16 text-center"
          min={min}
          max={max}
        />
        <Button variant="outline" size="sm" onClick={() => handleChange(count + 1)}>+</Button>
      </div>
    </div>
  );
}