import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Navigation } from '@/app/navigation'

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Carneau",
  description: "Ticket sales",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com"/>
          <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
          <link href="https://fonts.googleapis.com/css2?family=Geist+Mono:wght@100..900&family=TASA+Orbiter:wght@400..800&display=swap" rel="stylesheet"/>
      </head>
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased min-h-screen flex flex-col`}>
        
        <Navigation />

        {/* Main content */}
        <div className="flex-1 relative mx-auto w-full max-w-7xl px-10
          before:absolute before:top-0 before:bottom-0 before:left-0 before:w-px before:bg-gray-300
          after:absolute after:top-0 after:bottom-0 after:right-0 after:w-px after:bg-gray-300">

          <div className="">
            {children}
          </div>

        </div>

        {/* Footer inside body */}
        <footer className="bg-gray-800 text-white p-4 text-center">
          © 2025 Carneau
        </footer>

      </body>
    </html>
  );
}
