local fs        = require "lemoon.fs"
local sys       = require "lemoon.sys"
local class     = require "lemoon.class"
local throw     = require "lemoon.throw"
local filepath  = require "lemoon.filepath"
local logger    = class.new("lemoon.log","gsmake")

task.resources = function(self)

    local ok, go_tools_path = sys.lookup("go")

    if not ok then
        print(string.format("golang tools not found,visit website:https://golang.org/dl/ for more information"))
        return true
    end

    go = go_tools_path

    properties        = self.Owner.Properties.golang
    local sync              = self.Owner.Loader.Sync
    local defaultversion    = self.Owner.Loader.Config.DefaultVersion

    local tmp = self.Owner.Loader.Temp

    gopath = filepath.join(tmp,"golang")

    if not fs.exists(gopath) then
        fs.mkdir(gopath,true)
    end

    for _,dep in ipairs(properties.dependencies) do

        if not dep.version then dep.version = defaultversion end

        print(string.format("sync package [%s:%s]",dep.name,dep.version))

        local path = sync:sync(dep.name,dep.version)

        local linked = filepath.join(gopath,"src",dep.name)

        if not fs.exists(filepath.dir(linked)) then
            fs.mkdir(filepath.dir(linked),true)
        end

        if fs.exists(linked) then
            fs.rm(linked)
        end

        fs.symlink(path,linked)

        print(string.format("sync package [%s:%s] -- success",dep.name,dep.version))
    end

    linked = filepath.join(gopath,"src",self.Owner.Name)

    if not fs.exists(filepath.dir(linked)) then
        fs.mkdir(filepath.dir(linked),true)
    end

    if fs.exists(linked) then
        fs.rm(linked)
    end

    fs.symlink(self.Owner.Path,linked)

end
task.resources.Desc = "prepare dependencies package"


task.compile = function(self)

    outputdir = filepath.join(gopath,"bin")
    sys.setenv("GOPATH",gopath)

    for _,binary in ipairs(properties.binaries or {}) do
        local exec = sys.exec(go)

        if type(binary) == "table" then
            path = filepath.join(linked,binary.path)
            name = binary.name
        else
            path = filepath.join(linked,binary)
            name = binary
        end

        exec:dir(path)

        exec:start("build","-o",outputdir .. name .. sys.EXE_NAME)

        if 0 ~= exec:wait() then
            print("run golang build -- failed")
            return true
        end
    end
end

task.compile.Desc = "clang package compile task"
task.compile.Prev = "resources"

task.test = function(self,name,...)

    sys.setenv("GOPATH",gopath)

    local tests = properties.tests or {}

    if name ~= "" and name ~= nil then
        tests = { name }
    end


    for _,test in pairs( tests ) do
        local exec = sys.exec(go)
        exec:dir(filepath.join(linked,test))
        exec:start("test",...)

        if 0 ~= exec:wait() then
            print("run golang test(%s) -- failed",test)
            return true
        end
    end



end
task.test.Desc = "run golang test command"
task.test.Prev = "resources"

task.install = function(self,install_path)
    fs.list(outputdir,function(entry)
        if entry == "." or entry == ".." then
            return
        end
        fs.copy_dir(filepath.join(outputdir,entry),filepath.join(install_path,"bin",entry),fs.update_existing)
    end)
end
task.install.Desc = "clang package install package"
task.install.Prev = "compile"
