local myTable = {
    'fib?n=',
}

function delay()
    return 4
end

request = function()
    path = '/' .. myTable[ math.random( #myTable ) ] .. math.random(10,10)
    -- path = '/test'
    return wrk.format("GET", path)
end
