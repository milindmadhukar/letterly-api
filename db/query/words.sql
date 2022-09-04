-- name: GetRandomWord :one
SELECT word
FROM wordlist
WHERE length(word) > $1
ORDER BY random()
LIMIT 1;

-- name: GetRandomCommonWord :one
SELECT word
FROM wordlist
WHERE is_common = true
ORDER BY random()
LIMIT 1;

-- name: IsPresent :one
SELECT EXISTS(
  SELECT 1
  FROM wordlist
  WHERE word = $1
);


-- name: IsCommonPresent :one
SELECT EXISTS(
  SELECT 1
  FROM wordlist
  WHERE word = $1 AND is_common = true
);
