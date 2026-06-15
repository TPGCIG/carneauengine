import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';

// Tell k6 that 409 is an expected outcome, not a failure.
// Without this, http_req_failed counts every 409 as an error,
// which makes the threshold useless for detecting real problems.
http.setResponseCallback(http.expectedStatuses(200, 409));

// ---------------------------------------------------------------------------
// Config — set TICKET_TYPE_ID to a real ticket_type id with limited quantity.
// Run this SQL first to create a controlled test ticket type:
//
//   INSERT INTO ticket_types (event_id, name, price, total_quantity, sold_quantity)
//   VALUES (<your_event_id>, 'k6 Test Ticket', 0.00, 5, 0)
//   RETURNING id;
//
// Then paste the returned id below.
// ---------------------------------------------------------------------------
const TICKET_TYPE_ID = __ENV.TICKET_TYPE_ID ? parseInt(__ENV.TICKET_TYPE_ID) : 1;
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Custom counters so we can see exactly how many succeeded vs were blocked
const reservedOK  = new Counter('tickets_reserved');
const reservedNo  = new Counter('tickets_rejected_oversell');
const serverError = new Counter('server_errors');

export const options = {
  scenarios: {
    spike: {
      executor: 'shared-iterations',
      vus: 50,        // 50 virtual users fire simultaneously
      iterations: 50, // one attempt each — all at once
      maxDuration: '30s',
    },
  },
  thresholds: {
    // After setResponseCallback, http_req_failed only counts truly unexpected responses (5xx, timeouts).
    // 409s are expected and excluded.
    http_req_failed: ['rate<0.01'],
    // 95% of all requests finish within 2s
    http_req_duration: ['p(95)<2000'],
    // Successful reservations must be <= ticket supply (5 in this setup)
    // k6 will fail the test if this threshold is exceeded
    tickets_reserved: ['count<=5'],
    // Every request must be either a successful hold (200) or a blocked one (409)
    server_errors: ['count==0'],
  },
};

export default function () {
  const url = `${BASE_URL}/create-checkout-session`;

  const payload = JSON.stringify({
    items: [{ ticket_id: TICKET_TYPE_ID, quantity: 1 }],
    // Unique email per VU so each gets a distinct guest user row
    email: `k6testuser${__VU}@example.com`,
  });

  const res = http.post(url, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  const got200 = res.status === 200;
  const got409 = res.status === 409;
  const got5xx = res.status >= 500;

  check(res, {
    'no server error (not 5xx)': () => !got5xx,
    'response is 200 or 409':    () => got200 || got409,
  });

  if (got200) reservedOK.add(1);
  if (got409) reservedNo.add(1);
  if (got5xx) serverError.add(1);
}
