-- 发送到的 key,也就是 code:业务:手机号
local key=KEYS[1]
-- 使用次数,也就是验证次数
local cntKey = key..":cnt"
--用户输入的验证码
local expectCode = ARGV[1]

local cnt = tonumber(redis.call('get', cntKey))
local code = redis.call('get', key)

if cnt == nil or cnt <= 0 then
    -- 验证次数已经用完
    return -1
end

if code == expectCode then
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    return -2
end