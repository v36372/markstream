'use strict'
var MarkStream = angular.module('MarkStream',[]);

MarkStream.factory('decode', function($q) {
    var Bit2Char = function(bits){
        var sum;
        var msg = ""
        var last = 0
        for (var i =0;i<bits.length;i++){
            sum <<= 1;
            sum += bits[i] - '0';
            if ((i-last+1)%8 == 0) {
                msg += (String.fromCharCode(sum));
                sum = 0;
                last = i + 1;
            }
        }
        return msg;
    };
    
    var QIMDecode = function(mag,phs){
        var step = [];
        step[0] = Math.PI/18;
        step[1] = Math.PI/14;
        step[2] = Math.PI/10;
        step[3] = Math.PI/6;
        step[4] = Math.PI/2;
        var stepsize = findStep(mag);
        var integer = Math.floor(phs / (step[stepsize] / 2));
        var r = phs/(step[stepsize]/2) - Math.floor(phs/(step[stepsize]/2));
        if (r < 0.5) {
            if (integer%2 == 0) {
                return 0;
            } else {
                return 1;
            }
        }
        if (r >= 0.5 ){
            if (integer%2 == 0) {
                return 1;
            } else {
                return 0;
            }
        }
        return 0;
    };

    var findStep = function(mag){
        var sMag = mag / (0.005);
        var group = Math.ceil(sMag / 0.2);
        if (group == 0) {
            group = 0;
        }
        if (group > 4 ){
            group = 4;
        }
        return group;
    };
    
	return {
		QIMDecode: function(audioChunk) {
            var deferred = $q.defer();
			var real = [];
            var img = [];
            var mag = [];
            var phs = [];

            for (var i = 0; i < audioChunk.length; i++) {
                real.push(audioChunk[i]);
                img.push(0);
            }

            transform(real,img);
            for (var i = 0; i < real.length; i++) {
                mag.push(Math.sqrt(real[i]*real[i]+img[i]*img[i]));
                phs.push(Math.atan2(img[i],real[i]));
            }
            var countone = 0;
            var countzero = 0;
            var str = "";
            for (var i = 0; i < 800; i++) {
                if (mag[i] < 0.0001){
                    continue;
                }
                var bit = QIMDecode(mag[i], phs[i]);
                if (bit == 1) {
                    countone++;
                } else {
                    countzero++;
                }
                if (countzero+countone == 5 ){
                    if (countzero > countone ){
                        str += "0";
                    } else {
                        str += "1";   
                    }
                    countzero = 0
                    countone = 0
                }
                if(str.length == 16){
                    var tmp = Bit2Char(str);
                    if(tmp[0]=='0' && tmp[1] != '0' && tmp[1] != '1')
                    {
                        deferred.reject('no wm');
                        return deferred.promise;
                    }
                }
            }
            str = Bit2Char(str);
            var tmp = str.lastIndexOf('\n');
            if(tmp == -1)
                deferred.resolve(str);
            else
                deferred.resolve(str.substr(0,tmp+1));
            return deferred.promise;
        }
	}
});

