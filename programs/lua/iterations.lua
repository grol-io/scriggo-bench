size = 400;
s = {};

for i = 0, size, 1 do
    s[i] = i
end

for i, v in ipairs(s) do
    for j, t in ipairs(s) do
        s[j] = t + v
    end
end
