
const size = 400;
var s = [];

for (var i = 0; i < size; i++ ) {
    s[i] = i
}

for (var i = 0; i < s.length; i++) {
    for (var j = 0; j < s.length; j++) {
        s[j] += s[i];
    }
}
