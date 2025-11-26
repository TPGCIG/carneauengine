// app/not-found.tsx
"use client";

import Link from "next/link";

export default function NotFound() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100">
      <h1 className="text-6xl font-bold mb-4">404</h1>
      <p className="text-xl mb-8">This page does not exist.</p>
      <Link href="/" className="text-blue-600 underline hover:text-blue-800">
        Go back home
      </Link>
    </div>
  );
}
