package stage

type ShellScriptStage struct {
	command CommandStage
}

func (self *ShellScriptStage) Run() bool {
	return self.command.Run()
}

func (self *ShellScriptStage) AddScript(scriptFile string) {
	// TODO: validate the existance of scriptFile
	// and flush log when the file does not exist.
	self.command.AddCommand("bash", scriptFile)
}

func NewShellScriptStage() *ShellScriptStage {
	return &ShellScriptStage{}
}
