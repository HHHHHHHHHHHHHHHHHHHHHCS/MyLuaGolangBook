local x = { 3 }

function x:ForPrint(len)
    for i = x[1], len do
        print(i .. " : Hello , World!")
    end
end

x:ForPrint(6)

----------------------------------------


function div0(a, b)
    if b == 0 then
        error("DIV BY ZERO !")
    else
        return a / b
    end
end

function div1(a, b)
    return div0(a, b)
end

--测试传播
function div2(a, b)
    return div1(a, b)
end

ok, result = pcall(div2, 4, 2)
print(ok, result)

ok, err = pcall(div2, 5, 0)
print(ok, err)

ok, err = pcall(div2, {}, {})
print(ok, err)



----------------------------------------



local mt = {}

function vector(x, y)
    local v = { x = x, y = y }
    setmetatable(v, mt)
    return v
end

mt.__add = function(v1, v2)
    return vector(v1.x + v2.x, v1.y + v2.y)
end

mt.__sub = function(v1, v2)
    return vector(v1.x - v2.x, v1.y - v2.y)
end

mt.__mul = function(v1, n)
    return vector(v1.x * n, v1.y * n)
end

mt.__div = function(v1, n)
    return vector(v1.x / n, v1.y / n)
end

mt.__len = function(v)
    return (v.x * v.x + v.y * v.y) ^ 0.5
end

mt.__eq = function(v1, v2)
    return v1.x == v2.x and v1.y == v2.y
end

mt.__index = function(v, k)
    if k == "print" then
        return function()
            print("[" .. v.x .. " , " .. v.y .. "]")
        end
    end
end

mt.__call = function(v)
    print("[" .. v.x .. " , " .. v.y .. "]")
end


v1 = vector(1, 2)
v1:print()
v2 = vector(3, 4)
v2:print()
v3 = v1 * 2
v3.print()
v4 = v1 + v3
v4:print()
print(#v2)
print(v1 == v2)
print(v2 == vector(3, 4))
v4()


--------------------------

local t = { "a", "b", "c" }
t[2] = "B"
t["foo"] = "Bar"
local s = t[3] .. t[2] .. t[1] .. t["foo"] .. #t