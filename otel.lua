local myTable = {
    'loop_fib?n=',
}

request = function()
    path = '/' .. myTable[ math.random( #myTable ) ] .. math.random(10,20)
    return wrk.format("GET", path)
end
