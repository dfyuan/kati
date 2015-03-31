package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Func is a make function.
// http://www.gnu.org/software/make/manual/make.html#Functions
// TODO(ukai): return error instead of panic?
type Func func(*Evaluator, []string) string

func funcSubst(ev *Evaluator, args []string) string {
	Log("subst %q", args)
	// TODO: Actually, having more than three arguments is valid.
	if len(args) != 3 {
		panic(fmt.Sprintf("*** insufficient number of arguments (%d) to function `subst'."))
	}
	from := ev.evalExpr(args[0])
	to := ev.evalExpr(args[1])
	text := ev.evalExpr(args[2])
	return strings.Replace(text, from, to, -1)
}

func funcPatsubst(ev *Evaluator, args []string) string {
	Log("patsubst %q", args)
	// TODO: Actually, having more than three arguments is valid.
	if len(args) != 3 {
		panic(fmt.Sprintf("*** insufficient number of arguments (%d) to function `patsubst'."))
	}
	pat := ev.evalExpr(args[0])
	repl := ev.evalExpr(args[1])
	texts := splitSpaces(ev.evalExpr(args[2]))
	for i, text := range texts {
		texts[i] = substPattern(pat, repl, text)
	}
	return strings.Join(texts, " ")
}

// http://www.gnu.org/software/make/manual/make.html#File-Name-Functions
func funcWildcard(ev *Evaluator, args []string) string {
	Log("wildcard %q", args)
	pattern := ev.evalExpr(strings.Join(args, ","))
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	return strings.Join(files, " ")
}

// http://www.gnu.org/software/make/manual/make.html#Shell-Function
func funcShell(ev *Evaluator, args []string) string {
	Log("shell %q", args)
	arg := ev.evalExpr(strings.Join(args, ","))
	cmdline := []string{"/bin/sh", "-c", arg}
	cmd := exec.Cmd{
		Path: cmdline[0],
		Args: cmdline,
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	re, err := regexp.Compile(`\s`)
	if err != nil {
		panic(err)
	}
	return string(re.ReplaceAllString(string(out), " "))
}

// http://www.gnu.org/software/make/manual/make.html#Make-Control-Functions
func funcWarning(ev *Evaluator, args []string) string {
	Log("warning %q", args)
	arg := ev.evalExpr(strings.Join(args, ","))
	fmt.Printf("%s:%d: %s\n", ev.filename, ev.lineno, arg)
	return ""
}
