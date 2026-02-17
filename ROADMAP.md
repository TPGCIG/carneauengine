# Project Roadmap: Carneau Engine - Ticket Sales

This document outlines the architectural overview, current state, and a detailed roadmap to complete the ticket sales engine, including steps for production readiness.

## Architectural Overview

The project consists of two main components:

1.  **Frontend (Client):** A Next.js application (TypeScript, Tailwind CSS, shadcn/ui) responsible for the user interface, event browsing, cart management, and initiating the checkout process.
2.  **Backend (Server):** A Go application using the Gin web framework, responsible for handling API requests, interacting with the PostgreSQL database, and integrating with external services like Stripe.

### Key Components

*   **Database (PostgreSQL):**
    *   `users`: User authentication and profiles.
    *   `organisations` & `organisation_members`: Multi-tenancy for event organizers.
    *   `events` & `event_images`: Event details and associated media.
    *   `ticket_types`: Defines different ticket categories for an event (e.g., VIP, General Admission, Early Bird) with price and quantity.
    *   `purchases`: Records of user transactions, including payment status and Stripe reference.
    *   `tickets`: Individual tickets generated upon successful purchase, linked to a user, purchase, and ticket type, including a QR code for verification.
*   **Go Backend:**
    *   `main.go`: Entry point, sets up Gin router, CORS, database connection, and registers handlers.
    *   `db/connection.go`: Handles PostgreSQL database connection.
    *   `handlers/`: Contains API logic for events, checkout, tickets, Stripe, and users.
    *   `proto/ticket.proto`: Placeholder for potential gRPC or protobuf definitions.
    *   **Current Endpoints:**
        *   `GET /api/events`: Fetch summarized event listings.
        *   `GET /api/events/:id`: Fetch details for a specific event.
        *   `POST /api/ticketTypes`: Retrieve ticket types based on IDs.
        *   `POST /create-checkout-session`: Initiates a Stripe checkout session.
*   **Next.js Frontend:**
    *   `app/events/`: Event listing and detail views.
    *   `app/events/[id]/cart/`: Shopping cart and checkout flow.
    *   `components/ui/`: Reusable UI components.
    *   Integrates with `@stripe/react-stripe-js` and `@stripe/stripe-js` for frontend payment processing.

## Current State & Gaps

*   The core data model for events, tickets, purchases, and users is defined in `schema.sql`.
*   Basic API endpoints for fetching events and creating Stripe checkout sessions are implemented.
*   The frontend can display events and initiate a checkout.
*   **Major Gaps:**
    *   **Payment Fulfillment:** The system can initiate a Stripe checkout, but lacks the backend logic to confirm payment success, create individual tickets, and update inventory.
    *   **Stripe Webhook:** No endpoint or logic to receive and process Stripe webhooks (e.g., `checkout.session.completed`).
    *   **Ticket Generation & QR Codes:** No process to generate `tickets` records and their associated `qr_code` values after a successful purchase.
    *   **Inventory Management:** While `sold_quantity` exists, the logic for atomically updating it and preventing overselling is not fully implemented in the checkout flow.
    *   **User Authentication:** The `users` table exists, but robust authentication middleware and user session management are not evident.
    *   **User Experience:** No dedicated pages for viewing purchased tickets or receiving email confirmations.
    *   **Admin/Organizer Features:** No functionality for event creation/management or ticket validation.

---

## Roadmap to Completion

### Phase 1: Core Checkout & Order Fulfillment (Backend Focus)

This phase ensures the payment flow is secure and tickets are correctly issued.

*   **Step 1: Implement Stripe Webhook Handling.**
    *   Create a new Go handler (`/webhook` or similar) to receive and process Stripe webhook events.
    *   Verify webhook signatures for security.
    *   Focus on handling `checkout.session.completed` events.
*   **Step 2: Fulfill the Purchase.**
    *   Inside the webhook handler, after confirming payment:
        *   Update the `purchases` table (`payment_status = 'succeeded'`, `stripe_payment_id`).
        *   Create individual `tickets` records based on the `ticket_types` purchased in the session.
        *   Update the `ticket_types.sold_quantity` atomically for each purchased ticket type.
*   **Step 3: Generate QR Codes.**
    *   For each `ticket` created in Step 2, generate a unique, secure QR code string (e.g., UUID or a signed token) and store it in `tickets.qr_code`.
*   **Step 4: Add Pre-Checkout Inventory Validation.**
    *   Before creating a Stripe checkout session, implement robust checks to ensure that the requested `quantity` for each `ticket_type` is available (`total_quantity - sold_quantity >= requested_quantity`).
    *   **(Optional/Advanced) Temporary Ticket Reservation:** Consider using Redis to temporarily hold tickets for a user during their checkout process (e.g., for 15-30 minutes) to prevent multiple users from purchasing the last few tickets simultaneously.

