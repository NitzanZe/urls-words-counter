# Go Urls Words Counter

### Given a list of urls containing essays and a list of valid words, the service will output he top 10 words from all the essays combined.
### The service will utilize the number of CPU cores and create workers that will work in parallel (With rate limitation)

### Prerequisites:
    * 'Make' installed
    * Go version 1.21 or later installed (Needed to be able to build the service)
    
### To run the service follow the next steps using the provided Makefile.

```
* Build the binary
   make build

* Run the service using provided path for the words list file and for the urls list file
  make run ARGS="<words_file_path urls_file_path"
  
 
 for example:
   make run ARGS="/home/user/words.txt /home/user/endg-urls-short"
   
* Optional: Remove the binary after usage
   make clean
   
* The service have the next environemnt variables with default values that can be changed by the user:

        LOG_FORMATTING                   default:"console" //Modes: json, console (console is for development needs for nicer look on the console)
	LOG_PATH                         default:"/var/logs"
	LOG_ENABLE_STD_OUTPUT            default:"true"
	LOG_ENABLE_FILE_OUTPUT           default:"true"
	LOG_LEVEL                        default:"info"
	
	MAXIMUM_WORKERS_REQUESTS_SECONDS default:"10"
	GET_TOP_N_WORDS                  default:"10"
	WORKERS_MULTIPLIER               default:"1"
	
* Exapmle for running with a changed environment variable:
    GET_TOP_N_WORDS=3 make run ARGS="/home/user/words.txt /home/user/endg-urls-short"

```
