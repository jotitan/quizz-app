# Quizz app

Application to make quizz between peoples like a kahoot.

## Technical details

### Stack 

* Backend in Go 1.7
* Frontend in React JS
* Use SSE (Server Side Event) to maange communication between players

### How to start

Frontend : 
* cd resources
* npm start (by default on port 3000)

Backend : 
* cd src
  * create conf.yaml in starting folder with configuration
    ```yaml
      port: 9001 (port where to run server)
      storage: filer (kind of storage, only filer available)
      path: my_path (path where store data)
      ffmpeg: my_ffmpeg_path (path where ffmpeg are, used to cut mp3 files)
      resources: my_resources_path (to serve frontend, folder where resources are)
    ```
* go run main/run_quizz.go (by default start on port 9001)


### How to build

Frontend : npm install  
Backend : go build main/run_quizz.go
