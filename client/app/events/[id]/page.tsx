"use client";

import * as React from "react";
import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

import EventDisplay from "@/app/events/[id]/EventInfo";
import { TicketTypeRow } from "@/app/events/[id]/TicketCount";

interface TicketType {
  id: number;
  name: string;
  price: number;        // in dollars
  total_quantity: number;
}

interface Event {
  id: number;
  organisation_name: string;
  title: string;
  description: string;
  location: string;
  start_time: string;
  end_time: string;
  total_capacity: number;
  image_urls?: string[];
  ticket_types: TicketType[];
}

interface TicketSelection {
  [ticketId: number]: number;
}

const EventPage = () => {
  const params = useParams(); // { id: "X" }
  const eventId = Number(params.id);

  const [loading, setLoading] = useState<boolean>(true);
  const [eventInfo, setEventInfo] = useState<Event | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Track selected quantities for each ticket type
  const [ticketSelection, setTicketSelection] = useState<TicketSelection>({});

  const router = useRouter();

  useEffect(() => {
    setLoading(true);
    fetch(`http://localhost:8080/api/events/${eventId}`)
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch");
        return res.json() as Promise<Event>;
      })
      .then((data) => setEventInfo(data))
      .catch((err: any) => setError(err.message))
      .finally(() => setLoading(false));
  }, [eventId]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!eventInfo) return <div>No event found.</div>;
  if (!eventInfo.image_urls) return <p>Error loading event data</p>;

  // Update parent state when ticket quantity changes
  const handleQuantityChange = (ticketId: number, quantity: number) => {
    setTicketSelection((prev) => ({ ...prev, [ticketId]: quantity }));
  };

  async function handleSubmit() {
    console.log(ticketSelection);
    if (!eventInfo) return;
    if (!ticketSelection) return;

    sessionStorage.setItem("ticketSelection", JSON.stringify(ticketSelection));
    router.push(`/events/${eventId}/cart`)




    // try {
    //   const res = await fetch("localhost:8080/checkout/create-checkout-session", {
    //     method: "POST",
    //     headers:  { "ContentType": "application/json" }, 
    //     body: JSON.stringify({
    //       eventId: eventInfo.id,
    //       tickets: ticketSelection
    //     }),
    //   });

    //   const data = await res.json();

    //   if (!res.ok) {
    //     throw new Error(data.error || "Failed to create checkout session");
    //   }

    //   window.location.href = data.url;
    // } catch (err: any) {
    //   console.log("Checkout error: ",err.message);
    //   alert("Failed to start checkout session. Please try again later.");
    // }
  };

  return (
    <div className="mx-auto">
      <Carousel className="w-9/10 mx-auto">
        <CarouselContent>
          {eventInfo.image_urls.length === 0 && <p>No images for this event</p>}
          {eventInfo.image_urls.map((image_url, index) => (
            <CarouselItem key={index}>
              <div className="p-1">
                <img src={image_url} alt={`Event image ${index + 1}`} />
              </div>
            </CarouselItem>
          ))}
        </CarouselContent>
        <CarouselPrevious />
        <CarouselNext />
      </Carousel>

      <Separator className="my-8" />

      <div className="grid grid-cols-3 gap-4">
        <div className="col-span-2">
          <EventDisplay event={eventInfo} />
        </div>

        <div>
          <Sheet>
            <SheetTrigger asChild>
              <div className="relative w-full max-w-xs mx-auto">
                {/* Outer outline slightly bigger than button */}
                <span className="absolute -inset-3 border border-gray-300 rounded-md"></span>
                <button className="relative w-full bg-purple-500 text-white rounded-md px-6 py-5 font-semibold hover:bg-purple-600 transition">
                  Buy Tickets
                </button>
              </div>
            </SheetTrigger>

            <SheetContent>
              <SheetHeader>
                <SheetTitle>Ticket Selection</SheetTitle>
                <SheetDescription>
                  Choose which tickets you would like to purchase.
                </SheetDescription>
              </SheetHeader>

              <div className="grid flex-1 auto-rows-min gap-6 px-4">
                {/* Example ticket type */}
                {eventInfo.ticket_types.map((ticket) => (
                    <TicketTypeRow key={ticket.id}
                    ticketType={ticket.name}
                    price={ticket.price}
                    count={ticketSelection[ticket.id] || 0}
                    onQuantityChange={(qty) => handleQuantityChange(ticket.id, qty)}
                    />
                ))}
              </div>

              <SheetFooter>
                <Button type="submit" onClick={handleSubmit}>
                  Checkout ({Object.values(ticketSelection).reduce((a, b) => a + b, 0)}{" "}
                  tickets)
                </Button>
                <SheetClose asChild>
                  <Button variant="outline">Close</Button>
                </SheetClose>
              </SheetFooter>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </div>
  );
};

export default EventPage;
