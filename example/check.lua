local filterKey = KEYS[1]
local shortKey = KEYS[2]
local val = ARGV[1]

local res = redis.call('BF.EXISTS', filterKey, val)
if res > 0 then
    return redis.call('get', shortKey)
end
return ""
