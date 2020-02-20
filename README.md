# HTTP log monitor

This application monitors incomming http traffic by reading a log file.
The log file must be formatted with [Common Log Format](https://www.w3.org/Daemon/User/Config/Logging.html).

Its default behaviour is the following :
- reads the log file from `/tmp/access.log`, creates it with unix right `0644` if not found 
- outputs error messages to `/usr/local/var/log/http_log_monitor.log`
- uses a data-refresh-rate of 10 seconds
- the alert measuring time is of 10 minutes with a threshold of 10 req/s

## Getting started
### Run the app locally
To read the documentation:
```bash
# Log monitor
go run cmd/logmonitor/main.go --help

# Log writter
go run cmd/writelog/main.go --help
```

To run the test with the default parameters specified in the email please type the following commands:
```bash
go run cmd/logmonitor/main.go

# Run from another shell to write logs to the file
go run cmd/writelog/main.go --lines=3000 --duration=180 --path=/tmp/access.log
```
*Quit the app using* `ESC` or `CTRL C`

To run the test with custom parameters:
```bash
go run cmd/logmonitor/main.go --path=/tmp/foo.log --update=1s --alert-period=1s --alert-threshold=5

# Run from another shell to write logs to the file
go run cmd/writelog/main.go --lines=20 --duration=1 --path=/tmp/foo.log
```
*Quit the app using* `ESC` or `CTRL C`

### Run with docker
```bash
docker build -t logmonitor .
docker run --name logmonitor logmonitor

# Pass specific flags to customise logmonitor's behaviour (see above section)
docker exec -it logmonitor bash -c /app/logmonitor

# Exit with CTRL C from the interactive pane then
docker stop logmonitor
```

### Run the tests
```bash
go test ./...
```

To have more information on the alert-specific tests please read the file `pkg/task/alert_test.go`.

## Improvements
- More tests need to be implemented. At the moment, some readers and a few tasks have their tests implemented
- Add a proper sleep mechanism in the backend loop to reduce CPU use
- Improve the `task` and `reader` interfaces to avoid losing input-types at compile time (replace ...interface{} by well identified parameters) 
- Format the metrics computed by the app to export them to a `prometheus` instance. This app would be converted into an exporter.
- The current UI would be disabled and data visualisation would be done by `grafana`
- Be able to choose the computation method for the rates, whether based on the input rate or on the request-time logged in the file
- Possibly make alerts fetch their own `measure rates` to decouple them from the `measure rates` task. 
It is more CPU, RAM and time consumming but that way all tasks could be executed in separate goroutines.
- Formalise and implement a task-dependency-system for the `TaskEnv` struct. It would enable to define task sequences and link task-input and output together
- Read the tasks executed by the backend and their relation from yaml file. That way the app would highly configurable. That way new tasks could be easily implemented
and integrated. It would be even more useful if the app becomes an exporter.
- Implement a content analysis of the log file before capturing live input-data. It would give context to the reader to interpret the metrics exposed by the tool
- Be able to read several log files and aggregate their content
- Be able to set up several alerts at the same time
- Make alerts capable of being switched on on any task-output value
- Alert notification methods should not be limited to displaying a text but should rather send an email, a mobile text, a slack notification...
- Implement other reader types so that it would be possible to read from other sources than files. Interesting options would be to read from sockets, a gRPC API, a REST API...
- Implement readers capable of connecting to relational databases (Postgres) or NoSQL ones (Elastic Search)
- Be able to parse several HTTP log formats

## Design presentation
First of all let's talk about the language choice - I chose `go` for the performance it offers, ability to easily express
parallel/concurent code, quality of the standard and third-party libraries to quickly implement cloud-related-services. Also,
in the long run go has native constructs (interfaces and easy-to-use-composition) that makes software extensible rather easily over time.
Finally, the quality of the tooling in go is generally a good help to be on track fast.

The implementation makes a wide use of OOP concepts (even though go is not object oriented, but it has composition) and a tad of functional programming. 
Object Oriented Programming was used as it allows for maintainability (low coupling), extensibility and testability. 

The app was implemented using a traditional MVC pattern. Below is what makes up each component :
- The `model` is the log.Info structure representing a log line represented in Common Log Format. A []log.Info hence representing a whole log file.
Also all data output from `tasks` are part of the model.
- The controllers are the tasks present in `pkg/task`. They take a model (log file, previous command output) and manipulate them to output other data.
- Finally the view is made up of the `app.ViewFrame` and all structs related to the user interface (ui files present in `pkg/app`).

More in details, the app was made to be highly parallel. That is, it uses `tasks` (we generally call them workers in our cloud-distributed environments)
that can be executed asynchronously and generally not depending on any other. That way, by design, it is scalable as tasks could be executed on dedicated
worker machines. Using tasks to perform computation reduces coupling which is a great benefit for testing the system and reduces breakdowns and cost
maintaining highly coupled objects. Last but not least, due to its architecture most of the application achieves concurrence management in a lock-free manner
(see the `dbufreader` in `pkg/reader` which implements a technique used for screen rendering to achieve high and safe throughput - the double buffering).

The whole system is periodicaly synchronous meaning that after a certain always-the-same amount of time, the app reaches a coherent state. 
This state is reached by default after `10 seconds`. This is the update duration. This is called a `frame` in the video game industry
and would be named that way in the rest of the document.

### Project layout
Here is the project layout (without the files) :
```
.
├── cmd                 // binary directory
│   ├── logmonitor      // app's binary
│   └── writelog        // helper to write logs to input log file
├── containers          // container folder
│   └── test_nginx      // test nginx writing logs to /tmp/access.log
└── pkg                 // go code directory
    ├── app             // app's code, frontend and backend
    ├── log             // log-representation-related source
    ├── logger          // the app's logger
    ├── reader          // objects able to read content and return logs
    ├── task            // objects generating metrics
    └── timer           // custom timer, essentially used to ease testing
```

Test files can be found next to the files containing their system-under-tests.

### Package descriptions
#### App
This package contains the backend and frontend code.

The frontend is made up of a renderer (`render.go`) and has its name suggest renders the UI. The Updates are carried out by `Renderer.update` polling
a `chan ViewFrame` returning a renderable-frame's content. The remaining part of the frontend is the `layout.go` file that sets the app's widget layout.
The renderer has three functions to control its execution – `init`, `run` and `shutdown`.

Now let's talk about the backend. Its structure is simple, all its content sits in `backend.go`. It implements a simple set of functions (they could have
been made into an interface but I did not need it at the time) to control its execution – `init`, `run` and `shutdown`. The `run` function is executed in 
the background in a separate goroutine to allow for seemless updates.

