"use client"

import Link from "next/link";

import {ArrowRight} from 'lucide-react';
import Image from "next/image";

import { Separator } from "@radix-ui/react-separator";
import logo from "@/public/icon.png"; // Ensure this path is correct
import {DotLottieReact} from "@lottiefiles/dotlottie-react";

export default function Home() {
  return (
    <div className="min-h-screen items-center">
            <div className="relative min-h-[60vh] w-full pt-12 sm:pt-16 overflow-hidden">
              <div className="absolute left-4 sm:left-8 top-12 sm:top-16 flex flex-col items-start">
                <div>
                  <h1 className="">Tickets for Giving,<br/>Not Taking</h1>
                </div>
              </div>

              <div className="absolute left-4 sm:left-8 bottom-12 sm:bottom-16 flex flex-col items-start">
                <div>
                  <span className="">
                    <h2>KEEPING FUNDS INSIDE CLUBS</h2>
                  </span>
                  <span className="">
                    <h2>AND STUDENT'S WALLETS BY</h2>
                  </span>
                  <span className="">
                    <h2>CHARGING LESS FOR MORE</h2>
                  </span>
                  
                  
                  
                  <div className="h-8" />
                  <Link href="/events">
                      <button className="relative inline-flex w-full h-12 items-center justify-center gap-2 whitespace-nowrap bg-primary px-4 py-2 text-lg font-mono text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50">
                        <span className="inline-flex items-center gap-2">
                          This Is How <ArrowRight className="h-5 w-5"/>
                        </span>
                      </button>
                  </Link>
                </div>
              </div>

              <div className="absolute right-0 bottom-0 h-[250px] w-[250px] sm:h-[600px] sm:w-[600px] xl:h-[600px] xl:w-[500px] overflow-visible">

                <DotLottieReact
                  src="https://lottie.host/6f159aa8-84cd-4f4f-af32-d99a26643083/QXQ3FARZcl.lottie"
                  loop
                  autoplay
                  className="w-[500px] h-[500px]"
                />
              </div>
            </div>
            <div className="relative w-full">
    <hr className="border-t border-gray-300  w-auto -mx-10 mt-6" />



  </div>
      

    </div>
  );
}
