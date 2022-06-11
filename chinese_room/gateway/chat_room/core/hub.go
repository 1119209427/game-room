package core

import (
	"fmt"
	"gateway/chat_room/cache"
	"gateway/chat_room/chess_judge"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)
var( upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
)

type ChessMsg struct {
	Hflag bool  `json:"hflag"`//表示红是否出棋
	Bflag bool `json:"bflag"`//表示黑是否出棋子
	RoomId string `json:"room_id"`
	Start string `json:"start"`
	End string `json:"end"`
	ChesserType int64   `json:"chesser_type "`//1表示红方，2表示黑方
	ChessType int64   `json:"chess_type"`//1=帅(将)，2=仕(士)，3=相(象),4=马,5=炮，6=车，7=兵(卒)
	Type int64  `json:"type"` //1=求和，2=认输
	Ready int64 `json:"ready"`
}
/*type ChessParam struct {
	Hflag string  `form:"hflag"`//表示红是否出棋
	Bflag string`form:"bflag"`//表示黑是否出棋子
	ChesserType string  `form:"chesser_type"`//1表示红方，2表示黑方
	ChessType string `form:"chess_type"`//1=帅(将)，2=仕(士)，3=相(象),4=马,5=炮，6=车，7=兵(卒)
	Type string `form:"type"`//1=求和，2=认输
}*/

type ChessConn struct {
	ws            *websocket.Conn
	send          chan []byte
	ChesserType int64  //1表示红方，2表示黑方
	ChessType int64  //1=帅(将)，2=仕(士)，3=相(象),4=马,5=炮，6=车，7=兵(卒)
	hflag bool //表示红是否出棋
	bflag bool //表示黑是否出棋子
	Type int64 //1=求和，2=认输,3=继续对弈
	Ready int64
}
type Connection struct {
	ws            *websocket.Conn
	send          chan []byte
	numberv       int
	forbiddenword bool
	timelog       int64
}

// ReadPump 将聊天的消息读入管道
func (m Message)ReadPump(){
	c:=m.conn
	//防止浪费资源
	defer func() {
		H.unprepare<-m
		H.unlogin<- m
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for{
		_,msg,err:=c.ws.ReadMessage()
		if err!=nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			//log.Printf("error: %v", err)
			break
		}
		go m.Kickout(msg)
	}

}
func(cm ChessMessage) RedaChessMessage() {
	c := cm.Conn
	defer func() {
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	//初始化棋盘
	chess_judge.InitChess()
	for {
		msg := new(ChessMsg)
		if err := c.ws.ReadJSON(&msg); err != nil {
			log.Println("读取json数据失败")
			H.unregister <- cm
			c.ws.Close()
			break
		}
		msg.Hflag = true
		msg.Bflag = false
		//判断是否合理
		//首次红方先出棋子
		if msg.ChesserType == 1 && msg.Hflag { //说明是红方

			if msg.ChessType == 1 {
				//说明是将，判断逻辑
				flag := chess_judge.KingMove(msg.Start, msg.End, msg.ChesserType)
				if flag { //移动逻辑正确，把消息放入
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}
				}
			}
				if msg.ChessType == 2 {
					//说明是士，判断逻辑
					flag := chess_judge.MandarinsMove(msg.Start, msg.End, msg.ChesserType)
					if flag {
						H.register <- cm
						//将消息存入redis
						err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
						if err != nil {
							log.Println("数据存入缓存失败")
						}

					}
				}
				if msg.ChessType == 3 {
					//说明是相
					flag := chess_judge.ElephantsMove(msg.Start, msg.End, msg.ChesserType)
					if flag {
						H.register <- cm
						//将消息存入redis
						err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
						if err != nil {
							log.Println("数据存入缓存失败")
						}

					}
				}
				if msg.ChessType == 4 {
					//说明是马
					flag := chess_judge.KingMove(msg.Start, msg.End, msg.ChesserType)
					if flag {
						H.register <- cm
						//将消息存入redis
						err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
						if err != nil {
							log.Println("数据存入缓存失败")
						}
					}
				}
					if msg.ChessType == 5 {
						flag := chess_judge.CannonsMove(msg.Start, msg.End, msg.ChesserType)
						if flag {
							H.register <- cm
							//将消息存入redis
							err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
							if err != nil {
								log.Println("数据存入缓存失败")
							}
						}
					}
					if msg.ChessType==6{
						flag:=chess_judge.RooksMove(msg.Start, msg.End, msg.ChesserType)
						if flag {
							H.register <- cm
							//将消息存入redis
							err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
							if err != nil {
								log.Println("数据存入缓存失败")
							}
						}
					}
					if msg.ChessType==7{
						flag:=chess_judge.PawnsMove(msg.Start, msg.End, msg.ChesserType)
						if flag{
							H.register <- cm
							//将消息存入redis
							err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
							if err != nil {
								log.Println("数据存入缓存失败")
							}
						}
					}
					//赋值让黑方走了
					msg.Bflag = true
				}
		if msg.ChessType == 2 && msg.Bflag { //说明是黑方
			if msg.ChessType == 1 {
				//说明是将，判断逻辑
				flag := chess_judge.KingMove(msg.Start, msg.End, msg.ChesserType)
				if flag { //移动逻辑正确，把消息放入
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}
				}
			}
			if msg.ChessType == 2 {
				//说明是士，判断逻辑
				flag := chess_judge.MandarinsMove(msg.Start, msg.End, msg.ChesserType)
				if flag {
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}

				}
			}
			if msg.ChessType == 3 {
				//说明是相
				flag := chess_judge.ElephantsMove(msg.Start, msg.End, msg.ChesserType)
				if flag {
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}

				}
			}
			if msg.ChessType == 4 {
				//说明是马
				flag := chess_judge.KingMove(msg.Start, msg.End, msg.ChesserType)
				if flag {
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}
				}
			}
			if msg.ChessType == 5 {
				flag := chess_judge.CannonsMove(msg.Start, msg.End, msg.ChesserType)
				if flag {
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}
				}
			}
			if msg.ChessType==6{
				flag:=chess_judge.RooksMove(msg.Start, msg.End, msg.ChesserType)
				if flag {
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}
				}
			}
			if msg.ChessType==7{
				flag:=chess_judge.PawnsMove(msg.Start, msg.End, msg.ChesserType)
				if flag{
					H.register <- cm
					//将消息存入redis
					err := cache.Set(msg.RoomId, msg.Start+msg.End, 7*time.Hour)
					if err != nil {
						log.Println("数据存入缓存失败")
					}
				}
			}
		}
		if msg.ChessType==1&&msg.Bflag{ //红方想在黑方出棋的时候出棋
			H.illegalChess<-cm
			break
		}
		if msg.ChessType==2&&msg.Hflag{//黑方想在红方出棋的时候出棋
			H.illegalChess<-cm
			break
		}
	}
}




