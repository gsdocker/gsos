name "github.com/gsdocker/gsos"

plugin "github.com/gsmake/golang"


golang = {
    dependencies = {
        { name = "github.com/gsdocker/gserrors" };
        { name = "github.com/gsdocker/gslogger" };
        { name = "github.com/go-fsnotify/fsnotify" };
    };

    tests = { "fs","uuid" }
}
