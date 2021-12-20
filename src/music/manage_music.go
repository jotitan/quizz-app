package music

import (
	"fmt"
	"github.com/quizz-app/logger"
	"os"
	"os/exec"
	"strings"
)

type Cutter struct {
	ffmpegPath string
}

func NewCutter(ffmpegPath string)Cutter {
	return Cutter{ffmpegPath}
}

func (c Cutter)Cut(input,output string, from,to int)error{
	params := strings.Split(fmt.Sprintf("-i %s -ss %d -to %d %s",input,from,to,output)," ")
	logger.GetLogger2().Info("Convert to " + output)
	//cmd := exec.Command(path)
	cmd := exec.Command(c.ffmpegPath,params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
