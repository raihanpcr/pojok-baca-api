import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
      vus: 100,
      duration: '15s',
};

export default function () {
      const res = http.get('https://pojok-baca-api-fb30b0912dab.herokuapp.com/api/products');
      check(res, {
            'status was 200': (r) => r.status === 200,
      });
      sleep(1);
}