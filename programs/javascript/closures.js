const calls = 90000;
var b = 0;
for (var i = 0; i < calls; i++) {
    (function (x) {
        b += x;
    })(i);
}