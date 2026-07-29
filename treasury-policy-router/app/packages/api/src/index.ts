import { Hono } from 'hono';
import { PolicyHookClient } from '@tpr/sdk';

const app = new Hono();

app.get('/health', (c) => c.json({ status: 'ok' }));

app.get('/reports/policy-summary', async (c) => {
  const start = c.req.query('start');
  const end = c.req.query('end');
  return c.json({
    total: 0,
    approved: 0,
    rejected: 0,
    flagged: 0,
    dateRange: { start, end },
  });
});

export default app;
