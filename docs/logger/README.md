# pkg/logger

Everything you need to know about the logger package.

## Details

My package uses [logrus](https://github.com/sirupsen/logrus) as the "*behind the scenes*" logging package.
Everything that [logrus](https://github.com/sirupsen/logrus) can get's exposed and is ready to use.

For configuration, add a **.env** file or add the required environment variables to your environment. (*checkout the docs/logger directory for more information*)

## Usage

To use the package, just ```go get``` and import it in your golang file.
After importing, just call the ```GetInstance()``` function. This will return a new logger instance or the existing one.

```golang
log, err := logger.GetInstance()
if err != nil {
    panic(err)
}
log.Info("Hello, World!")
```

## Configuration

To configure the logger and it's behaviour, it requires some environment variables. You can either set them yourself or add a **.env** file in your prjects root directory.

You can also checkout the *example.env* file in this directory, it should give you an idea of what to do.

### LOG_LEVEL

Set the loggers *logging level*. You have quit a few options:

- Trace
- Debug
- Warn/Warning
- Error
- Fatal (*calls os.Exit(1) after logging*)
- Panic (*calls panic() after logging*)

### LOG_FORMAT

You can decide between **Plaintext** (*PLAIN*) or **JSON** (*JSON*) format.

### LOG_OUTPUT

Define the output of the logger.
Choose between **File** (*FILE*) logging or the good ol' **Terminal** (*TERM*).

*For a production environment consider using file output.*

### LOG_PATH

Add the path where all the logs should get written to.
Make sure you have the rights to **write** to the direction you provide, otherwise the logger will crash and cannot create logs.

*This only applys if you specified file as your LOG_OUTPUT.*

## Todo

- Write tests
- Make and polish the docs