// Kickout 信息处理，不合法言论 禁言警告，超过3次，踢出游戏 ；
func(m Message) Kickout(msg []byte){
		c := m.conn
		// 判断是否有禁言时间,并超过5分钟禁言时间,没有超过进入禁言提醒
		nowT := int64(time.Now().Unix())
		if nowT-c.timelog < 300 {
			H.warnmsg <- m
		}
		if c.numberv < 3 {
			basestr := "死亡崩薨"
			teststr := string(msg[:])
			for _, i := range teststr {
				flag := strings.Contains(basestr, string(i))
				if flag {
					c.numberv += 1
					c.forbiddenword = true
					//记录禁言时间
					c.timelog = int64(time.Now().Unix())
					H.warnmsg <- m
					break
				}
			}
		}
		//不禁言可以发送
		if c.forbiddenword !=true{
			M := Message{msg, m.roomid, c}
			H.broadcast <- M
		} else {
			H.kickoutroom <- m
			fmt.Println("要被提出群聊了")
			c.ws.Close()
		}
	}
func(c *Connection) write(mt int, payload[]byte)error{
		c.ws.SetWriteDeadline(time.Now().Add(writeWait))
		return c.ws.WriteMessage(mt, payload)
	}
func (cc *ChessConn)writeChess(mt int,payload[]byte)error{
	cc.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return cc.ws.WriteMessage(mt, payload)
}
func (cs *ChessMessage)WriteChessPump(){
	c:=cs.Conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.writeChess(websocket.CloseMessage, []byte{})
				return

			}
			err := c.writeChess(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			if err := c.writeChess(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}

}
func(s *Message) WritePump() {
		c := s.conn
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			c.ws.Close()
		}()
		//写出来
		for {
			select {
			case message, ok := <-c.send:
				if !ok {
					c.write(websocket.CloseMessage, []byte{})
					return

				}
				err := c.write(websocket.TextMessage, message)
				if err != nil {
					log.Println(err)
					return
				}
			case <-ticker.C:
				if err := c.write(websocket.PingMessage, []byte{}); err != nil {
					return
				}

			}
		}
	}
func ServerWs(ctx * gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
		}
		roomid := ctx.Param("room_id")
		/*var chess ChessParam
		if ctx.ShouldBind(&chess)!=nil{
			ctx.JSON(401,gin.H{"status":"表单绑定失败"})
		}


		Type,err:=strconv.ParseInt(chess.Type,10,64)
		if err!=nil{
		log.Println("ParseInt",err)
		}
		chesserType,err:=strconv.ParseInt(chess.ChesserType,10,64)
		if err!=nil{
		log.Println("ParseInt",err)
		}
		chessType,err:=strconv.ParseInt(chess.ChessType,10,64)
		if err!=nil{
			log.Println("ParseInt",err)
		}
		hflag,err:=strconv.ParseBool(chess.Hflag)
		if err!=nil{
			log.Println("ParseBool",err)
		}
		bflag,err:=strconv.ParseBool(chess.Bflag)
		if err!=nil{
		log.Println("ParseBool",err)
		}
		ch:=&ChessConn{send: make(chan []byte, 256),ws: conn,ChessType: chessType,ChesserType:chesserType,Type: Type,hflag: hflag,bflag: bflag}
		cm:=ChessMessage{nil,roomid,ch}
		H.register<-cm
*/

		//H.prepare <- m
		c := &Connection{send: make(chan []byte, 256), ws: conn}
		m := Message{nil, roomid, c}
		H.login <- m
		go m.WritePump()
		go m.ReadPump()
		/*go cm.RedaChessMessage()*/
}
func ChessWs(ctx *gin.Context){
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	roomid := ctx.Query("room_id")
	c:=&ChessConn{send: make(chan []byte, 256), ws: conn}
	m:=ChessMessage{nil,roomid,c}
	H.register<-m
	go m.WriteChessPump()
	go m.WriteChessPump()
}



