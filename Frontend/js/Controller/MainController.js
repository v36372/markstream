MarkStream.controller('MainController',['$scope','$timeout',function($scope,$timeout){
    var ws = new WebSocket("ws://localhost:8081/stream");
    ws.binaryType = 'arraybuffer';
    var audio_context =  new AudioContext();

    ws.onopen = function(){  
        console.log("Socket has been opened!"); 
    };

    var queue = [];
    
    var embedd = false;
    ws.onmessage = function (event) {
        // console.log(event);
        if(event.data == "start"){
            // console.log("hehe");
            embedd = true;
            return;
        }
        if(event.data == "end"){
            // console.log("hieie");
            embedd = false;
            return;
        }
        var frame = new Int16Array(event.data);
//        console.log(frame);
        var floatframe = {};
        floatframe.buffer = new Float32Array(frame.length);
        for(var i=0;i<frame.length;i++){
            floatframe.buffer[i] = frame[i]/32767;
        }
        floatframe.embedd = embedd;
        queue.push(floatframe);
    };
    
    
    
    var closed = false;
    
    var Decode = function(audioChunk){
      
    };

    var QIMDecode = function(){

    };

    var findStep = function(){

    };

    ws.onclose = function () {
        console.log("closed");
        var startTime = audio_context.currentTime;

        for (var i = 0; i<queue.length; ++i) {
          // Create/set audio buffer for each chunk
          var audioChunk = queue[i].buffer;
          var audioBuffer = audio_context.createBuffer(1, 22050, 44100);
          audioBuffer.getChannelData(0).set(audioChunk);

          var source = audio_context.createBufferSource();
          source.buffer = audioBuffer;
          if(queue[i].embedd == true){
            var watermark = Decode(audioChunk);
            console.log(watermark);
          }
          source.start(startTime);
          source.connect(audio_context.destination);
          startTime += audioBuffer.duration;
        }

    };

}]);