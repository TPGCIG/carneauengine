"use client"

import { FC } from "react"

interface Event {
    id: number;
    title: string;
    organisation_id: string;
    description: string;
    image_url?: string;
}

interface EventListProps {
    events: Event[];
}

// Option 1: export default directly
const EventList: FC<EventListProps> = ({ events }) => {
    return (
        <div className="event-list">
            {events.map(event => (
                <div key={event.id} className="max-w-sm rounded overflow-hidden shadow-lg">
                    <img className="w-full" src={event.image_url} alt={event.title} />
                    <div className="px-6 py-4">
                        <div className="font-bold text-xl mb-2">{event.title}</div>
                        <p className="text-gray-700 text-base">{event.description}</p>
                    </div>
                    <div className="px-6 pt-4 pb-2">
                        <span className="inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-gray-700 mr-2 mb-2">H</span>
                    </div>
                </div>
            ))}
        </div>
    )
}

export default EventList
