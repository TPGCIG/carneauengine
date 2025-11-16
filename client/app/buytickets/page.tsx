"use client"

import  EventList  from '@/app/buytickets/EventList'
import { SearchForEvent } from '@/app/buytickets/Searchbar'
import { useState, useEffect } from "react"
import Fuse from "fuse.js"


export default function Home() {
	const [events, setEvents] = useState([]);
	const [filtered, setFiltered] = useState([]);


	useEffect(() => {
		fetch("http://localhost:8080/events")
		.then(res => res.json())
		.then(data => {
			setEvents(data);
			setFiltered(data);
		});
	}, []);

	const handleSearch = (query: string) => {
	    if (!query) {
		setFiltered(events);
		return;
	    }

	    const fuse = new Fuse(events, {
		keys: ["title", "description"], // make sure this matches your Event interface keys
		threshold: 0.4
	    });

	    const result = fuse.search(query).map(r => r.item); // <- note `.item`, not `.items`
	    setFiltered(result);
	};

	return (
    		<div className="flex min-h-screen items-center justify-center bg-neutral-900 font-sans dark:bg-black">
        		<SearchForEvent onSearch={handleSearch}/>
    
			{filtered.length > 0 && <EventList events={filtered}/>}
			{filtered.length == 0 && <p>No events found</p>}

		</div>
  );
}
