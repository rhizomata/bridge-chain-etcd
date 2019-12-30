package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhizomata/bridge-chain-etcd/kernel"
	"github.com/rhizomata/bridge-chain-etcd/kernel/job"
)

// BuiltinService ..
type BuiltinService struct {
	kernel *kernel.Kernel
}

func (service BuiltinService) health(context *gin.Context) {
	checkFrom := context.GetHeader("Check-From")
	fmt.Println("checkFrom : ", checkFrom)
	context.Writer.WriteString("OK")
	context.Writer.Flush()
}

func (service BuiltinService) addJob(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}

	job := service.kernel.GetJobManager().AddJob(job.NewJob(data))
	data, err = json.Marshal(job)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	context.Writer.Write(data)
	context.Writer.Flush()
}
