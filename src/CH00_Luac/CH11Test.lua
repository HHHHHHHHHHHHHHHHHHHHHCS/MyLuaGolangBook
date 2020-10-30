print(getmetatable("foo"))
print(getmetatable("bar"))
print(getmetatable(nil))
print(getmetatable(false))
print(getmetatable(100))
print(getmetatable({  }))
print(getmetatable(print))

print("====================================")

t = {}
mt = {}
setmetatable(t, mt)
print(getmetatable(t) == mt)
--对于非表类型的 可以用debug.setmetatable
debug.setmetatable(100, mt)
print(getmetatable(200) == mt)

print("====================================")

mt = {}
mt.__add = function(v1, v2)
    return vector(v1.x + v2.x, v1.y + v2.y)
end

function vector(x, y)
    local v = { x = x, y = y }
    setmetatable(v, mt)
    return v
end

v1 = vector(1, 2)
v2 = vector(3, 5)
v3 = v1 + v2
print(v3.x, v3.y)