Before passing to the next package description, let's talk about the `add` method. The backend stores a list of `task.Task` interface to automate
task execution and bring flexibility into the app's configuration. Please find below add's signature and a more detailed description of its action :
```go
func (b *Backend) add(tasks ...Taskenv)

type Taskenv struct {
    Task        task.Task
    InitParams  []interface{}
}
```
This architecture allows to define in the main (or from a config file... not implemented though) a list of tasks to be loaded by the backend and then executed.
The `InitParams` slice from `TaskEnv` is used to feed parameters in order to the `Task.Init` function. Other attributes should be defined in `TaskEnv` but
time is lacking and they are more complex to support. Also, task dependency would require to be implemented to fully take advantage of this system.

Before going on to the next description here is a simple call to define backend tasks for execution :
```go
b.add(
		Taskenv{
			Task:       &fetchLogs,
			InitParams: []interface{}{conf.LogFilePath, reader.CommonLogFormatParser(), conf.UpdateFrameDuration},
		},
		Taskenv{
			Task: &mostHits,
		},
		Taskenv{
			Task: &rates,
		},
		Taskenv{
			Task: &countCodes,
		},
		Taskenv{
			Task:       &alert,
			InitParams: []interface{}{conf.AlertFrameDuration, conf.AlertThreshold},
		},
	)
```

