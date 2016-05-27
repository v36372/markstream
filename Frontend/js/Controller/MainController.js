MarkStream.controller('MainController',['$scope','$timeout','$interval','decode',function($scope,$timeout,$interval,decode){
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
            promise = $interval(Process, 490);
        },500);
    };

    $scope.watermarks = [];

    ws.onmessage = function (event) {
        var frame = new Int16Array(event.data);
        var floatframe = {};
        floatframe.buffer = new Float32Array(frame.length);
        for(var i=0;i<frame.length;i++){
            floatframe.buffer[i] = frame[i]/32767;
        }
        // floatframe.embedd = embedd;
        // if(embedd){
        var promise = decode.QIMDecode(floatframe.buffer);
        promise.then(
            function(payload){
//                console.log(payload);
                if(payload != null && payload.length > 0)
                    floatframe.wm = payload;
            },
            function(errorPayload){
                console.log("error : " + errorPayload);
            });

        // }
        // else floatframe.wm = "";
        queue.push(floatframe);
    };

    var Process = function(){
        if(queue.length==0)
        {
            console.log("shit");
            return;
        }

        var audioChunk = queue[0].buffer;
        var audioBuffer = audio_context.createBuffer(1, 22050, 44100);
        audioBuffer.getChannelData(0).set(audioChunk);

        var source = audio_context.createBufferSource();
        source.buffer = audioBuffer;
        // console.log(queue[0]);
//        console.log(startTime);
//        var d2 = new Date().getTime();
//        console.log(d2 + " " + d1);
        source.start(startTime);
        source.connect(audio_context.destination);
        startTime += audioBuffer.duration;
        if(queue[0].wm != null){
//            console.log("haha");
            console.log(queue[0].wm);
//            var d1 = new Date().getTime();
//            console.log(d1);
            source.wm = queue[0].wm;
            source.onended = function(){
                
                console.log(this.wm);
//                var d2 = new Date().getTime();
//                console.log(d2);
                var wm = this.wm;

                if($scope.watermarks.length >0 && $scope.watermarks[$scope.watermarks.length-1].lastIndexOf('\n') == -1)
                    $scope.watermarks[$scope.watermarks.length-1] += wm;
                else
                    $scope.watermarks.push(wm.substr(1,wm.length-1));
            };
        }
        queue.shift();
    };
    var startTime;
    var closed = false;

    ws.onclose = function () {
        closed = true;
        $interval.cancel(promise);
        console.log("closed");
    };

}]);