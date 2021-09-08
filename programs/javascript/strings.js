var size = 100;
var s = '';
for ( var r = 0; r < size*2; r++ ) {
    if ( r%2 == 0 ) {
        s += String.fromCharCode(r);
    }
}
var n = 0;
for ( var r = 0; r < size*2; r++ ) {
    if ( s.includes(String.fromCharCode(r)) ) {
        n++
    }
}
