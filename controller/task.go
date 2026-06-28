package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-alpha/models"
	"go-alpha/response"
)

type TaskController struct{}

func (TaskController) ListTasks(c *gin.Context) {
	tasks := models.Task{}.GetAllTasks()
	response.Success("查询所有任务成功", tasks, c)
}

func (TaskController) AddTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.AddTask()
	response.Success("添加任务成功", &task, c)
}

func (TaskController) ToggleActive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	task := models.Task{}.GetTaskById(uint(id))
	if task.ID == 0 {
		response.Failed("任务不存在", c)
		return
	}

	var body struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ToggleActive(body.Active)
	response.Success("更新任务状态成功", task, c)
}

func (TaskController) DeleteTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	task := models.Task{}.GetTaskById(uint(id))
	if task.ID == 0 {
		response.Failed("任务不存在", c)
		return
	}
	models.Task{}.DeleteTask(uint(id))
	response.Success("删除任务成功", task, c)
}
