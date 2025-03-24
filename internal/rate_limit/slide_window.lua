local key = KEYS[1]
-- window interval(window size)
local interval = tonumber(ARGV[1])
-- rate limit
local rate = tonumber(ARGV[2])

-- current time(window end time)
local now = tonumber(ARGV[3])

-- unique id
-- used to solve the problem of counting multiple requests at the same time
local uid = ARGV[4]

-- window start time
local window_start = now - interval


redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)
local count = redis.call('ZCOUNT', key, window_start, '+inf')

if count >= rate then
    -- reject the request
    return 0
else 
    -- allow the request
    -- add the request to the window
    redis.call('ZADD', key, now, uid)
    redis.call('PEXPIRE', key, interval)
    return 1
end
