SELECT
  a.attname
FROM
  pg_index i
  JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
WHERE
  i.indrelid = ('"' || $1::text || '"."' || $2::text || '"')::regclass
  AND i.indisprimary
ORDER BY
  array_position(i.indkey, a.attnum)
