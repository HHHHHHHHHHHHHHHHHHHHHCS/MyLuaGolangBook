local x = { 3 }

function x:ForPrint(len)
    for i = x[1], len do
        print(i .. " : Hello , World!")
    end
end

x:ForPrint(6)