Then backend's `init` can do the following :
```go
for _, t := b.tasks {
    if err := t.Task.Init(t.InitParams...); err != nil {
        return err
    }
}
```
The same can be done with all Task methods except for Run due to some current (time) limitations.

#### Task
The task package contains all tasks executed by the backend. A task's role is to work out metrics then usable by the rest of the app. Tasks are made to be executed
independently. Although there is one example that reuses the result of a previous task. This action could be avoided however by recomputing necessary data.

All tasks implement a common interface whose description can be found below :
```go
type Task interface {
	// Init sets up the task. Think of it as a constructor.
	Init(...interface{}) error

	// BeforeRun sets some data before the task is run. Remember that
	// a task can be executed several times, accross several time frames.
	// This function and AfterRun is there to setup/cleanup the task's state
	// in case a long iterative work is planned
	BeforeRun(...interface{}) error

	// Run executes the task
	Run(...interface{}) error

	// AfterRun like BeforeRun except it is called after Run has been executed
	AfterRun() error

	// IsDone is true if the task's work has been completed
	IsDone() bool

	// Close shutdowns the task. Call Init to use it again
	Close() error
}
```

Now let us describe the tasks.

##### Fetch logs
This task reads the file asynchronously using an asynchronous double-buffering reader (find explanations in the next section). It is the one responsible of feeding
data to all other tasks. Text data is converted into `log.Info` data structured for other tasks to use.

To be able to feed data in a lock-free maner, the task writes to a dedicated write buffer while it reads from a dedicated read buffer. When
the frame is over, all tasks stop their action and the `fetch log task` swaps its read and write buffers. That way tasks can be fed the new lines on the next frame.
It is worth pointing out that, in addition to the double buffering technique, this task is always one frame ahead of the others. That is how concurrency is avoided `without using a mutex`.

##### Most hit sections
This task reads the log provided by `fetch logs`, indexes them by sections (using a map) then counts the section occurence and the HTTP request (`GET`, `POST`...) occurences. The result of this task is the following :
```go
// The task returns a slice of Hits
type Hit struct {
	Section string            // URL section in /compute/create/..., compute is the section
	Total   uint64            // number of section occurences
	Methods map[string]uint64 // a map of method count (i.e. POSTcount := Methods["POST"])
}
```

#### Count codes
This task counts the number of http-codes from `fetch logs`' input. It simply gathers all the logs (using a map) by http-return code and increment a
counter on each occurence. The map of counters is the returned from the task.

#### Alert
The alert task continuously checks the `average-request-rate` (req/s) is always below the alert `threshold` (this value is given as an argument to the CLI). The
average-request-rate is computed on a custom `period` of time (it also is an input from the CLI).

If during the measuring-stage the average-rate is above the threshold, then the alert is switched on.

Implementation-wise I reused the results from `measure rates` task. It introduces a dependency between the two tasks and could be avoided by making alerts recompute
all the values it needs (the second option being more CPU, RAM and time consumming but can be totally parallelised).

Here is an alert result:
```go
type AlertState struct {
	// IsOn is true when the alert is active.
	// It is false if there wasn't any alert or the system recovered
	IsOn bool

	// Duration is the time spent to check the alert, and the time it takes to cool down
	Duration time.Duration

	// Threshold is the value triggering the alert if the average rate is greater on a whole Duration-period
	Threshold uint64

	// Avg is the average req/s the alert was triggered at (always 0 if IsOn)
	Avg uint64

	// NbReqs is number of requests that triggered the alert. Always equals 0 when IsOn == false
	NbReqs uint64

	// Date is the time the alert has been switched on or off. It has a default value
	// in case the alert has never been activated.
	Date time.Time
}
```