### Phase 2: User Experience (Frontend & Backend Integration)

This phase focuses on enhancing the user journey post-purchase.

*   **Step 5: Implement User Authentication.**
    *   Choose an authentication strategy (e.g., JWT, session-based).
    *   Implement user registration, login, and secure session management.
    *   Add authentication middleware to protected backend routes.
*   **Step 6: Build a "My Tickets" Page (Frontend & Backend).**
    *   **Backend:** Create an API endpoint (`GET /api/users/:userId/tickets`) to fetch all tickets associated with a logged-in user.
    *   **Frontend:** Create a page where logged-in users can view a list of their purchased tickets, including event details and QR codes.
*   **Step 7: Send Email Confirmations.**
    *   After purchase fulfillment (triggered by the Stripe webhook), use an email service (see Production Readiness) to send a confirmation email to the user.
    *   Include purchase details and a way to access their tickets/QR codes (e.g., a link to the "My Tickets" page or attached QR code images).
*   **Step 8: Enhance Success and Cancel Pages (Frontend).**
    *   The `http://localhost:3000/success` page should display detailed purchase information retrieved from the backend.
    *   The `http://localhost:3000/cancel` page should provide clear guidance on what went wrong and how to retry.

### Phase 3: Organizer & Administration Features

This phase adds functionality for event hosts and operations.

*   **Step 9: Event Management Dashboard (Frontend & Backend).**
    *   **Backend:** Create API endpoints for creating, updating, and deleting events and ticket types, accessible only by authenticated organizers/admins.
    *   **Frontend:** Develop a dashboard interface for organizers to manage their events, view sales, and monitor ticket inventory.
*   **Step 10: Ticket Scanning/Verification System (Frontend & Backend).**
    *   **Backend:** Create a secure API endpoint (`POST /api/tickets/verify`) that accepts a QR code string. This endpoint should validate the ticket, check its status, and mark it as `redeemed`.
    *   **Frontend:** Develop a simple, mobile-friendly interface for event staff to scan QR codes (using a device's camera) and verify ticket validity.

---

## Production Readiness Tools & Practices

To ensure the Carneau Engine is robust, scalable, and secure in a production environment, consider incorporating the following:

1.  **Redis:**
    *   **Caching:** For frequently accessed data (e.g., event lists, event details).
    *   **Temporary Ticket Holds:** To manage short-term ticket reservations during checkout.
    *   **Session Store:** If using server-side sessions for user authentication.
    *   **Rate Limiting:** To protect your APIs from abuse.

2.  **Containerization (Docker):**
    *   Create `Dockerfile`s for both the Go backend and Next.js frontend to ensure consistent build and runtime environments.

3.  **Orchestration (Docker Compose / Kubernetes):**
    *   **Docker Compose:** For streamlined local development and testing of all services together.
    *   **Kubernetes (K8s):** For managing and scaling your containerized applications in production, providing high availability, auto-scaling, and self-healing capabilities.

4.  **Managed Database Service:**
    *   Deploy PostgreSQL to a managed service (e.g., AWS RDS, Google Cloud SQL, Azure Database for PostgreSQL) for enterprise-grade features like automated backups, replication, and scaling.

5.  **CI/CD Pipeline:**
    *   Set up Continuous Integration/Continuous Deployment (CI/CD) using tools like GitHub Actions, GitLab CI, or Jenkins. This will automate testing, building, and deploying your application changes.

6.  **Logging and Monitoring:**
    *   **Structured Logging:** Implement JSON logging in both backend and frontend for easier parsing and analysis.
    *   **Application Performance Monitoring (APM):** Tools like Prometheus/Grafana, Datadog, or New Relic to monitor application health, performance metrics, and resource utilization.
    *   **Error Tracking:** Integrate services like Sentry or Bugsnag for real-time error detection, reporting, and debugging.

7.  **Transactional Email Service:**
    *   Utilize a reliable email API service (e.g., SendGrid, Postmark, Resend, AWS SES) for sending critical transactional emails (ticket confirmations, password resets).

8.  **Environment Variable Management:**
    *   Securely store and inject environment variables into your production environment, avoiding `godotenv` for production secrets. Cloud providers offer dedicated secret management services.

9.  **Security Best Practices:**
    *   Enforce HTTPS/TLS for all communication.
    *   Implement robust input validation and sanitization.
    *   Regularly audit dependencies for vulnerabilities.
    *   Configure appropriate CORS policies.

This comprehensive roadmap should guide you through building a complete and production-ready ticket sales engine.
