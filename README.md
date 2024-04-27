# gct

## Terminal 1: build, run, server logs
>>make all
>>./gct

## Terminal 2: Post some number of jobs spaced by a few seconds (each post retruns hash instantly)
>>curl -X POST "http://localhost:8080/api/v1/job/high?inputpath=%2Ftmp%2Fsrc&outputpath=%2Ftmp%2Fdst&w=1920&h=1080"
{"id":c8e0ab68e1a226693de23eb81c2ffc49}
>>curl -X POST "http://localhost:8080/api/v1/job/medium?inputpath=%2Ftmp%2Fsrc&outputpath=%2Ftmp%2Fdst&w=1000&h=500"
{"id":7e4af802b6472cdfafb920e94006cf67}
>>curl -X POST "http://localhost:8080/api/v1/job/low?inputpath=%2Ftmp%2Fsrc&outputpath=%2Ftmp%2Fdst&w=3000&h=15000"
{"id":0c7ee460242d2b0076897d6e3cb4120e}

## Terminal 3: Monitor the number of jobs still processing
>>curl "http://localhost:8081/"
ongoing hashes=[49b252857f27f6faedef6fd7fab49524 9790281aa7676a3d57e60fc6368dd55b]
ongoing hashes=[]


## Terminal 2: Use instantly returned hash to probe the job's current status (which renditions are done)
>>curl  "http://localhost:8080/api/v1/probe/7e4af802b6472cdfafb920e94006cf67"
{"done":640x360}
{"done":640x360,768x432}
{"done":640x360,768x432,960x540}
{"done":640x360,768x432,960x540,1280x720}
{"done":7e4af802b6472cdfafb920e94006cf67}

## Terminal 2: inform the server to terminate after claiming all resources
>>curl -X POST "http://localhost:8080/api/v1/terminate"