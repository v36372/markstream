# Markstream

## What is this?
Markstream is an audio streaming application. It takes a music file as an input and start streaming audio data to the internet. To do streaming, I used websocket. 

Any client connect to the port of the websocket, will receive audio data and produce sound.

In runtime, the Markstream application can take in input as strings. Those strings will be broken down into bits. Those bits will be embedded into the frequencies of the audio. The frequencies combine together and produce an audio data, stream to clients normally.

The client receive the embedded audio data and start playing. There is no sound modification detected. Finally, the embedded information is retrieved and printed out.

## What does Markstream used?
The project used Golang as the server programming language. On the front-end, to iterate faster, I use simple AngularJS.

## What next?
There are many upgrades that I want to apply. I will also deploy my project to a VPS.

