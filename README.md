# Carneau Engine: Event Ticketing Platform

## Project Overview

Carneau Engine is a modern, full-stack event ticketing platform designed to provide a a robust and scalable solution for managing and selling event tickets. It features a responsive frontend for event browsing and purchasing, backed by a high-performance Go API, and leverages advanced techniques for ensuring data consistency and optimal user experience.

## Key Features

*   **Event Management:** Browse and view detailed information for various events.
*   **Ticket Sales:** Select and purchase tickets for desired events.
*   **Secure Payment Processing:** Seamless integration with Stripe for checkout sessions.
*   **Atomic Ticket Reservation:** Concurrency-safe system to prevent overselling of limited tickets.
*   **High Performance:** Caching mechanisms to optimise data delivery for frequently accessed event information.
*   **User Authentication:** (Planned) Secure user registration and login.
*   **Responsive Design:** Modern user interface built with Next.js and Tailwind CSS.

## Technical Architecture

The project is structured into a client-server architecture:

*   **Frontend (Client):** A Next.js application built with TypeScript, Tailwind CSS, and shadcn/ui, providing the user interface for event discovery, cart management, and payment initiation.
*   **Backend (Server):** A Go application utilising the Gin web framework, responsible for handling API requests, business logic, database interactions, and integrations with external services like Stripe.
*   **Database:** PostgreSQL serves as the primary data store for all persistent data including users, events, ticket types, purchases, and tickets.
*   **Caching & Concurrency Layer:** Redis is employed for critical performance enhancements and to ensure transactional integrity in high-concurrency scenarios.

## Core Implementations & Technical Highlights

*   **Atomic Ticket Reservation System:**
    *   Implemented a robust solution to prevent overselling of tickets, a common challenge in ticketing platforms.
    *   Utilises a combination of **PostgreSQL's `SELECT ... FOR UPDATE`** for pessimistic locking during initial availability checks and **atomic Redis Lua scripts** for temporary ticket holds during the checkout process.
    *   Ensures that requested tickets are reserved for a user for a set duration (e.g., 15 minutes) before payment, and are automatically released if the purchase is not completed, or permanently allocated upon successful Stripe payment.
*   **Event Data Caching:**
    *   Integrated Redis caching for read-heavy API endpoints (e.g., listing all events, fetching individual event details).
    *   Significantly reduces database load and improves API response times by serving cached data with a time-to-live (TTL), falling back to the database on cache misses.
    *   Includes cache invalidation strategies to ensure data freshness when events are updated.
*   **Stripe Payment Gateway Integration:**
    *   Implemented server-side creation of Stripe Checkout Sessions for secure and compliant payment processing.
    *   Developed a Stripe Webhook handler to asynchronously process payment outcomes (e.g., `checkout.session.completed`), update purchase statuses, generate unique tickets with QR codes, and decrement ticket inventory in the database.
*   **Modular Go API Design:** Structured handlers, models, and database interactions for maintainability and scalability.

## Technologies Used

**Frontend:**
*   Next.js (React)
*   TypeScript
*   Tailwind CSS
*   shadcn/ui
*   pnpm

**Backend:**
*   Go
*   Gin Web Framework
*   pgx (PostgreSQL driver)
*   Stripe Go Library
*   go-redis (Redis client)

**Databases & Services:**
*   PostgreSQL
*   Redis
*   Stripe (Payment Gateway)

## Setup and Local Development

To set up and run Carneau Engine locally, follow these steps:

### Prerequisites

*   Go (1.21 or newer)
*   Node.js (18 or newer)
*   pnpm (for Node.js package management)
*   PostgreSQL database server
*   Redis server
*   k6 (for load testing, optional)

### 1. Database Setup

1.  **Start PostgreSQL and Redis servers.**
2.  **Create a PostgreSQL database.** For example, `ticketing`.
3.  **Apply the database schema:**
    ```bash
    psql -U your_pg_user -d ticketing -f schema.sql
    ```
    (Replace `your_pg_user` and `ticketing` with your credentials/database name).

### 2. Environment Configuration

Create a `.env` file in the `server/` directory:
```env
# server/.env
DATABASE_URL="postgres://postgres:123@localhost:5432/ticketing?sslmode=disable"
REDIS_URL="redis://localhost:6379/0"
STRIPE_SECRET_KEY="sk_test_..." # Replace with your Stripe secret key
STRIPE_WEBHOOK_SECRET="whsec_..." # Replace with your Stripe webhook secret
SMTP_HOST="" # Optional: For email sending
SMTP_PORT=""
SMTP_USER=""
SMTP_PASSWORD=""
SENDER_EMAIL=""
```
*Remember to replace placeholder values with your actual credentials.*

### 3. Run the Backend

1.  Navigate to the `server/` directory:
    ```bash
    cd server
    ```
2.  Install Go dependencies:
    ```bash
    go mod tidy
    ```
3.  Build and run the server:
    ```bash
    make dev
    ```
    The API will be available at `http://localhost:8080`.

### 4. Run the Frontend

*(Note: At this point, the frontend might have TypeScript errors due to recent reverts, requiring fixes to run locally.)*

1.  Navigate to the `client/` directory:
    ```bash
    cd client
    ```
2.  Install Node.js dependencies:
    ```bash
    pnpm install
    ```
3.  Start the Next.js development server:
    ```bash
    pnpm dev
    ```
    The frontend will be available at `http://localhost:3000`.

## Deployment Strategy (Planned)

The planned deployment strategy for production involves a hybrid approach to leverage optimal services for each component:
*   **Frontend:** Deployed on Vercel for highly optimised Next.js hosting.
*   **Backend:** Deployed as a web service on an integrated PaaS like Render (for Go application).
*   **Databases:** Utilise free and persistent managed services like Neon.tech or Supabase for PostgreSQL, and Upstash for Redis.

## Future Enhancements (from ROADMAP.md)

*   **User Authentication:** Full implementation of user registration, login, and secure session management.
*   **"My Tickets" Page:** A dedicated user interface for viewing purchased tickets.
*   **Organiser Dashboard:** Features for event creators to manage events and ticket sales.
*   **Ticket Scanning System:** Functionality for event staff to verify tickets via QR codes.
*   **CI/CD Pipeline:** Automated testing, building, and deployment.

---
