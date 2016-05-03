MarkStream.controller('MainController',['$scope','$timeout',function($scope,$timeout){
    var ws = new WebSocket("ws://localhost:8081/stream");
    ws.binaryType = 'arraybuffer';
    var audio_context =  new AudioContext();

    ws.onopen = function(){  
        console.log("Socket has been opened!"); 
    };

    var queue = [];
    
    
    ws.onmessage = function (event) {
        var frame = {};
        frame.buffer = new Float32Array(event.data);
//        audio_context.decodeAudioData(event.data, function(buffer) {
        queue.push(frame);
//        }, null);
        
    };
    
    var closed = false;
    
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
          source.start(startTime);
          source.connect(audio_context.destination);
//          console.log(audioBuffer.duration);
          startTime += audioBuffer.duration;
          console.log(startTime);
        }

    };

}]);