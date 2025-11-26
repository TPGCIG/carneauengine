"use client"

import TicketCart from "./TicketCart"
import React from "react";
import { useEffect, useState } from "react"
import { UserInfoForm } from "./UserInfoForm"

interface TicketType {
  name: string;
  price: number;
}

function CheckoutPage() {
    const [ticketSelection, setTicketSelection] = useState<Record<number, number>>({});
    const [email, setEmail] = useState<string>("");
    const [phoneNumber, setPhoneNumber] = useState<string>("");

    useEffect(() => {
        const selection = sessionStorage.getItem("ticketSelection");
        if (selection) setTicketSelection(JSON.parse(selection));
    }, []);

    function handleCheckout() {

    }



    return (
        <div className="grid grid-cols-2 gap-4">
            <div className="flex justify-end">
                <UserInfoForm />
            </div>

            <div>
                <TicketCart ticketSelection={ticketSelection} onCheckout={handleCheckout}/>
            </div>
      </div>

    )

}

export default CheckoutPage;