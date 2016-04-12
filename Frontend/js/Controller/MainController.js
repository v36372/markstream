MarkStream.controller('MainController',['$scope','$timeout',function($scope,$timeout){
    var ws = new WebSocket("ws://localhost:8081/stream");
    
    ws.onopen = function(){  
        console.log("Socket has been opened!"); 
    };
        
    ws.onmessage = function (event) {
        console.log(event);
//        ws.send("x");
    };
    
    var closed = false;
    
    ws.onclose = function () {
        console.log("closed");
        closed = true;
    };

}]);