SELECT
  a.account_id,
  a.balance AS stored_balance,
  COALESCE(SUM(l.delta), 0) AS ledger_balance
FROM accounts a
LEFT JOIN ledger_entries l
  ON a.account_id = l.account_id
GROUP BY a.account_id, a.balance
LIMIT 20;