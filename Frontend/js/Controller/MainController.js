MarkStream.controller('MainController',['$scope','$websocket','$timeout',function($scope,$websocket,$timeout){
    var ws = $websocket.$new({
        url: 'ws://localhost:8081/stream',
        lazy: true
    });
    ws.onmessage = function(e){
        console.log("websocket : " + e.data);
    };
    ws.$on('$open', function () {
        console.log('The ngWebsocket has open!'); // It will print after 5 (or more) seconds
    })
    .$on('message', function (message) {
        console.log(message); // it prints 'dude, this is a custom message'
      });
    
    $timeout(function () {
        ws.$open(); // Open the connction only at this point. It will fire the '$open' event
    }, 1000);
}]);