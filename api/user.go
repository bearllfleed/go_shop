package api

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/bearllflee/go_shop/global"
	"github.com/bearllflee/go_shop/model"
	"github.com/bearllflee/go_shop/model/request"
	"github.com/bearllflee/go_shop/model/response"
	"github.com/bearllflee/go_shop/utils"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/bearllflee/go_shop/service"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req request.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// log.Println("参数错误: ", utils.Translate(err))
		global.Logger.Error("参数错误", zap.String("err", utils.Translate(err)))
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	user, err := service.UserServiceApp.Login(req)
	if err != nil {
		if errors.Is(err, global.ErrUserNotFound) || errors.Is(err, global.ErrPasswordIncorrect) {
			global.Logger.Error("登录失败", zap.String("err", err.Error()))
			response.FailWithMessage(err.Error(), c)
			return
		} else {
			log.Println("登录失败: ", err)
			response.FailWithMessage("登录失败", c)
			return
		}
	}
	// 生成token
	jwt := utils.NewJwt()
	claims := jwt.CreateClaims(model.BaseClaims{
		UserId:   user.ID,
		Username: user.Username,
	})
	token, err := jwt.GenerateToken(&claims)
	if err != nil {
		log.Println("生成token失败: ", err)
		response.FailWithMessage("生成token失败", c)
		return
	}
	// 将用户信息缓存到`redis`中，对应的操作应该是`HASH`。
	// userJSON, err := json.Marshal(user)
	// if err != nil {
	// 	log.Println("序列化用户信息失败: ", err)
	// 	return
	// }
	// err = global.Redis.HSet(context.Background(), "online_user", user.ID, userJSON).Err()
	// if err != nil {
	// 	log.Println("缓存用户信息失败: ", err)
	// }
	global.Logger.Info("登录成功", zap.String("username", user.Username))
	response.OkWithData(&response.LoginResponse{
		User:  user,
		Token: token,
	}, c)
}

func Register(c *gin.Context) {
	var req request.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("参数错误: ", utils.Translate(err))
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	user, err := service.UserServiceApp.Register(req)
	if err != nil {
		if errors.Is(err, global.ErrUserAlreadyExists) {
			response.FailWithMessage(err.Error(), c)
			return
		} else {
			log.Println("注册失败: ", err)
			response.FailWithMessage("注册失败", c)
			return
		}
	}
	response.OkWithData(user, c)
}

func UserList(c *gin.Context) {
	var req request.UserListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("参数错误: ", utils.Translate(err))
		response.FailWithMessage(utils.Translate(err), c)
		return
	}
	total, users, err := service.UserServiceApp.UserList(req)
	if err != nil {
		log.Println("获取用户列表失败: ", err)
		response.FailWithMessage("获取用户列表失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		Total:    total,
		List:     users,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = map[*Client]bool{}
var mu = sync.RWMutex{}

type Client struct {
	conn *websocket.Conn
}

type Message struct {
	Type    string
	Content string
}

func OnlineTool(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Logger.Error("建立连接失败：", zap.String("err", err.Error()))
		return
	}
	client := Client{
		conn: conn,
	}
	mu.Lock()
	clients[&client] = true
	mu.Unlock()
	go client.readMessage()
}

func (c *Client) readMessage() {
	defer func() {
		c.conn.Close()
		mu.Lock()
		delete(clients, c)
		mu.Unlock()
	}()
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			global.Logger.Error("读取消息失败：", zap.String("err", err.Error()))
			return
		}
		broadCast(data)
	}
}

func broadCast(data []byte) {
	for c := range clients {
		err := c.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			global.Logger.Error("发送消息失败", zap.String("client", c.conn.RemoteAddr().String()), zap.String("err", err.Error()))
			continue
		}
	}
}
