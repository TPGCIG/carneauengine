"use client"

import { FC, useEffect, useState } from "react"

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
}

interface EventDisplayProps {
    event: Event;
}

// Option 1: export default directly
function EventDisplay( { event } : EventDisplayProps ) {
    return (
        <div className="">
           <p className="hname">{event.title}</p>
           <p>{event.description}</p>
        </div>
    )
}

export default EventDisplay