#### Measure rates
This task measures all the traffic-related data. It measures them considering a per-frame basis and a global-basis (whole app execution). Before focusing on the
implementation description, let's describe what data is output:
```go
// Rates contains all type of rates and measures taken by the task
type Rates struct {
	Global GlobalRates
	Frame  FrameRates
}

// GlobalRates global measures taking into account the whole log file
type GlobalRates struct {
	AvgReqPerS uint64 // Average request per second since the app is on
	MaxReqPerS uint64 // Maximum number of requests since app startup
}

// FrameRates measures related to the current time-frame
type FrameRates struct {
	Duration   uint64 // Frame's duration expressed in seconds
	ReqPerS    uint64 // Frame's request-rate (req/s)
	NbRequests uint64 // Number of requests recorded during the frame's execution
	NbSuccess  uint64 // Number of successful requests recorded during the frame's execution
	NbFailures uint64 // Number of failed requests recorded during the frame's execution
}
```

Frame rates (for example the request-rate in a single time-frame) are measured using writing rate in the log file, not by reading the actual request time.
With this technique it is possible to rapidly deduce the request-rate. Indeed a frame is a constant period of time (let's call it `frame`) and during this 
frame there is a discrete and finite amount of lines that can be read (let's call this `nb_read`). Therefore the formula to get the request-rate is :
`request_rate (req/s) = nb_read (req) / frame (s)`

The measurement uncertainty (U) is :
`U (s) = server_log_write (s) + app_read_time (s)`

However given that a read or write operation from/to a file on a modern machine is in a microsecond order-of-magnitude only the server and app responsiveness 
have a serious impact. The worst case scenario would therefore be on extremely high load conditions. In this case the app would lose accuracy but in this
situation ops actions must be carried out anyway (the app would just give the cue for action). It is worth pointing out that the uncertainty impact grows
as the update time-frame shrinks.

The only incompatibility would be if the server would write its request logs in batches (retain them while it received several requests) but this is not helpful
for debugging so it wouldn't be the favoured method.

Last but not least, this method allows reading from stream-inputs such as sockets for live monitoring. 

#### Readers
This package contains input stream readers. They are used and combined to be able to load logs efficiently and asynchronously.

This package defines a `Reader` interface. It is then used by all the readers to combine their behaviour using composition. This whole
package follow the `IOC` principle (`Invertion Of Control`).
```go
// Reader is an interface to read data
type Reader interface {
	// Open prepares an object for reading
	Open(...interface{}) error
	// Read reads the object content and returns formatted logs if any
	Read() ([]log.Info, error)
	// Close closes the object and all resources used for reading the file
	Close()
}
```

##### File reader
Basic file reader, reads until `EOF` is reached. Like every reader it returns common-log-formatted-content.
*This reader is not used anymore as there is no way to move the cursor in a go scanner, so file streaming couldn't be implemented*

##### Tail reader
Basic file reader tailing a file content. It is equivalent to the linux command `tail -f`.

##### Async reader
This reader is composed of another reader to actually read the file. The only thing it adds is to be able to read a file asynchronously
(understand in a goroutine) and be able to control the reading process. To do so it defines new methods to control the reading flow.
The methods are :

This reader is not safe when calling Read so attention must be
paid when calling it (call `Stop` before).

For instance, this is how it's done to have a File reader reading asynchronously (errors are ignored for the sake of simplicity):
```go
async := reader.Async{}
async.Reader = &reader.File{}
async.Init("/tmp/access.log")

async.Start()
// .. Do something long while it's reading
async.Stop()

// The reading process has been stopped so it's OK to read the values
log := async.Read()

// Resume reading
async.Start()

// Close makes sure the reader is stopped
async.Close()
```

##### Async dbuf reader
This reader is composed of an `Async reader` composed of a `Tail reader`. To the chain it adds the implementation of
the double buffering technique. This reader is therefore able to `tail -f` a file asynchronously while being thread-safe thanks to
the double buffering technique. This is the reader used to get the logs in the application.

See the `Fetch logs` section to get more details on the technique.
