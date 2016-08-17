# Markstream - VNU HCMUS FIT APCS Graduation Thesis 

## Author

Trong-Tin Nguyen | nttin@apcs.vn | +84918644143

## What does this project is about?
Markstream is an audio streaming application. It takes a music file as input and start streaming audio data to the internet.

Any client connect to the port of the websocket, would receive audio data and produce sound.

Being combined with audio watermarking, in runtime, the Markstream application can take in input as strings. Those strings will be embedded into the audio data and stream to clients as normal.

The client receive the watermarked audio data and start watermark extraction. Once the watermarked sound is played, the watermark string is displayed too.

## How is this CD organized?
This CD is organized into folders with its specific reasons.

### `/MARKSTREAM`
This folder contains the project itself.

While `/SETUP` delivers what is needed to get this project up and running, `/SOURCE` folder contains the whole source code.

### `/THESIS`
This folder contains the document of this project. 

While `/LATEX` keep the `tex` source code inside, needed to build the `LaTeX` project, the `PDF` folder only contains the final official documents.

### `/SOFT`
This folder contains needed softwares for this project to run.

### `/REF`
This folder contains reference documents that this project used.

## How to get this project up and running?


1. Open up command line at the root folder.

2. Navigate to `MARKSTREAM/SETUP/SERVER` with: `cd MARKSTREAM/SETUP/SERVER`

3. Start up the server with an audio file as input: `./MarkStream RWC_002.wav`

4. Open another command line window.

5. Navigate to `MARKSTREAM/SETUP/CLIENT` with `cd MARKSTREAM/SETUP/CLIENT`

6. Start up the client: `./http-server`

7. Open the internet browser and go to the client application at `0.0.0.0:8080`
  

## How to test this project?


### Test watermarking correctness


1. Open up command line at the root folder.

2. Navigate to `MARKSTREAM/SOURCE/test`: `cd MARKSTREAM/SOURCE/test`

3. Execute command: `./execute 2` or `./execute 3` or `./execute 4` (Each of them is a different settings of watermarking)

4. After the test is done running. We can collect test results and view on a chart.

5. Execute the following commands: `./prepare_chartdata` and `./makechart`
 
6. Launch `Gochart` to sketch chart: `./gochart`

7. Open Internet browser and navigate to `0.0.0.0:8000`. The chart should be displayed there.


### Test audio quality

1. Open up command line at the root folder
 
2. Navigate to `MARKSTREAM/SOURCE/test`: `cd MARKSTREAM/SOURCE/test`

3. If you are on Linux, you would have to install wine: `sudo apt-install wine`. Otherwise, continue.
 
4. Execute the following commands to launch ResamAudio and PQevalAudio to test audio quality:
```
wine ./ResampAudio.exe -s 48000 RWC_60s/RWC_001.wav RWC_60s/RWC_001_sample.wav

wine ./ResampAudio.exe -s 48000 RWC_60s/RWC_001_2_wm.wav RWC_60s/RWC_001_2_wm_sample.wav

wine ./PQevalAudio.exe RWC_60s/RWC_001_sample.wav RWC_60s/RWC_001_2_wm_sample.wav

```

The final result is the score of our watermarked file compare to the original file.
