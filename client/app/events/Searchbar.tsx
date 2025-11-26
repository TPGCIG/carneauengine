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
  <div className="w-full max-w-4xl mx-auto mt-6 px-3">
    <div className="relative">
      {/* Magnifying glass */}
      <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-8 h-8 text-gray-400 pointer-events-none" />

      {/* Input */}
      <input
        value={value}
        onChange={handleInput}
        type="text"
        placeholder="Search events..."
        className="w-full pl-14 pr-4 py-4 text-xl border border-gray-300 bg-gray-100 focus:outline-none focus:ring-2 focus:ring-primary"
      />
    </div>
  </div>
)
}
