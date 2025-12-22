SELECT
  account_id,
  SUM(delta) AS computed_balance
FROM ledger_entries
GROUP BY account_id
ORDER BY account_id
LIMIT 20;