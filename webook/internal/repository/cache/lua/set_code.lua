-- 发送到的 key,也就是 code:业务:手机号
local key=KEYS[1]
-- 使用次数,也就是验证次数
local cntKey = key..":cnt"
local val = ARGV[1]
--验证码的有效时间是十分钟,600秒
local ttl = tonumber(redis.call('ttl', key))

if ttl == -1 then
    -- key存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- key不存在或者有效期小于540，可以发验证码
    redis.call('set', key, val)
    redis.call('expire', key, 600)
    redis.call('set', cntKey, 3)
    redis.call('expire', cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end