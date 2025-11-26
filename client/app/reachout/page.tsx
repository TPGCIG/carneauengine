// app/not-found.tsx
"use client";

import Link from "next/link";

export default function ReachOut() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center">
      <h1 className="text-6xl font-bold mb-4">This page is under construction!</h1>
      <p className="text-xl mb-8">Give us a sec.</p>
      <Link href="/" className="text-blue-600 underline hover:text-blue-800">
        Go back home
      </Link>
    </div>
  );
}
