"use client"

import TicketCart from "./TicketCart"
import React from "react";
import { useEffect, useState } from "react"
import { UserInfoForm } from "./UserInfoForm"
import { E164Number } from "libphonenumber-js/core";

interface TicketType {
  name: string;
  price: number;
}

function CheckoutPage() {
    const [ticketSelection, setTicketSelection] = useState<Record<number, number>>({});
    const [email, setEmail] = useState<string>("");
    const [phone, setPhone] = useState<E164Number | undefined>();

    useEffect(() => {
        const selection = sessionStorage.getItem("ticketSelection");
        if (selection) setTicketSelection(JSON.parse(selection));
    }, []);

    async function handleCheckout() {
        const items = Object.entries(ticketSelection).map(([id, quantity]) => ({
            ticket_id: Number(id),
            quantity: quantity,
        }));

        if (items.length === 0) {
            alert("Cart is empty");
            return;
        }

        if (!email) {
            alert("Please enter your email");
            return;
        }

        try {
            const response = await fetch("http://localhost:8080/create-checkout-session", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    items,
                    email,
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || "Checkout failed");
            }

            const data = await response.json();
            if (data.url) {
                window.location.href = data.url;
            } else {
                console.error("No checkout URL returned");
            }

        } catch (error) {
            console.error("Error during checkout:", error);
            alert("Checkout failed. Please try again.");
        }
    }

    return (
        <div className="grid grid-cols-2 gap-4 pt-10">
            <div className="flex justify-end">
                <UserInfoForm email={email} setEmail={setEmail} phone={phone} setPhone={setPhone} />
            </div>

            <div>
                <TicketCart ticketSelection={ticketSelection} onCheckout={handleCheckout}/>
            </div>
      </div>

    )
}

export default CheckoutPage;