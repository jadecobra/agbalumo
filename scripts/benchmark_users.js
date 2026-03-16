import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '1m', target: 10000 },  // Ramp up to 10k users
        { duration: '2m', target: 10000 },  // Sustain 10k users
        { duration: '3m', target: 100000 }, // Spike to 100k users
        { duration: '1m', target: 0 },      // Ramp down
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'], // http errors should be less than 1%
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    },
};

const BASE_URL = __ENV.TARGET_URL || 'https://localhost:8443';

export default function () {
    const rand = Math.random();
    let res;

    if (rand < 0.80) {
        // 80% Browse Home Page
        res = http.get(`${BASE_URL}/`);
        check(res, { 'status is 200': (r) => r.status === 200 });
    } else if (rand < 0.90) {
        // 10% Paginate through listings
        const page = Math.floor(Math.random() * 50) + 1; // Random page 1-50
        res = http.get(`${BASE_URL}/?page=${page}`);
        check(res, { 'status is 200': (r) => r.status === 200 });
    } else if (rand < 0.95) {
        // 5% Search for specific categories
        res = http.get(`${BASE_URL}/?category=Business`);
        check(res, { 'status is 200': (r) => r.status === 200 });
    } else {
        // 5% View single listing (using search as a proxy for complex queries)
        res = http.get(`${BASE_URL}/?q=tech`);
        check(res, { 'status is 200': (r) => r.status === 200 });
    }

    sleep(1); // Simulate user reading time
}
