import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 50 },   // Ramp-up to 50 users
    { duration: '1m', target: 100 },   // Ramp-up to 100 users
    { duration: '2m', target: 100 },   // Stay at 100 users
    { duration: '30s', target: 200 },  // Spike to 200 users
    { duration: '30s', target: 0 },    // Ramp-down to 0 users
  ],
  thresholds: {
    'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1s
    'http_req_failed': ['rate<0.01'],  // Error rate < 1%
    'errors': ['rate<0.05'],           // Custom error rate < 5%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

// Test data
const testUsers = [
  { email: 'loadtest1@example.com', password: 'Password123!' },
  { email: 'loadtest2@example.com', password: 'Password123!' },
  { email: 'loadtest3@example.com', password: 'Password123!' },
];

export default function () {
  const user = testUsers[Math.floor(Math.random() * testUsers.length)];

  // Test 1: Register (10% of requests)
  if (Math.random() < 0.1) {
    const registerPayload = JSON.stringify({
      email: `test${Date.now()}@example.com`,
      password: 'Password123!',
      name: 'Load Test User',
    });

    const registerRes = http.post(`${BASE_URL}/api/v1/auth/register`, registerPayload, {
      headers: { 'Content-Type': 'application/json' },
    });

    check(registerRes, {
      'register status is 201': (r) => r.status === 201,
      'register response has token': (r) => JSON.parse(r.body).data.token !== undefined,
    }) || errorRate.add(1);
  }

  // Test 2: Login (40% of requests)
  if (Math.random() < 0.4) {
    const loginPayload = JSON.stringify({
      email: user.email,
      password: user.password,
    });

    const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, loginPayload, {
      headers: { 'Content-Type': 'application/json' },
    });

    const loginSuccess = check(loginRes, {
      'login status is 200': (r) => r.status === 200,
      'login response has token': (r) => {
        try {
          return JSON.parse(r.body).data.token !== undefined;
        } catch (e) {
          return false;
        }
      },
    });

    if (!loginSuccess) {
      errorRate.add(1);
    } else {
      // Extract token for authenticated requests
      const token = JSON.parse(loginRes.body).data.token;

      // Test 3: Get profile
      const profileRes = http.get(`${BASE_URL}/api/v1/users/me`, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      check(profileRes, {
        'profile status is 200': (r) => r.status === 200,
        'profile has user data': (r) => JSON.parse(r.body).data.user !== undefined,
      }) || errorRate.add(1);
    }
  }

  // Test 4: Health check (50% of requests)
  if (Math.random() < 0.5) {
    const healthRes = http.get(`${BASE_URL}/health`);

    check(healthRes, {
      'health status is 200': (r) => r.status === 200,
      'health response is valid': (r) => JSON.parse(r.body).status === 'healthy',
    }) || errorRate.add(1);
  }

  // Simulate think time
  sleep(Math.random() * 2 + 1); // 1-3 seconds
}

// Setup function (runs once at start)
export function setup() {
  console.log(`Starting load test against ${BASE_URL}`);

  // Verify application is running
  const healthRes = http.get(`${BASE_URL}/health`);
  if (healthRes.status !== 200) {
    throw new Error('Application is not healthy');
  }

  return { startTime: Date.now() };
}

// Teardown function (runs once at end)
export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000;
  console.log(`Load test completed in ${duration}s`);
}
