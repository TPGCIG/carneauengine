"use client"

import  EventList  from '@/app/events/EventList'
import { SearchForEvent } from '@/app/events/Searchbar'
import { useState, useEffect } from "react"
import Fuse from "fuse.js"
import Link from "next/link";


export default function Home() {
	const [events, setEvents] = useState([]);
	const [filtered, setFiltered] = useState([]);


	useEffect(() => {
		fetch("http://localhost:8080/api/events")
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
		<div className=''>
			<div className='flex flex-col items-center pt-15'>
				<h1 className=''>Looking for events?</h1>
				<h1 className=''>We've got a few.</h1>
				<h2 className='text-gray-600'><br/>Don't see your event on the list? <Link href="/reachout" className='text-primary underline'>We can change that.</Link></h2>
			</div>
			<div className="flex items-center justify-center font-sans dark:primary py-8">
				<SearchForEvent onSearch={handleSearch}/>
			</div>
			<div>
				{filtered.length > 0 && <EventList events={filtered}/>}
				{filtered.length == 0 && <p>No events found</p>}
			</div>
			
		</div>
  );
}
