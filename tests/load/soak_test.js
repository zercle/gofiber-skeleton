import http from 'k6/http';
import { check, sleep } from 'k6';

// Soak test - sustained load over time
export const options = {
  stages: [
    { duration: '2m', target: 400 },   // Ramp up to 400 users
    { duration: '3h56m', target: 400 }, // Stay at 400 for ~4 hours
    { duration: '2m', target: 0 },     // Ramp down to 0
  ],
  thresholds: {
    'http_req_duration': ['p(99)<1500'],
    'http_req_failed': ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export default function () {
  const res = http.get(`${BASE_URL}/health`);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
