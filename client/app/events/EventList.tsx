"use client"

import { FC } from "react"

interface Event {
    id: number;
    title: string;
    organisation_name: string;
    description: string;
    image_url?: string;
}

interface EventListProps {
    events: Event[];
}

// Option 1: export default directly
const EventList: FC<EventListProps> = ({ events }) => {
    return (
        <div className="event-list grid grid-cols-4 gap-4 ">
            {events.map(event => (
                <a className="block" href={`/events/${event.id}`} key={event.id}>
  <div
    className="max-w-sm h-66 overflow-hidden
               hover:bg-primary hover:text-primary-foreground
               transition-colors duration-300 ease-in-out
               flex flex-col
               border-2 border-dashed border-gray-400 "
  >
    {/* Image */}
    {event.image_url && (
      <img className="w-full h-38 object-cover" src={event.image_url} />
    )}

    {/* Title */}
    <div className="px-6 py-2 flex-1">
      <div className="text-lg mb-1 line-clamp-2">
        {event.title}
      </div>
    </div>

    {/* Organisation tag */}
    <div className="px-6 pt-1 pb-2">
      <span className="inline-block bg-gray-200 px-3 py-1 text-sm font-semibold text-gray-700 truncate">
        {event.organisation_name}
      </span>
    </div>
  </div>
</a>

            ))}
        </div>
    )
}

export default EventList
