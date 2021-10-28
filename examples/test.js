import sqs from 'k6/x/sqs';
import { check } from 'k6';
import { Trend, Rate } from 'k6/metrics';

const QUEUE_NAME = 'test';

let sqsTimeDuration = new Trend('sqs_time_duration');
let sqsFailed = new Rate('sqs_failed')

export function setup() {
    let client = sqs.new({
        url: 'http://localhost:8081/api',
        user_id: '1234567890',
    })
    client.createQueue(QUEUE_NAME)
}

export default function () {
    let client = sqs.new({})

    var start = new Date();
    let response = client.sendMessage({
        queue_name: QUEUE_NAME,
        message_body: "This is a test message"
    })
    sqsTimeDuration.add(new Date() - start)
    sqsFailed.add(response !== null)

    check(response, { 'Response has been received': (r) => r !== null });
}

export function teardown() {
    let client = sqs.new({})
    client.deleteQueue(QUEUE_NAME)
}

export let options = {
    thresholds: {
        sqs_failed: ['rate=1'],
        sqs_time_duration: ['p(95)<15'],
    },
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 800,
            timeUnit: '1s',
            duration: '30s',
            preAllocatedVUs: 5,
            maxVUs: 20,
        }
    }
}
