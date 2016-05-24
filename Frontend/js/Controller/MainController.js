MarkStream.controller('MainController',['$scope','$timeout','$interval',function($scope,$timeout,$interval){
    var ws = new WebSocket("ws://localhost:8081/stream");
    ws.binaryType = 'arraybuffer';
    var audio_context =  new AudioContext();

    ws.onopen = function(){  
        console.log("Socket has been opened!"); 
    };

    var queue = [];
    
    var embedd = false;

    $scope.play = function(){
        $timeout(function(){
            startTime = audio_context.currentTime+0.5;
            promise = $interval(Process, 400);
        },500);
    };

    $scope.watermarks = [];

    ws.onmessage = function (event) {
        if(event.data == "start"){
            embedd = true;
            return;
        }
        if(event.data == "end"){
            embedd = false;
            return;
        }
        var frame = new Int16Array(event.data);
        var floatframe = {};
        floatframe.buffer = new Float32Array(frame.length);
        for(var i=0;i<frame.length;i++){
            floatframe.buffer[i] = frame[i]/32767;
        }
        floatframe.embedd = embedd;
        if(embedd){
            floatframe.wm = Decode(floatframe.buffer);
        }
        else floatframe.wm = "";
        queue.push(floatframe);
    };

    var Process = function(){
        if(queue.length==0)
            return;
        var audioChunk = queue[0].buffer;
        var audioBuffer = audio_context.createBuffer(1, 22050, 44100);
        audioBuffer.getChannelData(0).set(audioChunk);

        var source = audio_context.createBufferSource();
        source.buffer = audioBuffer;
        // console.log(queue[0]);
        // console.log(startTime);
        source.start(startTime);
        source.connect(audio_context.destination);
        startTime += audioBuffer.duration;
        if(queue[0].embedd == true){
            source.wm = queue[0].wm;
            source.onended = function(){
                console.log(this.wm);
                var wm = this.wm;

                if($scope.watermarks.length >0 && $scope.watermarks[$scope.watermarks.length-1].lastIndexOf('\n') == -1)
                    $scope.watermarks[$scope.watermarks.length-1] += wm;
                else
                    $scope.watermarks.push(wm);
            };
        }
        queue.shift();
    };
    var startTime;
    var promise;
    var closed = false;
    
    var Decode = function(audioChunk){
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
                countone++
            } else {
                countzero++
            }
			if (countzero+countone == 5 ){
				if (countzero > countone ){
					str += "0"
				} else {
					str += "1"
				}
				countzero = 0
				countone = 0
			}
        }
        str = Bit2Char(str);
        console.log(str);
        var tmp = str.lastIndexOf('\n');
        if(tmp == -1)
            return str;
        else
            return str.substr(0,tmp);
    };

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

    ws.onclose = function () {
        closed = true;
        $interval.cancel(promise);
        console.log("closed");
    };

}]);