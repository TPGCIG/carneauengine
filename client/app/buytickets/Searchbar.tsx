"use client"

import { Input } from "@/components/ui/input"
import { MagnifyingGlassIcon } from "@radix-ui/react-icons";
import React from "react"
import { useState } from "react"
import { EventList } from ""

export function SearchForEvent( { onSearch } ) {
        const [value, setValue] = useState("");

	function handleInput(e) {
		const v = e.target.value;
		setValue(v)
		onSearch(v);
	}
	  return (
	    <div className="w-full max-w-6xl mx-auto mt-8 px-4 ">
	      <div className="relative ">
		{/* Magnifying glass */}
		<MagnifyingGlassIcon className="absolute left-6 top-1/2 -translate-y-1/2 w-12 h-12 text-gray-400 pointer-events-none" />

		{/* Input */}
		<input
		  value = {value}
		  onChange={handleInput}
		  type="text"
		  placeholder="Search events..."
		  className="w-full pl-20 pr-6 py-6 text-2xl rounded-full border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary"
		/>
	      </div>
	</div>
   )
}
