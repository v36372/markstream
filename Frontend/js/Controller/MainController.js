MarkStream.controller('MainController',['$scope','$timeout',function($scope,$timeout){
    var ws = new WebSocket("ws://localhost:8081/stream");
    ws.binaryType = 'blob';
    var audio_context;
    var gain_node;
    var streaming_node;

    var init_web_audio = function() {

       if (typeof audio_context !== "undefined") {

        return;     //      audio_context already defined
        }

        try {

            window.AudioContext = window.AudioContext       ||
            window.webkitAudioContext ||
            window.mozAudioContext    ||
            window.oAudioContext      ||
            window.msAudioContext;

            audio_context = new AudioContext();  //  cool audio context established
            console.log("yay");

        } catch (e) {

            var error_msg = "Web Audio API is not supported by this browser\n" +
            " ... http://caniuse.com/#feat=audio-api";
            console.error(error_msg);
            alert(error_msg);
            throw new Error(error_msg);
        }

        gain_node = audio_context.createGain(); // Declare gain node
        gain_node.connect(audio_context.destination); // Connect gain node to speakers

    };
    
    init_web_audio();

    ws.onopen = function(){  
        console.log("Socket has been opened!"); 
    };

    var frame = [];
    
    ws.onmessage = function (event) {
//        console.log(event.data);
        decodeData = atob(event.data);
//        console.log(decodeData.length);
        var i =0;
        while (i<decodeData.length){
            buffer = new ArrayBuffer(8);
            int8bit = new Int8Array(buffer);
            float64bit = new Float64Array(buffer);
            substr = decodeData.substr(i,8);
//            var bytes = [];
        
            for (var j = 0; j < substr.length; ++j)
            {
//                bytes.push(substr.charCodeAt(i));
                int8bit[j] = substr.charCodeAt(j);
            }
//            console.log(int8bit);
            frame.push(float64bit[0]);
            i+=8;
        }
    };
    
    var closed = false;
    
    ws.onclose = function () {
        console.log("closed");
        closed = true;
        console.log(frame);
    };

}]);