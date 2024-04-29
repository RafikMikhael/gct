# gct
web server running on localhost that manages transcode jobs from src to an ABR ladder in destination

## Requirements  
- go v1.21 installed  
- zip file extracted into ~/go/src/github.com/RafikMikhael/gct

# gct
web server running on localhost that manages transcode jobs from src to an ABR ladder in destination
## Terminal 1: build, run, view server logging
>>make all  (or simply "go build")  
>>./gct  

Note: User can dictate the port number to use via "./gct -p 8085".  
If no port number is given, 8080 will be used as default.

## Terminal 2: Post some number of jobs spaced by a few seconds (each post retruns hash instantly)
```
curl -X POST "http://localhost:8080/api/v1/job/high?inputpath=%2Ftmp%2Fsrc&outputpath=%2Ftmp%2Fdst&w=1920&h=1080"  
```
{"id":c8e0ab68e1a226693de23eb81c2ffc49}  
```
curl -X POST "http://localhost:8080/api/v1/job/medium?inputpath=%2Ftmp%2Fsrc&outputpath=%2Ftmp%2Fdst&w=1000&h=500"  
```
{"id":7e4af802b6472cdfafb920e94006cf67}  
```
curl -X POST "http://localhost:8080/api/v1/job/low?inputpath=%2Ftmp%2Fsrc&outputpath=%2Ftmp%2Fdst&w=3000&h=15000"  
```
{"id":0c7ee460242d2b0076897d6e3cb4120e}  

## Terminal 3: Monitor the number of jobs still processing
```
curl "http://localhost:8080/api/v1/monitor"
```  
ongoing hashes=[49b252857f27f6faedef6fd7fab49524 9790281aa7676a3d57e60fc6368dd55b]  
ongoing hashes=[]



## Terminal 2: Use instantly returned hash to probe the job's current status (which renditions are done)
```
curl  "http://localhost:8080/api/v1/probe/7e4af802b6472cdfafb920e94006cf67"  
```
{"done":640x360}  
{"done":640x360,768x432}  
{"done":640x360,768x432,960x540}  
{"done":640x360,768x432,960x540,1280x720}  
204 (when done)  

## Terminal 2: inform the server to terminate after finishing current jobs and claiming all resources
```
curl "http://localhost:8080/api/v1/terminate"
```

## Terminal 1: inform the server to terminate instantly after claiming all resources
```
ctrl+c 
```

=====================================================================================================

# Directions:
- There are no restrictions on your choice of language, framework, nor data store so choose the ones you are most comfortable/expert in.
- Please commit your changes at least 24 hours before your pairing interview and confirm your submission with your point of contact or let them know if you have any questions.
- Make sure you get your solution to a working state. During the pairing interview we are going to ask you to build some new features on top of what you have done.
- We think it should normally take about a couple of hours to complete this assignment, but you are free to dedicate as much time as you see fit.
- We'll be running your software on our machines, so please ensure any prerequisites and necessary steps are documented.  Please don’t use anything private or anything that would require a login.

## Simulate a Live Video Transcoder

For this exercise, you’ll be building a small utility for managing a simulated transcoding process that can create multiple HLS ladders in parallel to simulate a live-streaming environment.  Note we’ll be mocking the actual transcode since we’re not testing you on ffmpeg, third party, or cloud knowledge.

- You could write the utility as an API, a code library, or a CLI/script, etc. - all would be acceptable and your choice is not part of the assessment.
- The script/library/api should be long-lived and remain running even while it has no jobs to do
- The script/library/api should be able to handle multiple simultaneous jobs (different input files)
- The script/library/api should include a method for starting a transcode process
    - inputs:
        - A location on disk of a video file
        - A value for quality/bandwidth of the outputs (`High` , `Medium` , `Low`)
        - a folder to place outputs for a CDN
- The program should simulate a live streaming scenario of generating a Live HLS ladder. Given the input file and quality setting, it should, in parallel, generate the transcode ladder for all lower resolutions using the table below to determine bitrate settings for each lower resolution
    
    | Resolution | Low | Medium | High |
    | --- | --- | --- | --- |
    | 640 x 360 | 120 | 145 | 160 |
    | 768 x 432 | 280 | 300 | 360 |
    | 960 x 540 | 1400 | 1600 | 1930 |
    | 1280 x 720 | 3080 | 3400 | 4080 |
    | 1920 x 1080 | 4500 | 5800 | 7000 |
- Since we’re just simulating the transcode, there is no need to actually integrate ffmpeg or do any transcoding or file handling.  Instead of implementing these, you can have the method log the desired resolution, bitrate, and other parameters passed in and sleep for some amount of time to simulate the transcode and prove it’s using the correct settings.  Please don’t actually implement FFMPEG or actually do the transcode. For example, both of these methods could be mocked (pseudocode):
    
    ```python
    class FfmpegWrapper
    {
      transcode(source_file, destination, width, height, bitrate): {Success|Error}
      probe(source_file): {result:string} (examples:'1920x1080', '1024x720', 'Error')
    }
    ```
