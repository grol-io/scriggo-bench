calls = 90000;
b = 0;
for i = 0, calls, 1 do
    (function (x)
        b = b + x;
     end
    )(i);
end
