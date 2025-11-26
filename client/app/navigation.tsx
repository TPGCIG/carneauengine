"use client"

import React from "react";
import Link from "next/link"
import Image from "next/image";
import logo from "@/public/logo.png"; // Ensure this path is correct
import { 
  Ticket, 
  CalendarDays, 
  Sparkles, 
  BookOpen, 
  LogIn 
} from "lucide-react"; // Install lucide-react if you haven't

export function Navigation() {
  return (
    <>
      {/* 1. Outer Wrapper 
        - h-14: Fixed height like the source
        - border-b: The subtle bottom border
        - bg-background: Solid background (add /95 for slight transparency)
      */}
      <nav className="fixed top-0 left-0 right-0 z-50 h-14 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        
        {/* 2. Inner Container 
          - max-w-[1550px]: The exact width from the source
          - mx-auto: Centers the wide container
        */}
        <div className="mx-auto flex h-full w-full max-w-[1550px] items-center justify-between px-4 sm:px-12 lg:px-32 xl:px-40">
          
          {/* Left: Logo Area */}
          <div className="flex items-center">
             <Link href="/" className="flex items-center gap-2">
                <Image 
                  src={logo} 
                  alt="Carneau Logo" 
                  placeholder="blur" 
                  className="h-8 w-auto" 
                />
             </Link>
          </div>

          {/* Middle: Links (The "Greptile" Look) */}
          {/* hidden on mobile (lg:flex), flex-row, monospace font */}
          <div className="hidden flex-row items-center lg:flex">
            
            {/* Link 1 */}
            <NavLink href="/events" icon={<Ticket className="h-3.5 w-3.5"/>}>
              Buy Tickets
            </NavLink>

            {/* Separator */}
            <DashedSeparator />

            {/* Link 2 */}
            <NavLink href="/hosting" icon={<CalendarDays className="h-3.5 w-3.5"/>}>
              Host an Event
            </NavLink>

            {/* Separator */}
            <DashedSeparator />

            {/* Link 3 */}
            <NavLink href="/firstevent" icon={<Sparkles className="h-3.5 w-3.5"/>}>
              Big First Event
            </NavLink>

            {/* Separator */}
            <DashedSeparator />

            {/* Link 4 */}
            <NavLink href="/docs" icon={<BookOpen className="h-3.5 w-3.5"/>}>
              Docs
            </NavLink>

          </div>

          {/* Right: Sign In Action */}
          <div className="flex items-center gap-4">
             {/* The "Start Now" / Sign In button style */}
             <Link href="/hosting">
                <button className="inline-flex h-9 items-center justify-center gap-2 whitespace-nowrap bg-secondary px-4 py-2 text-sm font-mono text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50">
                  <span>Sell Tickets</span>
                </button>
             </Link>
             <Link href="/events">
                <button className="inline-flex h-9 items-center justify-center gap-2 whitespace-nowrap bg-primary px-4 py-2 text-sm font-mono text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50">
                  <span>Buy Tickets</span>
                </button>
             </Link>
          </div>

        </div>
      </nav>
      
      {/* Spacer */}
      <div className="h-14" />
    </>
  )
}

// --- Helper Components to keep code clean ---

// 1. The Dashed Vertical Line
function DashedSeparator() {
  return (
    <div className="mx-2 h-4 w-px border-r border-dashed border-border opacity-50" />
  );
}

// 2. The Link Item Style
function NavLink({ href, icon, children }: { href: string; icon: React.ReactNode; children: React.ReactNode }) {
  return (
    <Link 
      href={href} 
      className="flex items-center gap-1.5 text-sm font-mono text-muted-foreground transition-colors hover:text-foreground"
    >
      {icon}
      <span>{children}</span>
    </Link>
  )
}