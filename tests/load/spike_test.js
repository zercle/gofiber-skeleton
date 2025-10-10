import http from 'k6/http';
import { check, sleep } from 'k6';

// Spike test configuration
export const options = {
  stages: [
    { duration: '10s', target: 100 },  // Below normal load
    { duration: '1m', target: 100 },
    { duration: '10s', target: 1400 }, // Spike to 1400 users
    { duration: '3m', target: 1400 },  // Stay at 1400 for 3 minutes
    { duration: '10s', target: 100 },  // Scale down to normal
    { duration: '3m', target: 100 },
    { duration: '10s', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    'http_req_duration': ['p(99)<2000'], // 99% under 2s even during spike
    'http_req_failed': ['rate<0.05'],    // Error rate < 5%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export default function () {
  const res = http.get(`${BASE_URL}/health`);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.5);
}
