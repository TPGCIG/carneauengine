"use client"

import React from "react";


import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuIndicator,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  NavigationMenuViewport,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu"
import Link from "next/link"

export function Navigation() {

    return (
        <div className="w-full">
        <NavigationMenu className="relative mx-auto rounded-full">
            <NavigationMenuList className="flex-wrap rounded-full">
                <NavigationMenuItem className="rounded-full">
                    <div className="rounded-full">
                    <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
                        <Link href="/buytickets">Buy Tickets</Link>
                    </NavigationMenuLink></div>
                </NavigationMenuItem>
                <NavigationMenuItem>
                    <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
                        <Link href="/hosting">Host an Event</Link>
                    </NavigationMenuLink>
                </NavigationMenuItem>
                <NavigationMenuItem>
                    <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
                        <Link href="/firstevent">Big First Event</Link>
                    </NavigationMenuLink>
                </NavigationMenuItem>
                <NavigationMenuItem>
                    <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
                        <Link href="/docs">Docs</Link>
                    </NavigationMenuLink>
                </NavigationMenuItem>
            </NavigationMenuList>
        </NavigationMenu>
        </div>
    )

};