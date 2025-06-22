local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local timestamp = tonumber(ARGV[2])
local id = ARGV[3]
local ttlSeconds = tonumber(ARGV[4])
local count = redis.call("zcard", key)
local allowed = count < capacity

if allowed then
  redis.call("zadd", key, timestamp, id)
  redis.call("expire", key, ttlSeconds)
end

return { allowed, count }